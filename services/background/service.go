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

	supervisors []govnr.ShutdownWaiter
}

func NewBackgroundIndexer(cfg *config.Config, logger log.Logger) BackgroundIndexer {
	return &service{
		cfg:    cfg,
		logger: logger.WithTags(log.Service("indexer")),
	}
}

func (s *service) Start(ctx context.Context) error {
	for _, vcid := range s.cfg.VirtualChains {
		supervisor, err := s.indexVchain(ctx, vcid)
		if err != nil {
			return err
		}
		s.supervisors = append(s.supervisors, supervisor)
	}

	for _, supervisor := range s.supervisors {
		supervisor.WaitUntilShutdown(ctx)
	}

	return nil
}

func (s *service) indexVchain(ctx context.Context, vcid uint32) (govnr.ShutdownWaiter, error) {
	serviceLogger := s.logger.WithTags(log.Uint32("vcid", vcid))

	handle := govnr.Forever(ctx, fmt.Sprintf("vchain %d handler", vcid), config.NewErrorHandler(serviceLogger), func() {
		client := orbs.NewClient(s.cfg.Endpoint, vcid, codec.NETWORK_TYPE_TEST_NET)
		account, _ := orbs.CreateAccount()
		db, err := storage.NewStorageForChain(serviceLogger, vcid, false)
		defer db.Shutdown()

		if err != nil {
			serviceLogger.Error("failed to access storage", log.Error(err))
			return
		}

		lastProcessedBlock := db.GetBlockHeight()
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

			lastProcessedBlock, err = events.ProcessEvents(client, lastProcessedBlock+1, finalBlock, db.StoreEvents)
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
