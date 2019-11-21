package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/types"
)

type EventProcessingCallback func(blockHeight uint64, timestamp uint64, eventList []*types.IndexedEvent) error

func ProcessEvents(client *orbs.OrbsClient, start uint64, end uint64, callback EventProcessingCallback) (uint64, error) {
	for blockHeight := start; blockHeight <= end; blockHeight++ {
		timestamp, events, err := GetBlockEvents(client, blockHeight)
		if err != nil {
			return blockHeight, err
		}

		if err := callback(blockHeight, timestamp, events); err != nil {
			return blockHeight, err
		}
	}

	return end, nil
}
