package bootstrap

import (
	"context"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/services/background"
	"github.com/orbs-network/orbs-network-events-service/services/indexer"
	"github.com/orbs-network/orbs-network-events-service/services/storage"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/scribe/log"
)

type Node struct {
	vcid uint32

	db  storage.Storage
	bg  background.BackgroundIndexer
	api types.Indexer
}

func NewNode(ctx context.Context, cfg *config.Config, logger log.Logger, vcid uint32) (*Node, error) {
	nodeLogger := logger.WithTags(log.Uint32("vcid", vcid))

	db, err := storage.NewStorageForChain(nodeLogger, cfg.DB, vcid, false)
	if err != nil {
		return nil, err
	}

	bg := background.NewBackgroundIndexer(cfg, nodeLogger, db, vcid)
	if err := bg.Start(ctx); err != nil {
		return nil, err
	}

	api, err := indexer.NewIndexer(cfg, logger, db)
	if err != nil {
		return nil, err
	}

	return &Node{
		vcid: vcid,
		db:   db,
		bg:   bg,
		api:  api,
	}, nil
}
