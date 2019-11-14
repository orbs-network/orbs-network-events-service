package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
)

func GetBlockEvents(client *orbs.OrbsClient, height uint64) (events []*codec.Event, err error) {
	res, err := client.GetBlock(height)
	if err != nil {
		return
	}

	for _, tx := range res.Transactions {
		events = append(events, tx.OutputEvents...)
	}

	return
}
