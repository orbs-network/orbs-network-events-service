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

var arizonaArgsNative = []interface{}{
	"Raising Arizona", uint32(1987), "Nicolas Cage",
}

var arizonaArgs, _ = protocol.ArgumentArrayFromNatives(arizonaArgsNative)

var DEFAULT_EVENT = (&protocol.IndexedEventBuilder{
	ContractName:    primitives.ContractName("SomeContract"),
	EventName:       "MovieRelease",
	BlockHeight:     DEFAULT_BLOCK_HEIGHT,
	ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
	Timestamp:       DEFAULT_TIME,
	Arguments:       arizonaArgs.Raw(),
}).Build()

var NEXT_EVENT = (&protocol.IndexedEventBuilder{
	ContractName:    primitives.ContractName("SomeContract"),
	EventName:       "MovieRelease",
	BlockHeight:     DEFAULT_BLOCK_HEIGHT + 100,
	ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
	Timestamp:       DEFAULT_TIME + 5000,
	Arguments:       arizonaArgs.Raw(),
}).Build()

const DATA_SOURCE = "test.bolt"

const DEFAULT_BLOCK_HEIGHT = primitives.BlockHeight(1974)

var DEFAULT_TIME = primitives.TimestampNano(time.Now().UnixNano())

func removeDB() {
	os.RemoveAll(DATA_SOURCE)
}

func TestArs(t *testing.T) {
	t.Skip("does not work and does not plan to work")
	arizonaArgsNative := []interface{}{
		"Raising Arizona", uint32(1987), "Nicolas Cage",
	}

	packedArgs, err := protocol.PackedInputArgumentsFromNatives(arizonaArgsNative)
	require.NoError(t, err)

	args, err := protocol.PackedOutputArgumentsToNatives(packedArgs)
	require.NoError(t, err)

	require.EqualValues(t, arizonaArgsNative, args)
}

func TestArguments(t *testing.T) {
	argArray, err := protocol.ArgumentArrayFromNatives(arizonaArgsNative)
	require.NoError(t, err)
	t.Log(argArray.String())

	arr := protocol.ArgumentArrayReader(argArray.Raw())
	for i := arr.ArgumentsIterator(); i.HasNext(); {
		t.Log(i.NextArguments().String())
	}

	args, err := protocol.PackedOutputArgumentsToNatives(argArray.Raw())
	require.NoError(t, err)

	require.EqualValues(t, arizonaArgsNative, args)
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

	args, _ := protocol.PackedOutputArgumentsToNatives(eventList[0].Arguments())
	require.EqualValues(t, arizonaArgsNative, args)

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

func TestStorage_StoreEventUpdatesBlockHeight(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE, false)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*protocol.IndexedEvent{})
	require.NoError(t, err, "could not store event")

	blockHeight := storage.GetBlockHeight()
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT, blockHeight)

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT+100, DEFAULT_TIME, []*protocol.IndexedEvent{})
	require.NoError(t, err, "could not store event")

	updatedBlockHeight := storage.GetBlockHeight()
	require.NoError(t, err)
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT+100, updatedBlockHeight)
}
