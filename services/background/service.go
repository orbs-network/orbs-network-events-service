package background

import (
	"context"
	"fmt"
	"github.com/orbs-network/govnr"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/events"
	"github.com/orbs-network/orbs-network-events-service/services/storage"
	"github.com/orbs-network/scribe/log"
	"time"
)

type BackgroundIndexer interface {
	Start(ctx context.Context) error
}

type service struct {
	cfg    *config.Config
	logger log.Logger
	db     storage.Storage

	vcid       uint32
	supervisor govnr.ShutdownWaiter
}

func NewBackgroundIndexer(cfg *config.Config, logger log.Logger, db storage.Storage, vcid uint32) BackgroundIndexer {
	return &service{
		cfg:    cfg,
		logger: logger.WithTags(log.Service("indexer")),
		vcid:   vcid,
		db:     db,
	}
}

func (s *service) Start(ctx context.Context) error {
	_, err := s.indexVchain(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) indexVchain(ctx context.Context) (govnr.ShutdownWaiter, error) {
	serviceLogger := s.logger.WithTags(log.Uint32("vcid", s.vcid))

	handle := govnr.Forever(ctx, fmt.Sprintf("vchain %d handler", s.vcid), config.NewErrorHandler(serviceLogger), func() {
		client := orbs.NewClient(s.cfg.Endpoint, s.vcid, codec.NETWORK_TYPE_TEST_NET)
		account, _ := orbs.CreateAccount()

		lastProcessedBlock := s.db.GetBlockHeight()
		serviceLogger.Info("starting the sync process", log.Uint64("blockHeight", uint64(lastProcessedBlock)))

		for {
			time.Sleep(s.cfg.PollingInterval)

			finalBlock, err := events.GetBlockHeight(client, account)
			if err != nil {
				serviceLogger.Error("failed to get last block height", log.Error(err))
				return
			}
			if finalBlock <= lastProcessedBlock+1 {
				continue
			}

			lastProcessedBlock, err = events.ProcessEvents(client, lastProcessedBlock+1, finalBlock, s.db.StoreEvents)
			if err != nil {
				serviceLogger.Error("failed to store events", log.Error(err))
				return
			}
		}
	})

	supervisor := &govnr.TreeSupervisor{}
	supervisor.Supervise(handle)
	return supervisor, nil
}
