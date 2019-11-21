package e2e

import (
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/stretchr/testify/require"
	"os"
	"time"
)

func removeDB() {
	os.RemoveAll("./vchain-42.bolt")
}

func deployEventEmitterContract(t require.TestingT, client *orbs.OrbsClient, account *orbs.OrbsAccount) (contractName string, blockHeight uint64) {
	code, err := orbs.ReadSourcesFromDir("./_contracts")
	require.NoError(t, err)
	contractName = fmt.Sprintf("EventEmitter%d", time.Now().UnixNano())

	deployTx, _, _ := client.CreateDeployTransaction(account.PublicKey, account.PrivateKey, contractName, orbs.PROCESSOR_TYPE_NATIVE, code...)
	res, err := client.SendTransaction(deployTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, res.ExecutionResult)

	return contractName, uint64(res.BlockHeight)
}

func sendArizonaTransaction(t require.TestingT, orbsClient *orbs.OrbsClient, account *orbs.OrbsAccount, contractName string) *codec.SendTransactionResponse {
	arizonaTx, _, _ := orbsClient.CreateTransaction(account.PublicKey, account.PrivateKey, contractName, "release",
		"Raising Arizona", uint32(1987), "Nicolas Cage")
	arizonaRes, err := orbsClient.SendTransaction(arizonaTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, arizonaRes.ExecutionResult)

	return arizonaRes
}

func sendVampireTransaction(t require.TestingT, orbsClient *orbs.OrbsClient, account *orbs.OrbsAccount, contractName string) *codec.SendTransactionResponse {
	vampireTx, _, _ := orbsClient.CreateTransaction(account.PublicKey, account.PrivateKey, contractName, "release",
		"Vampire's Kiss", uint32(1989), "Nicolas Cage")

	vampireRes, err := orbsClient.SendTransaction(vampireTx)
	require.NoError(t, err)
	require.EqualValues(t, codec.EXECUTION_RESULT_SUCCESS, vampireRes.ExecutionResult)

	return vampireRes
}
