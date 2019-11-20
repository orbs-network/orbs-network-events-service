package e2e

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/events"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessEvents(t *testing.T) {
	client := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)

	account, _ := orbs.CreateAccount()
	contractName, startingBlock := deployEventEmitterContract(t, client, account)

	arizonaTx, _, _ := client.CreateTransaction(account.PublicKey, account.PrivateKey, contractName, "release",
		"Raising Arizona", uint32(1987), "Nicolas Cage")

	arizonaRes, err := client.SendTransaction(arizonaTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, arizonaRes.ExecutionResult)

	vampireTx, _, _ := client.CreateTransaction(account.PublicKey, account.PrivateKey, contractName, "release",
		"Vampire's Kiss", uint32(1989), "Nicolas Cage")

	vampireRes, err := client.SendTransaction(vampireTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, vampireRes.ExecutionResult)

	blockHeight, err := events.GetBlockHeight(client, account)
	require.NoError(t, err)

	var eventList []*protocol.IndexedEvent
	lastProcessedBlock, err := events.ProcessEvents(client, startingBlock, blockHeight, func(blockHeight primitives.BlockHeight, timestamp primitives.TimestampNano, list []*protocol.IndexedEvent) error {
		eventList = append(eventList, list...)
		return nil
	})
	require.NoError(t, err)
	require.EqualValues(t, blockHeight, lastProcessedBlock)

	require.Len(t, eventList, 2)

	arizonaArgs, _ := protocol.ArgumentArrayFromNatives([]interface{}{
		"Raising Arizona", uint32(1987), "Nicolas Cage",
	})
	require.EqualValues(t, (&protocol.IndexedEventBuilder{
		ContractName:    primitives.ContractName(contractName),
		EventName:       "MovieRelease",
		BlockHeight:     primitives.BlockHeight(arizonaRes.BlockHeight),
		ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
		Timestamp:       primitives.TimestampNano(arizonaRes.BlockTimestamp.UnixNano()),
		Arguments:       arizonaArgs.Raw(),
	}).Build(), eventList[0])

	vampireArgs, _ := protocol.ArgumentArrayFromNatives([]interface{}{
		"Vampire's Kiss", uint32(1989), "Nicolas Cage",
	})

	require.EqualValues(t, (&protocol.IndexedEventBuilder{
		ContractName:    primitives.ContractName(contractName),
		EventName:       "MovieRelease",
		BlockHeight:     primitives.BlockHeight(vampireRes.BlockHeight),
		ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
		Timestamp:       primitives.TimestampNano(vampireRes.BlockTimestamp.UnixNano()),
		Arguments:       vampireArgs.Raw(),
	}).Build(), eventList[1])
}
