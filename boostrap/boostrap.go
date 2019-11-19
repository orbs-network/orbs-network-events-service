package boostrap

import (
	"context"
	"github.com/orbs-network/orbs-network-events-service/boostrap/httpserver"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/services/background"
	indexer2 "github.com/orbs-network/orbs-network-events-service/services/indexer"
	"github.com/orbs-network/scribe/log"
)

func NewNode(ctx context.Context, cfg *config.Config, logger log.Logger) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	backgroundIndexer := background.NewBackgroundIndexer(cfg, logger)
	if err := backgroundIndexer.Start(ctx); err != nil {
		logger.Error("failed to start background service", log.Error(err))
		cancel()
	}

	api, err := indexer2.NewIndexer(cfg, logger)
	if err != nil {
		logger.Error("failed to start api service", log.Error(err))
		cancel()
	}

	httpserver.NewHttpServer(ctx, httpserver.NewServerConfig("0.0.0.0:9201", false), logger, api)
}
