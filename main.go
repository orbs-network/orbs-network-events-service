package main

import (
	"context"
	"flag"
	"github.com/orbs-network/orbs-network-events-service/boostrap"
	"github.com/orbs-network/orbs-network-events-service/config"
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

	boostrap.NewNode(ctx, cfg, logger)
}
