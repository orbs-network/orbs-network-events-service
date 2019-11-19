package main

import (
	"context"
	"flag"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/services/indexer"
	"github.com/orbs-network/scribe/log"
	"io/ioutil"
)

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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	indexer := indexer.NewIndexer(cfg, logger)
	if err := indexer.Start(ctx); err != nil {
		logger.Error("failed to start indexer service", log.Error(err))
		cancel()
	}
}
