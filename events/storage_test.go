package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

var DEFAULT_EVENT = &codec.Event{
	ContractName: "SomeContract",
	EventName:    "MovieRelease",
	Arguments: []interface{}{
		"Raising Arizona", "1987-03-06", "Nicolas Cage",
	},
}

const DATA_SOURCE = "test.sqlite3"
const DATA_SOURCE_MODE = "?mode=rw"

const DEFAULT_BLOCK_HEIGHT = uint64(1974)

var DEFAULT_TIME = time.Now().UnixNano()

func removeDB() {
	os.RemoveAll(DATA_SOURCE)
}

func TestStoreEvent(t *testing.T) {
	removeDB()

	storage, err := NewStorage(DATA_SOURCE + DATA_SOURCE_MODE)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvent(DEFAULT_BLOCK_HEIGHT, DEFAULT_TIME, DEFAULT_EVENT)
	require.NoError(t, err, "could not store event")

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
}
