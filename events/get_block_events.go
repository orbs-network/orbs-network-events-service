package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/pkg/errors"
)

func GetBlockEvents(client *orbs.OrbsClient, height uint64) (timestamp uint64, events []*types.IndexedEvent, err error) {
	res, err := client.GetBlock(uint64(height))
	if err != nil {
		return
	}

	timestamp = uint64(res.BlockTimestamp.UnixNano())
	for _, tx := range res.Transactions {
		for i, event := range tx.OutputEvents {
			arguments, err := protocol.ArgumentArrayFromNatives(event.Arguments)
			if err != nil {
				return 0, nil, err
			}
			executionResult, err := decodeExecutionResult(tx.ExecutionResult)
			if err != nil {
				return 0, nil, err
			}

			indexedEvent := (&types.IndexedEventBuilder{
				ContractName:    event.ContractName,
				EventName:       event.EventName,
				BlockHeight:     uint64(height),
				Timestamp:       timestamp, // tx time is irrelevant
				Index:           uint32(i),
				ExecutionResult: executionResult,
				Arguments:       arguments.Raw(),
			}).Build()

			events = append(events, indexedEvent)
		}
	}

	return
}

func decodeExecutionResult(input codec.ExecutionResult) (types.ExecutionResult, error) {
	switch input {
	case codec.EXECUTION_RESULT_SUCCESS:
		return 1, nil
	case codec.EXECUTION_RESULT_ERROR_SMART_CONTRACT:
		return 2, nil
	case codec.EXECUTION_RESULT_ERROR_INPUT:
		return 3, nil
	case codec.EXECUTION_RESULT_ERROR_CONTRACT_NOT_DEPLOYED:
		return 4, nil
	case codec.EXECUTION_RESULT_ERROR_UNEXPECTED:
		return 5, nil
	case codec.EXECUTION_RESULT_NOT_EXECUTED:
		return 6, nil
	}

	return 0, errors.Errorf("could not decode execution result from value %s", input)
}
