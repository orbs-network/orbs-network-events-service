package main

import (
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/events"
	"github.com/orbs-network/scribe/log"
	"os"
	"time"
)

const MAX = uint(100000)

func main() {
	logger := config.GetLogger()

	code, err := orbs.ReadSourcesFromDir("./test/e2e/_contracts")
	if err != nil {
		logger.Error("could not read the contract")
		os.Exit(1)
	}

	orbsClient := orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)
	account, _ := orbs.CreateAccount()

	contractName := fmt.Sprintf("EventEmitter%d", time.Now().UnixNano())
	deployTx, _, _ := orbsClient.CreateDeployTransaction(account.PublicKey, account.PrivateKey, contractName, orbs.PROCESSOR_TYPE_NATIVE, code...)
	sendTxAndWaitForSuccess(logger, orbsClient, deployTx)
	logger.Info("deployed contract", log.String("contractName", contractName))

	go func() {
		for i := uint(0); i < MAX; i++ {
			arizonaTx, _, _ := orbsClient.CreateTransaction(account.PublicKey, account.PrivateKey, contractName, "release",
				"Raising Arizona", uint32(1987), "Nicolas Cage")
			sendTxAndWaitForSuccess(logger, orbsClient, arizonaTx)

			if i%1000 == 0 {
				logger.Info("batch sent", log.Uint("txSent", i), log.Uint("txTotal", MAX))
			}
		}
	}()

	go func() {
		for {
			time.Sleep(15 * time.Second)
			blockHeight, err := events.GetBlockHeight(orbsClient, account)
			if err == nil {
				logger.Info("blockHeight advanced", log.Uint64("blockHeight", uint64(blockHeight)))
			}

		}
	}()

	var forever chan interface{}
	<-forever
}

func sendTxAndWaitForSuccess(logger log.Logger, orbsClient *orbs.OrbsClient, tx []byte) {
	res, err := orbsClient.SendTransaction(tx)
	if err != nil {
		logger.Error("sendTx failed", log.Error(err))
	}

	if res.ExecutionResult != codec.EXECUTION_RESULT_SUCCESS {
		logger.Error("smart contract error", log.Stringable("executionResult", res.ExecutionResult))
	}
}
