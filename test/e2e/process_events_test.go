package e2e

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/events"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessEvents(t *testing.T) {
	client := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)

	account, _ := orbs.CreateAccount()
	contractName, startingBlock := deployEventEmitterContract(t, client, account)

	arizonaRes := sendArizonaTransaction(t, client, account, contractName)
	vampireRes := sendVampireTransaction(t, client, account, contractName)

	blockHeight, err := events.GetBlockHeight(client, account)
	require.NoError(t, err)

	var eventList []*types.IndexedEvent
	lastProcessedBlock, err := events.ProcessEvents(client, startingBlock, blockHeight, func(blockHeight uint64, timestamp uint64, list []*types.IndexedEvent) error {
		eventList = append(eventList, list...)
		return nil
	})
	require.NoError(t, err)
	require.EqualValues(t, blockHeight, lastProcessedBlock)

	require.Len(t, eventList, 2)

	arizonaArgs, _ := protocol.ArgumentArrayFromNatives([]interface{}{
		"Raising Arizona", uint32(1987), "Nicolas Cage",
	})
	require.EqualValues(t, (&types.IndexedEventBuilder{
		ContractName:    contractName,
		EventName:       "MovieRelease",
		BlockHeight:     uint64(arizonaRes.BlockHeight),
		ExecutionResult: types.EXECUTION_RESULT_SUCCESS,
		Timestamp:       uint64(arizonaRes.BlockTimestamp.UnixNano()),
		Arguments:       arizonaArgs.Raw(),
	}).Build(), eventList[0])

	vampireArgs, _ := protocol.ArgumentArrayFromNatives([]interface{}{
		"Vampire's Kiss", uint32(1989), "Nicolas Cage",
	})

	require.EqualValues(t, (&types.IndexedEventBuilder{
		ContractName:    contractName,
		EventName:       "MovieRelease",
		BlockHeight:     uint64(vampireRes.BlockHeight),
		ExecutionResult: types.EXECUTION_RESULT_SUCCESS,
		Timestamp:       uint64(vampireRes.BlockTimestamp.UnixNano()),
		Arguments:       vampireArgs.Raw(),
	}).Build(), eventList[1])
}
