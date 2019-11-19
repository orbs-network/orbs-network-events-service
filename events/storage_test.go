package events

import (
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var arizonaArgs, _ = protocol.PackedInputArgumentsFromNatives([]interface{}{
	"Raising Arizona", uint32(1987), "Nicolas Cage",
})

var DEFAULT_EVENT = (&protocol.IndexedEventBuilder{
	ContractName:    primitives.ContractName("SomeContract"),
	EventName:       "MovieRelease",
	BlockHeight:     DEFAULT_BLOCK_HEIGHT,
	ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
	Timestamp:       DEFAULT_TIME,
	Arguments:       arizonaArgs,
}).Build()

var NEXT_EVENT = (&protocol.IndexedEventBuilder{
	ContractName:    primitives.ContractName("SomeContract"),
	EventName:       "MovieRelease",
	BlockHeight:     DEFAULT_BLOCK_HEIGHT + 100,
	ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
	Timestamp:       DEFAULT_TIME + 5000,
	Arguments:       arizonaArgs,
}).Build()

const DATA_SOURCE = "test.bolt"

const DEFAULT_BLOCK_HEIGHT = primitives.BlockHeight(1974)

var DEFAULT_TIME = primitives.TimestampNano(time.Now().UnixNano())

func removeDB() {
	os.RemoveAll(DATA_SOURCE)
}

func TestStorage_StoreEvent(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(uint64(DEFAULT_BLOCK_HEIGHT), int64(DEFAULT_TIME), []*protocol.IndexedEvent{DEFAULT_EVENT})
	require.NoError(t, err, "could not store event")

	blockHeight := storage.GetBlockHeight()
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT, blockHeight)

	eventList, err := storage.GetEvents(&FilterQuery{
		ContractName: DEFAULT_EVENT.ContractName().String(),
		EventNames:   []string{DEFAULT_EVENT.EventName()},
	})
	require.NoError(t, err)
	require.Len(t, eventList, 1)
	require.EqualValues(t, DEFAULT_EVENT.Raw(), eventList[0].Raw())

	err = storage.StoreEvents(uint64(DEFAULT_BLOCK_HEIGHT+100), int64(DEFAULT_TIME), []*protocol.IndexedEvent{NEXT_EVENT})
	require.NoError(t, err, "could not store event")

	updatedBlockHeight := storage.GetBlockHeight()
	require.NoError(t, err)
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT+100, updatedBlockHeight)

	eventList, err = storage.GetEvents(&FilterQuery{
		ContractName: NEXT_EVENT.ContractName().String(),
		EventNames:   []string{NEXT_EVENT.EventName()},
	})
	require.NoError(t, err)
	require.Len(t, eventList, 2)
}
