package storage

import (
	"fmt"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/scribe/log"
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
	if filterQuery.FromTime() != 0 || filterQuery.ToTime() != 0 {
		fromBlock := s.getBlockHeightByTimestamp(filterQuery.FromTime(), true)
		toBlock := s.getBlockHeightByTimestamp(filterQuery.ToTime(), false)

		newQuery := types.IndexerRequestReader(filterQuery.Raw())
		newQuery.MutateFromTime(0)
		newQuery.MutateToTime(0)
		newQuery.MutateFromBlock(fromBlock)
		newQuery.MutateToBlock(toBlock)

		return s.GetEvents(newQuery)
	}

	var filters []*protocol.ArgumentArray
	for i := filterQuery.FiltersIterator(); i.HasNext(); {
		filters = append(filters, protocol.ArgumentArrayReader(i.NextFilters()))
	}

	err = s.db.View(func(tx *bolt.Tx) error {
		for iterator := filterQuery.EventNameIterator(); iterator.HasNext(); {
			eventName := iterator.NextEventName()
			tableName := getEventsBucketName(filterQuery.ContractName(), eventName)
			eventsBucket := tx.Bucket([]byte(tableName))
			if eventsBucket == nil {
				// return empty array
				return nil
			}

			if filterQuery.FromBlock() != 0 || filterQuery.ToBlock() != 0 {
				cursor := eventsBucket.Cursor()
				toBlock := filterQuery.ToBlock()
				if toBlock == 0 {
					toBlock = s.GetBlockHeight()
				}

				lastBlock := uint64(0)
				for i := filterQuery.FromBlock(); i <= toBlock; i++ {
					blockHeightRaw, indexedEventRaw := cursor.Seek(ToBytes(i))

					if blockHeightRaw == nil {
						break
					} else if blockHeight := ReadUint64(blockHeightRaw); blockHeight == lastBlock {
						break
					} else if blockHeight <= filterQuery.ToBlock() {
						event := types.IndexedEventReader(indexedEventRaw)
						lastBlock = blockHeight
						i = event.BlockHeight()
						if matchEvent(event, filters) {
							events = append(events, event)
						}
					} else {
						break
					}
				}
			} else {
				eventsBucket.ForEach(func(blockHeightRaw, indexedEventRaw []byte) error {
					event := types.IndexedEventReader(indexedEventRaw)
					if matchEvent(event, filters) {
						events = append(events, event)
					}
					return nil
				})
			}
		}

		return nil
	})

	return
}

func (s *storage) GetBlockHeight() (blockHeight uint64) {
	blockHeight, _ = s.getLastBlockHeightAndTimestamp()
	return
}

func (s *storage) getLastBlockHeightAndTimestamp() (blockHeight uint64, timestamp uint64) {
	s.db.View(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte("blocks"))
		if blocksBucket == nil {
			return nil
		}

		timestampRaw, blockHeightRaw := blocksBucket.Cursor().Last()
		blockHeight = ReadUint64(blockHeightRaw)
		timestamp = ReadUint64(timestampRaw)

		return nil
	})

	return
}

func (s *storage) storeBlockHeight(tx *bolt.Tx, blockHeight uint64, timestamp uint64) (err error) {
	blocksBucket, err := tx.CreateBucketIfNotExists([]byte("blocks"))
	if err != nil {
		return err
	}

	return blocksBucket.Put(ToBytes(int64(timestamp)), ToBytes(uint64(blockHeight)))
}

func (s *storage) getBlockHeightByTimestamp(timestamp uint64, forward bool) (blockHeight uint64) {
	s.db.View(func(tx *bolt.Tx) error {
		blocksBucket := tx.Bucket([]byte("blocks"))
		if blocksBucket == nil {
			return nil
		}

		lastBlockHeight, lastTimestamp := s.getLastBlockHeightAndTimestamp()
		if timestamp > lastTimestamp {
			blockHeight = lastBlockHeight
			return nil
		}

		closestTimestampRaw, blockHeightRaw := blocksBucket.Cursor().Seek(ToBytes(timestamp))
		closestTimestamp := ReadUint64(closestTimestampRaw)
		blockHeight = ReadUint64(blockHeightRaw)

		if closestTimestamp != timestamp {
			if forward {
				// FIXME edge cases
			} else {
				blockHeight -= 1
			}
		}

		return nil
	})

	return
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
