package storage

import (
	"fmt"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/scribe/log"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
	"strings"
	"time"
)

type Storage interface {
	StoreEvents(blockHeight uint64, timestamp uint64, events []*types.IndexedEvent) error
	GetBlockHeight() uint64
	GetEvents(query *types.IndexerRequest) ([]*types.IndexedEvent, error)
	Shutdown() error
}

type storage struct {
	logger log.Logger
	db     *bolt.DB
}

func NewStorage(logger log.Logger, dataSource string, readOnly bool) (Storage, error) {
	boltDb, err := bolt.Open(dataSource, 0600, &bolt.Options{
		Timeout:  5 * time.Second,
		ReadOnly: readOnly,
	})
	if err != nil {
		return nil, err
	}

	return &storage{
		logger,
		boltDb,
	}, nil
}

func NewStorageForChain(logger log.Logger, dbPath string, vcid uint32, readOnly bool) (Storage, error) {
	return NewStorage(logger, fmt.Sprintf("%s/vchain-%d.bolt", dbPath, vcid), readOnly)
}

func (s *storage) StoreEvents(blockHeight uint64, timestamp uint64, events []*types.IndexedEvent) error {
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			s.logger.Error("rolling back!")
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

func (s *storage) storeEvent(tx *bolt.Tx, blockHeight uint64, timestamp uint64, event *types.IndexedEvent) error {
	tableName := getEventsBucketName(event.ContractName(), event.EventName())
	eventsBucket, err := tx.CreateBucketIfNotExists([]byte(tableName))
	if err != nil {
		return err
	}

	s.logger.Info("Storing event",
		log.Int64("blockTimestamp", int64(timestamp)),
		log.Uint64("blockHeight", uint64(blockHeight)),
		log.String("contractName", event.ContractName()),
		log.String("eventName", event.EventName()))

	return eventsBucket.Put(ToBytes(uint64(blockHeight)), event.Raw())
}

func (s *storage) GetEvents(filterQuery *types.IndexerRequest) (events []*types.IndexedEvent, err error) {
	err = s.db.View(func(tx *bolt.Tx) error {
		for iterator := filterQuery.EventNameIterator(); iterator.HasNext(); {
			eventName := iterator.NextEventName()
			tableName := getEventsBucketName(filterQuery.ContractName(), eventName)
			eventsBucket := tx.Bucket([]byte(tableName))
			if eventsBucket == nil {
				return errors.Errorf("bucket %s not found", tableName)
			}

			eventsBucket.ForEach(func(blockHeightRaw, indexedEventRaw []byte) error {
				events = append(events, types.IndexedEventReader(indexedEventRaw))

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
		value = uint64(ReadUint64(blockHeightRaw))

		return nil
	})

	return
}

func (s *storage) storeBlockHeight(tx *bolt.Tx, blockHeight uint64, timestamp uint64) (err error) {
	blocksBucket, err := tx.CreateBucketIfNotExists([]byte("blocks"))
	if err != nil {
		return err
	}

	return blocksBucket.Put(ToBytes(uint64(blockHeight)), ToBytes(int64(timestamp)))
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
