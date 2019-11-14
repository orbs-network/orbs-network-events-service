package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
)

type EventProcessingCallback func(eventMap []*codec.Event) error

func ProcessEvents(client *orbs.OrbsClient, start uint64, end uint64, callback EventProcessingCallback) (uint64, error) {
	for i := start; i <= end; i++ {
		events, err := GetBlockEvents(client, i)
		if err != nil {
			return i, err
		}

		if err := callback(events); err != nil {
			return i, err
		}
	}

	return end, nil
}
