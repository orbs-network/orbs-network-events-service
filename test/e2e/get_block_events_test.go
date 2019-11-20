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

func TestGetBlockEvents(t *testing.T) {
	client := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)

	account, _ := orbs.CreateAccount()
	contractName, _ := deployEventEmitterContract(t, client, account)

	arizonaTx, _, _ := client.CreateTransaction(account.PublicKey, account.PrivateKey, contractName, "release",
		"Raising Arizona", uint32(1987), "Nicolas Cage")

	res, err := client.SendTransaction(arizonaTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, res.ExecutionResult)

	_, eventList, err := events.GetBlockEvents(client, primitives.BlockHeight(res.BlockHeight))
	require.NoError(t, err)

	require.Len(t, eventList, 1)
	args, _ := protocol.ArgumentArrayFromNatives([]interface{}{
		"Raising Arizona", uint32(1987), "Nicolas Cage",
	})
	require.EqualValues(t, (&protocol.IndexedEventBuilder{
		ContractName:    primitives.ContractName(contractName),
		EventName:       "MovieRelease",
		BlockHeight:     primitives.BlockHeight(res.BlockHeight),
		ExecutionResult: protocol.EXECUTION_RESULT_SUCCESS,
		Timestamp:       primitives.TimestampNano(res.BlockTimestamp.UnixNano()),
		Arguments:       args.Raw(),
	}).Build(), eventList[0])
}

func deployEventEmitterContract(t *testing.T, client *orbs.OrbsClient, account *orbs.OrbsAccount) (contractName string, blockHeight primitives.BlockHeight) {
	code, err := orbs.ReadSourcesFromDir("./_contracts")
	require.NoError(t, err)
	contractName = fmt.Sprintf("EventEmitter%d", time.Now().UnixNano())

	deployTx, _, _ := client.CreateDeployTransaction(account.PublicKey, account.PrivateKey, contractName, orbs.PROCESSOR_TYPE_NATIVE, code...)
	res, err := client.SendTransaction(deployTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, res.ExecutionResult)

	return contractName, primitives.BlockHeight(res.BlockHeight)
}
