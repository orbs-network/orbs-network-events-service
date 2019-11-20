package e2e

import (
	"context"
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/boostrap"
	"github.com/orbs-network/orbs-network-events-service/client"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func removeDB() {
	os.RemoveAll("./vchain-42.bolt")
}

func TestFullFlow(t *testing.T) {
	removeDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := boostrap.NewCluster(ctx, &config.Config{
		Endpoint:        "http://localhost:8080",
		VirtualChains:   []uint32{42},
		DB:              "./",
		PollingInterval: 10 * time.Millisecond,
	}, config.GetLogger())
	require.NoError(t, err)

	orbsClient := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)

	account, _ := orbs.CreateAccount()
	contractName, _ := deployEventEmitterContract(t, orbsClient, account)

	arizonaTx, _, _ := orbsClient.CreateTransaction(account.PublicKey, account.PrivateKey, contractName, "release",
		"Raising Arizona", uint32(1987), "Nicolas Cage")

	arizonaRes, err := orbsClient.SendTransaction(arizonaTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, arizonaRes.ExecutionResult)

	require.Eventually(t, func() bool {
		events, err := client.GetEvents(fmt.Sprintf("http://localhost:%d", server.Port()), client.GetEventsQuery{
			VirtualChainId: 42,
			ContractName:   contractName,
			EventName:      []string{"MovieRelease"},
		})

		if err != nil {
			t.Log(err)
			return false
		}

		if len(events) == 0 {
			return false
		}

		arguments, err := protocol.PackedOutputArgumentsToNatives(events[0].RawArguments())
		if err != nil {
			return false
		}

		return arguments[2].(string) == "Nicolas Cage"
	}, 10*time.Second, 100*time.Millisecond, "indexer api should return events")

}
