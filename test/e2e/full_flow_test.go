package e2e

import (
	"context"
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/bootstrap"
	"github.com/orbs-network/orbs-network-events-service/client"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestFullFlow(t *testing.T) {
	removeDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	server, err := bootstrap.NewCluster(ctx, &config.Config{
		Endpoint:        "http://localhost:8080",
		VirtualChains:   []uint32{42},
		DB:              "./",
		PollingInterval: 10 * time.Millisecond,
	}, config.GetLogger())
	require.NoError(t, err)

	orbsClient := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)

	account, _ := orbs.CreateAccount()
	contractName, _ := deployEventEmitterContract(t, orbsClient, account)

	sendArizonaTransaction(t, orbsClient, account, contractName)
	sendVampireTransaction(t, orbsClient, account, contractName)

	require.Eventually(t, func() bool {
		events, err := client.GetEvents(fmt.Sprintf("http://localhost:%d", server.Port()), client.GetEventsQuery{
			VirtualChainId: 42,
			ContractName:   contractName,
			EventName:      []string{"MovieRelease"},
			Filters: [][]interface{}{
				{"Vampire's Kiss"}, {uint32(1989)},
			},
			FromBlock: 0,
			ToBlock:   9999,
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

		println("got", arguments[0].(string))
		return arguments[0].(string) == "Vampire's Kiss"
	}, 10*time.Second, 100*time.Millisecond, "indexer api should return events")

}
