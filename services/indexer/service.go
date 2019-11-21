package indexer

import (
	"context"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/services/storage"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/scribe/log"
	"github.com/pkg/errors"
)

type service struct {
	cfg    *config.Config
	logger log.Logger

	db storage.Storage
}

func NewIndexer(cfg *config.Config, logger log.Logger, db storage.Storage) (types.Indexer, error) {
	serviceLogger := logger.WithTags(log.Service("http"))

	return &service{
		cfg:    cfg,
		logger: serviceLogger,
		db:     db,
	}, nil
}

func (s *service) GetEvents(ctx context.Context, input *types.GetEventsInput) (*types.GetEventsOutput, error) {
	if input.ClientRequest().ContractName() == "" {
		return nil, errors.New("contract name is required")
	}

	if names := input.ClientRequest().EventNameIterator(); !names.HasNext() {
		return nil, errors.New("event name is required")
	}

	vcid := input.ClientRequest().VirtualChainId()
	if vcid == 0 {
		return nil, errors.New("virtual chain id is required")
	}

	events, err := s.db.GetEvents(input.ClientRequest())
	if err != nil {
		return nil, err
	}

	var clientResponseEvents []*types.IndexedEventBuilder
	for _, event := range events {
		clientResponseEvents = append(clientResponseEvents, types.IndexedEventBuilderFromRaw(event.Raw()))
	}

	return (&types.GetEventsOutputBuilder{
		ClientResponse: &types.IndexerResponseBuilder{
			Events: clientResponseEvents,
		},
	}).Build(), nil
}
