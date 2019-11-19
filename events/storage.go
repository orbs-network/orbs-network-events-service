package events

import (
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/scribe/log"
	bolt "go.etcd.io/bbolt"
	"strings"
	"time"
)

type storage struct {
	logger log.Logger
	db     *bolt.DB
}

func NewStorage(logger log.Logger, dataSource string) (Storage, error) {
	boltDb, err := bolt.Open(dataSource, 0600, &bolt.Options{
		Timeout:  5 * time.Second,
		ReadOnly: false,
	})
	if err != nil {
		return nil, err
	}

	return &storage{
		logger,
		boltDb,
	}, nil
}

func (s *storage) StoreEvents(blockHeight uint64, timestamp int64, events []*codec.Event) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			println("rolling back!")
			tx.Rollback()
		}
	}()

	if err := s.storeBlockHeight(tx, blockHeight, timestamp); err != nil {
		return err
	}

	for _, event := range events {
		if err := s.storeEvent(tx, blockHeight, timestamp, event); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *storage) storeEvent(tx *bolt.Tx, blockHeight uint64, timestamp int64, event *codec.Event) error {
	tableName := getEventsBucketName(event.ContractName, event.EventName)
	eventsBucket, err := tx.CreateBucketIfNotExists([]byte(tableName))
	if err != nil {
		return err
	}

	argumentsWithTimestamp := []interface{}{
		uint64(timestamp),
	}
	argumentsWithTimestamp = append(argumentsWithTimestamp, event.Arguments...)
	arguments, err := protocol.ArgumentArrayFromNatives(argumentsWithTimestamp)

	s.logger.Info("Storing event",
		log.Int64("blockTimestamp", timestamp),
		log.Uint64("blockHeight", blockHeight),
		log.String("contractName", event.ContractName),
		log.String("eventName", event.EventName),
		log.String("arguments", fmt.Sprintf("%v", event.Arguments)))

	if err != nil {
		return err
	}

	return eventsBucket.Put(ToBytes(blockHeight), arguments.Raw())
}

func (s *storage) GetEvents(filterQuery *FilterQuery) (events []*StoredEvent, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		for _, eventName := range filterQuery.EventNames {
			tableName := getEventsBucketName(filterQuery.ContractName, eventName)
			eventsBucket := tx.Bucket([]byte(tableName))

			eventsBucket.ForEach(func(blockHeightRaw, argumentsWithTimestampRaw []byte) error {
				argumentsWithTimestamp, err := protocol.PackedOutputArgumentsToNatives(argumentsWithTimestampRaw)
				if err != nil {
					return err
				}

				events = append(events, &StoredEvent{
					ContractName: filterQuery.ContractName,
					EventName:    eventName,
					BlockHeight:  ReadUint64(blockHeightRaw),
					Timestamp:    int64(argumentsWithTimestamp[0].(uint64)),
					Arguments:    argumentsWithTimestamp[1:],
				})

				return nil
			})
		}

		return nil
	})

	return
}

func (s *storage) GetBlockHeight() (value uint64) {
	s.db.View(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte("blocks"))
		if blocksBucket == nil {
			return nil
		}

		blockHeightRaw, _ := blocksBucket.Cursor().Last()
		value = ReadUint64(blockHeightRaw)

		return nil
	})

	return
}

func (s *storage) storeBlockHeight(tx *bolt.Tx, blockHeight uint64, timestamp int64) (err error) {
	blocksBucket, err := tx.CreateBucketIfNotExists([]byte("blocks"))
	if err != nil {
		return err
	}

	return blocksBucket.Put(ToBytes(blockHeight), ToBytes(timestamp))
}

func (s *storage) Shutdown() (err error) {
	if err = s.db.Sync(); err != nil {
		s.logger.Error("failed to synchronize storage on shutdown")
	}

	if err = s.db.Close(); err != nil {
		s.logger.Error("failed to close storage on shutdown")
	}

	s.logger.Info("storage shut down")

	return
}

func getEventsBucketName(contractName string, eventName string) string {
	return strings.Join([]string{"events", contractName, eventName}, ".")
}
