package e2e

import (
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/events"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestProcessEvents(t *testing.T) {
	client := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)

	code, err := orbs.ReadSourcesFromDir("./_contracts")
	require.NoError(t, err)
	contractName := fmt.Sprintf("EventEmitter%d", time.Now().UnixNano())
	account, _ := orbs.CreateAccount()

	deployTx, _, _ := client.CreateDeployTransaction(account.PublicKey, account.PrivateKey, contractName, orbs.PROCESSOR_TYPE_NATIVE, code...)
	res, err := client.SendTransaction(deployTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, res.ExecutionResult)

	startingBlock := res.BlockHeight

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
	lastProcessedBlock, err := events.ProcessEvents(client, startingBlock, blockHeight, func(blockHeight uint64, timestamp int64, list []*protocol.IndexedEvent) error {
		eventList = append(eventList, list...)
		return nil
	})
	require.NoError(t, err)
	require.EqualValues(t, blockHeight, lastProcessedBlock)

	require.Len(t, eventList, 2)

	arizonaArgs, _ := protocol.PackedInputArgumentsFromNatives([]interface{}{
		"Raising Arizona", uint32(1987), "Nicolas Cage",
	})
	require.EqualValues(t, (&protocol.IndexedEventBuilder{
		ContractName:    primitives.ContractName(contractName),
		EventName:       "MovieRelease",
		BlockHeight:     primitives.BlockHeight(arizonaRes.BlockHeight),
		ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
		Timestamp:       primitives.TimestampNano(arizonaRes.BlockTimestamp.UnixNano()),
		Arguments:       arizonaArgs,
	}).Build(), eventList[0])

	vampireArgs, _ := protocol.PackedInputArgumentsFromNatives([]interface{}{
		"Vampire's Kiss", uint32(1989), "Nicolas Cage",
	})

	require.EqualValues(t, (&protocol.IndexedEventBuilder{
		ContractName:    primitives.ContractName(contractName),
		EventName:       "MovieRelease",
		BlockHeight:     primitives.BlockHeight(vampireRes.BlockHeight),
		ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
		Timestamp:       primitives.TimestampNano(vampireRes.BlockTimestamp.UnixNano()),
		Arguments:       vampireArgs,
	}).Build(), eventList[1])
}
