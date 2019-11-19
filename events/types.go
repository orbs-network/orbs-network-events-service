package events

import (
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

type FilterQuery struct {
	ContractName string
	EventNames   []string

	FromBlock uint64
	ToBlock   uint64

	FromTime uint64
	ToTime   uint64
}

type Storage interface {
	StoreEvents(blockHeight primitives.BlockHeight, timestamp primitives.TimestampNano, events []*protocol.IndexedEvent) error
	GetBlockHeight() primitives.BlockHeight
	GetEvents(query *FilterQuery) ([]*protocol.IndexedEvent, error)
	Shutdown() error
}
