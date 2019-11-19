package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/orbs-network/govnr"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/events"
	"github.com/orbs-network/scribe/log"
	"io/ioutil"
	"time"
)

type stdoutErrorer struct{}

func (s *stdoutErrorer) Error(err error) {
	println(err.Error())
}

func main() {
	logger := config.GetLogger()
	logger.Info("Starting signer service")

	configPath := flag.String("config", "./config.json", "path to config")
	flag.Parse()

	configData, err := ioutil.ReadFile(*configPath)
	if err != nil {
		panic(configData)
	}

	cfg, err := config.Parse(configData)
	if err != nil {
		panic(err)
	}

	errorHandler := &stdoutErrorer{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, chainId := range cfg.VirtualChains {
		handle := govnr.Forever(ctx, fmt.Sprintf("vchain %d handler", chainId), errorHandler, func() {
			client := orbs.NewClient(cfg.Endpoint, chainId, codec.NETWORK_TYPE_TEST_NET)
			account, _ := orbs.CreateAccount()
			vchainLogger := logger.WithTags(log.Uint32("vcid", chainId))
			storage, err := events.NewStorage(vchainLogger, fmt.Sprintf("./data/vchain-%d.bolt", chainId))
			defer storage.Shutdown()

			if err != nil {
				logger.Error("failed to access storage", log.Error(err))
				return
			}

			lastProcessedBlock := storage.GetBlockHeight()
			logger.Info("starting the sync process", log.Uint64("blockHeight", lastProcessedBlock))

			for {
				time.Sleep(cfg.PollingInterval)

				finalBlock, err := events.GetBlockHeight(client, account)
				if err != nil {
					logger.Error("failed to get last block height", log.Error(err))
					return
				}
				if finalBlock <= lastProcessedBlock+1 {
					continue
				}

				lastProcessedBlock, err = events.ProcessEvents(client, lastProcessedBlock+1, finalBlock, storage.StoreEvents)
				if err != nil {
					logger.Error("failed to store events", log.Error(err))
					return
				}
			}
		})

		supervisor := &govnr.TreeSupervisor{}
		supervisor.Supervise(handle)

		supervisor.WaitUntilShutdown(ctx)
	}
}
