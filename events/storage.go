package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	bolt "go.etcd.io/bbolt"
	"strings"
	"time"
)

type storage struct {
	boltDB *bolt.DB
}

func NewStorage(dataSource string) (Storage, error) {
	boltDb, err := bolt.Open(dataSource, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}

	return &storage{
		boltDb,
	}, nil
}

func (s *storage) StoreEvent(blockHeight uint64, timestamp int64, event *codec.Event) error {
	tx, err := s.boltDB.Begin(true)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			println("rolling back!")
			tx.Rollback()
		}
	}()

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
	if err != nil {
		return err
	}
	err = eventsBucket.Put(SimpleSerialization(blockHeight), arguments.Raw())
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *storage) GetEvents(filterQuery *FilterQuery) (events []*StoredEvent, err error) {
	err = s.boltDB.View(func(tx *bolt.Tx) error {
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

func (s *storage) GetBlockHeight() (value uint64, err error) {
	err = s.boltDB.View(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte("blocks"))
		blockHeightRaw, _ := blocksBucket.Cursor().First()
		value = ReadUint64(blockHeightRaw)

		return nil
	})

	return
}

func (s *storage) StoreBlockHeight(blockHeight uint64, timestamp int64) (err error) {
	tx, err := s.boltDB.Begin(true)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			println("rolling back!")
			tx.Rollback()
		}
	}()

	blocksBucket, err := tx.CreateBucketIfNotExists([]byte("blocks"))
	if err != nil {
		return err
	}

	err = blocksBucket.Put(SimpleSerialization(blockHeight), SimpleSerialization(timestamp))
	if err != nil {
		return err
	}

	return tx.Commit()
}

func getEventsBucketName(contractName string, eventName string) string {
	return strings.Join([]string{"events", contractName, eventName}, ".")
}
