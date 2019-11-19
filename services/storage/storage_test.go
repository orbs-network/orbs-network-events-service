package storage

import (
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
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

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE, false)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*protocol.IndexedEvent{DEFAULT_EVENT})
	require.NoError(t, err, "could not store event")

	blockHeight := storage.GetBlockHeight()
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT, blockHeight)

	eventList, err := storage.GetEvents((&client.IndexerRequestBuilder{
		ContractName: primitives.ContractName(DEFAULT_EVENT.ContractName().String()),
		EventName:    []string{DEFAULT_EVENT.EventName()},
	}).Build())
	require.NoError(t, err)
	require.Len(t, eventList, 1)
	require.EqualValues(t, DEFAULT_EVENT.Raw(), eventList[0].Raw())

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT+100, DEFAULT_TIME, []*protocol.IndexedEvent{NEXT_EVENT})
	require.NoError(t, err, "could not store event")

	updatedBlockHeight := storage.GetBlockHeight()
	require.NoError(t, err)
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT+100, updatedBlockHeight)

	eventList, err = storage.GetEvents((&client.IndexerRequestBuilder{
		ContractName: primitives.ContractName(DEFAULT_EVENT.ContractName().String()),
		EventName:    []string{NEXT_EVENT.EventName()},
	}).Build())
	require.NoError(t, err)
	require.Len(t, eventList, 2)
}
