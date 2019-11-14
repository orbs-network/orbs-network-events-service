package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
)

type EventMap map[string][]*codec.Event

func GetBlockEvents(client *orbs.OrbsClient, height uint64) (EventMap, error) {
	res, err := client.GetBlock(height)
	if err != nil {
		return nil, err
	}

	events := make(EventMap)
	for _, tx := range res.Transactions {
		es := events[tx.ContractName]
		es = append(es, tx.OutputEvents...)
		events[tx.ContractName] = es
	}

	return events, nil
}