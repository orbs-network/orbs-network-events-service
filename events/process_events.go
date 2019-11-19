package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

type EventProcessingCallback func(blockHeight primitives.BlockHeight, timestamp primitives.TimestampNano, eventList []*protocol.IndexedEvent) error

func ProcessEvents(client *orbs.OrbsClient, start primitives.BlockHeight, end primitives.BlockHeight, callback EventProcessingCallback) (primitives.BlockHeight, error) {
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
