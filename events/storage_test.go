package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

var DEFAULT_EVENT = &codec.Event{
	ContractName: "SomeContract",
	EventName:    "MovieRelease",
	Arguments: []interface{}{
		"Raising Arizona", uint32(1987), "Nicolas Cage",
	},
}

const DATA_SOURCE = "test.sqlite3"
const DATA_SOURCE_MODE = "?mode=rw"

func removeDB() {
	os.RemoveAll(DATA_SOURCE)
}

func TestStoreEvent(t *testing.T) {
	removeDB()

	storage, err := NewStorage(DATA_SOURCE + DATA_SOURCE_MODE)
	require.NoError(t, err, "could not create new data source")

	err = storage.StoreEvent(DEFAULT_EVENT)
	require.NoError(t, err, "could not store event")

	eventList, err := storage.GetEvents(DEFAULT_EVENT.ContractName, DEFAULT_EVENT.EventName)
	require.NoError(t, err)
	require.Len(t, eventList, 1)
}
