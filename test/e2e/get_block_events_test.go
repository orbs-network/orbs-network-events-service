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

func TestGetBlockEvents(t *testing.T) {
	client := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)

	account, _ := orbs.CreateAccount()
	contractName, _ := deployEventEmitterContract(t, client, account)

	res := sendArizonaTransaction(t, client, account, contractName)

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
