package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var DEFAULT_EVENT = &codec.Event{
	ContractName: "SomeContract",
	EventName:    "MovieRelease",
	Arguments: []interface{}{
		"Raising Arizona", uint32(1987), "Nicolas Cage",
	},
}

const DATA_SOURCE = "test.bolt"

const DEFAULT_BLOCK_HEIGHT = uint64(1974)

var DEFAULT_TIME = time.Now().UnixNano()

func removeDB() {
	os.RemoveAll(DATA_SOURCE)
}

func TestStorage_StoreEvent(t *testing.T) {
	removeDB()

	storage, err := NewStorage(config.GetLogger(), DATA_SOURCE)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, []*codec.Event{DEFAULT_EVENT})
	require.NoError(t, err, "could not store event")

	blockHeight := storage.GetBlockHeight()
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT, blockHeight)

	eventList, err := storage.GetEvents(&FilterQuery{
		ContractName: DEFAULT_EVENT.ContractName,
		EventNames:   []string{DEFAULT_EVENT.EventName},
	})
	require.NoError(t, err)
	require.Len(t, eventList, 1)
	require.EqualValues(t, &StoredEvent{
		ContractName: DEFAULT_EVENT.ContractName,
		EventName:    DEFAULT_EVENT.EventName,
		BlockHeight:  DEFAULT_BLOCK_HEIGHT,
		Timestamp:    DEFAULT_TIME,
		Arguments:    DEFAULT_EVENT.Arguments,
	}, eventList[0])

	err = storage.StoreEvents(DEFAULT_BLOCK_HEIGHT+100, DEFAULT_TIME+5000, []*codec.Event{DEFAULT_EVENT})
	require.NoError(t, err, "could not store event")

	updatedBlockHeight := storage.GetBlockHeight()
	require.NoError(t, err)
	require.EqualValues(t, DEFAULT_BLOCK_HEIGHT+100, updatedBlockHeight)

	eventList, err = storage.GetEvents(&FilterQuery{
		ContractName: DEFAULT_EVENT.ContractName,
		EventNames:   []string{DEFAULT_EVENT.EventName},
	})
	require.NoError(t, err)
	require.Len(t, eventList, 2)
}
