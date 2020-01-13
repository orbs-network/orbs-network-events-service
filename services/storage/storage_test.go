package storage

import (
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var arizonaArgsNative = []interface{}{
	"Raising Arizona", uint32(1987), "Nicolas Cage",
}

var vampireArgsNative = []interface{}{
	"Vampire's Kiss", uint32(1989), "Nicolas Cage",
}

var mashupArgsNative = []interface{}{
	"Vampire's Kiss", uint32(1987), "Nicolas Cage",
}

var arizonaArgs, _ = protocol.ArgumentArrayFromNatives(arizonaArgsNative)
var vampireArgs, _ = protocol.ArgumentArrayFromNatives(vampireArgsNative)
var mashupArgs, _ = protocol.ArgumentArrayFromNatives(mashupArgsNative)

var ARIZONA_EVENT = (&types.IndexedEventBuilder{
	ContractName:    "SomeContract",
	EventName:       "MovieRelease",
	BlockHeight:     DEFAULT_BLOCK_HEIGHT,
	ExecutionResult: types.EXECUTION_RESULT_SUCCESS,
	Timestamp:       DEFAULT_TIME,
	Arguments:       arizonaArgs.Raw(),
}).Build()

var VAMPIRE_EVENT = (&types.IndexedEventBuilder{
	ContractName:    "SomeContract",
	EventName:       "MovieRelease",
	BlockHeight:     DEFAULT_BLOCK_HEIGHT + 100,
	ExecutionResult: types.EXECUTION_RESULT_SUCCESS,
	Timestamp:       DEFAULT_TIME + 5000,
	Arguments:       vampireArgs.Raw(),
}).Build()

var MASHUP_EVENT = (&types.IndexedEventBuilder{
	ContractName:    "SomeContract",
	EventName:       "MovieRelease",
	BlockHeight:     DEFAULT_BLOCK_HEIGHT + 100,
	ExecutionResult: types.EXECUTION_RESULT_SUCCESS,
	Timestamp:       DEFAULT_TIME + 5000,
	Arguments:       mashupArgs.Raw(),
}).Build()

const DATA_SOURCE = "test.bolt"

const DEFAULT_BLOCK_HEIGHT = 1974

var DEFAULT_TIME = uint64(time.Now().UnixNano())

func removeDB() {
	os.RemoveAll(DATA_SOURCE)
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

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*types.IndexedEvent{ARIZONA_EVENT})
	require.NoError(t, err, "could not store event")

	blockHeight := storage.GetBlockHeight()
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT, blockHeight)

	eventList, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
	}).Build())
	require.NoError(t, err)

	require.Len(t, eventList, 1)
	require.EqualValues(t, ARIZONA_EVENT.Raw(), eventList[0].Raw())

	args, _ := protocol.PackedOutputArgumentsToNatives(eventList[0].Arguments())
	require.EqualValues(t, arizonaArgsNative, args)

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT+100, DEFAULT_TIME, []*types.IndexedEvent{VAMPIRE_EVENT})
	require.NoError(t, err, "could not store event")

	updatedBlockHeight := storage.GetBlockHeight()
	require.NoError(t, err)
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT+100, updatedBlockHeight)

	eventList, err = storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{VAMPIRE_EVENT.EventName()},
	}).Build())
	require.NoError(t, err)
	require.Len(t, eventList, 2)
}

func TestStorage_StoreEventUpdatesBlockHeight(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE, false)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*types.IndexedEvent{})
	require.NoError(t, err, "could not store event")

	blockHeight := storage.GetBlockHeight()
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT, blockHeight)

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT+100, DEFAULT_TIME, []*types.IndexedEvent{})
	require.NoError(t, err, "could not store event")

	updatedBlockHeight := storage.GetBlockHeight()
	require.NoError(t, err)
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT+100, updatedBlockHeight)
}

func TestStorage_GetEventsForNonExistentContract(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE, false)
	require.NoError(t, err, "could not create new data source")

	events, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: "Unbelievable",
	}).Build())
	require.NoError(t, err)
	require.Len(t, events, 0)
}

func TestStorage_GetEventsFilterByBlockHeight(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE, false)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(0, 0, []*types.IndexedEvent{})
	require.NoError(t, err, "could not store event")

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*types.IndexedEvent{ARIZONA_EVENT})
	require.NoError(t, err, "could not store event")

	emptyEvents, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
		FromBlock:    0,
		ToBlock:      DEFAULT_BLOCK_HEIGHT - 100,
	}).Build())
	require.NoError(t, err)
	require.Len(t, emptyEvents, 0)

	events, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
		FromBlock:    DEFAULT_BLOCK_HEIGHT,
		ToBlock:      DEFAULT_BLOCK_HEIGHT + 100,
	}).Build())
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.EqualValues(t, events[0].Raw(), ARIZONA_EVENT.Raw())
}

func TestStorage_GetEventsFilterByBlockHeightWithRange(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE, false)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*types.IndexedEvent{ARIZONA_EVENT})
	require.NoError(t, err, "could not store event")

	emptyEvents, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
		FromBlock:    1,
		ToBlock:      DEFAULT_BLOCK_HEIGHT - 100,
	}).Build())
	require.NoError(t, err)
	require.Len(t, emptyEvents, 0)

	events, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
		FromBlock:    1,
		ToBlock:      DEFAULT_BLOCK_HEIGHT + 9000,
	}).Build())
	require.NoError(t, err)
	println(len(events))
	require.Len(t, events, 1)
	require.EqualValues(t, events[0].Raw(), ARIZONA_EVENT.Raw())
}

func TestStorage_GetEventsFilterByTimestamp(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE, false)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(0, 0, []*types.IndexedEvent{})
	require.NoError(t, err, "could not store event")

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*types.IndexedEvent{ARIZONA_EVENT})
	require.NoError(t, err, "could not store event")

	emptyEvents, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
		FromTime:     0,
		ToTime:       DEFAULT_TIME - 100,
	}).Build())
	require.NoError(t, err)
	require.Len(t, emptyEvents, 0)

	events, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
		FromTime:     DEFAULT_TIME - 100,
		ToTime:       DEFAULT_TIME + 100,
	}).Build())
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.EqualValues(t, events[0].Raw(), ARIZONA_EVENT.Raw())
}

func TestStorage_GetEventsFilterWithFieldMatching(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE, false)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*types.IndexedEvent{ARIZONA_EVENT})
	require.NoError(t, err, "could not store event")

	vampireFilter, _ := protocol.ArgumentArrayFromNatives([]interface{}{"Vampire's Kill"})

	emptyEvents, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
		Filters: [][]byte{
			vampireFilter.Raw(),
		},
	}).Build())
	require.NoError(t, err)
	require.Len(t, emptyEvents, 0)

	arizonaFilter, _ := protocol.ArgumentArrayFromNatives([]interface{}{"Raising Arizona"})

	events, err := storage.GetEvents((&types.IndexerRequestBuilder{
		ContractName: ARIZONA_EVENT.ContractName(),
		EventName:    []string{ARIZONA_EVENT.EventName()},
		Filters: [][]byte{
			arizonaFilter.Raw(),
		},
	}).Build())
	require.NoError(t, err)
	require.Len(t, events, 1)
	require.EqualValues(t, events[0].Raw(), ARIZONA_EVENT.Raw())
}
