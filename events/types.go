package events

import "github.com/orbs-network/orbs-spec/types/go/protocol"

type StoredEvent struct {
	ContractName string
	EventName    string

	BlockHeight uint64
	Timestamp   int64
	Arguments   []interface{}
}

type FilterQuery struct {
	ContractName string
	EventNames   []string

	FromBlock uint64
	ToBlock   uint64

	FromTime uint64
	ToTime   uint64
}

type Storage interface {
	StoreEvents(blockHeight uint64, timestamp int64, events []*protocol.IndexedEvent) error
	GetBlockHeight() uint64
	GetEvents(query *FilterQuery) ([]*protocol.IndexedEvent, error)
	Shutdown() error
}
