package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/pkg/errors"
)

func GetBlockEvents(client *orbs.OrbsClient, height primitives.BlockHeight) (timestamp primitives.TimestampNano, events []*protocol.IndexedEvent, err error) {
	res, err := client.GetBlock(uint64(height))
	if err != nil {
		return
	}

	timestamp = primitives.TimestampNano(res.BlockTimestamp.UnixNano())
	for _, tx := range res.Transactions {
		for i, event := range tx.OutputEvents {
			arguments, err := protocol.PackedInputArgumentsFromNatives(event.Arguments)
			if err != nil {
				return 0, nil, err
			}
			executionResult, err := decodeExecutionResult(tx.ExecutionResult)
			if err != nil {
				return 0, nil, err
			}

			indexedEvent := (&protocol.IndexedEventBuilder{
				ContractName:    primitives.ContractName(event.ContractName),
				EventName:       event.EventName,
				BlockHeight:     primitives.BlockHeight(height),
				Timestamp:       timestamp, // tx time is irrelevant
				Index:           uint32(i),
				ExecutionResult: executionResult,
				Arguments:       arguments,
			}).Build()

			events = append(events, indexedEvent)
		}
	}

	return
}

func decodeExecutionResult(input codec.ExecutionResult) (protocol.ExecutionResult, error) {
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
