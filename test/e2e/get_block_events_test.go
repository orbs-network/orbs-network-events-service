package e2e

import (
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/events"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetBlockEvents(t *testing.T) {
	client := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)

	code, err := orbs.ReadSourcesFromDir("./_contracts")
	require.NoError(t, err)
	contractName := fmt.Sprintf("EventEmitter%d", time.Now().UnixNano())
	account, _ := orbs.CreateAccount()

	deployTx, _, _ := client.CreateDeployTransaction(account.PublicKey, account.PrivateKey, contractName, orbs.PROCESSOR_TYPE_NATIVE, code...)
	res, err := client.SendTransaction(deployTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, res.ExecutionResult)

	arizonaTx, _, _ := client.CreateTransaction(account.PublicKey, account.PrivateKey, contractName, "release",
		"Raising Arizona", uint32(1987), "Nicolas Cage")

	res, err = client.SendTransaction(arizonaTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, res.ExecutionResult)

	eventMap, err := events.GetBlockEvents(client, res.BlockHeight)
	require.NoError(t, err)

	require.Len(t, eventMap[contractName], 1)
	require.EqualValues(t, &codec.Event{
		ContractName: contractName,
		EventName:    "MovieRelease",
		Arguments: []interface{}{
			"Raising Arizona", uint32(1987), "Nicolas Cage",
		},
	}, eventMap[contractName][0])
}
