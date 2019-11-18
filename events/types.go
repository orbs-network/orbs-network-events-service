package events

import "github.com/orbs-network/orbs-client-sdk-go/codec"

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
	StoreEvent(blockHeight uint64, timestamp int64, event *codec.Event) error
	GetBlockHeight() (uint64, error)
	StoreBlockHeight(blockHeight uint64, timestamp int64) error
	GetEvents(query *FilterQuery) ([]*StoredEvent, error)
}
