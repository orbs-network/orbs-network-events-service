package http_api

import (
	"context"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/services/storage"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/orbs-network/scribe/log"
	"github.com/pkg/errors"
)

type service struct {
	cfg    *config.Config
	logger log.Logger

	dbs map[uint32]storage.Storage
}

func NewHttpAPI(cfg *config.Config, logger log.Logger) (services.Indexer, error) {
	serviceLogger := logger.WithTags(log.Service("http"))

	dbs := make(map[uint32]storage.Storage)
	for _, vcid := range cfg.VirtualChains {
		db, err := storage.NewStorageForChain(serviceLogger, vcid, true)
		if err != nil {
			return nil, err
		}
		dbs[vcid] = db
	}

	return &service{
		cfg:    cfg,
		logger: logger,
		dbs:    dbs,
	}, nil
}

func (s *service) GetEvents(ctx context.Context, input *services.GetEventsInput) (*services.GetEventsOutput, error) {
	if input.ClientRequest().ContractName() == "" {
		return nil, errors.New("contract name is required")
	}

	if names := input.ClientRequest().EventNameIterator(); names.HasNext() {
		return nil, errors.New("event name is required")
	}

	vcid := input.ClientRequest().VirtualChainId()
	if vcid == 0 {
		return nil, errors.New("virtual chain id is required")
	}

	db, found := s.dbs[uint32(vcid)]
	if !found {
		return nil, errors.New("virtual chain not found")
	}

	events, err := db.GetEvents(input.ClientRequest())
	if err != nil {
		return nil, err
	}

	var clientResponseEvents []*protocol.IndexedEventBuilder
	for _, event := range events {
		clientResponseEvents = append(clientResponseEvents, protocol.IndexedEventBuilderFromRaw(event.Raw()))
	}

	return (&services.GetEventsOutputBuilder{
		ClientResponse: &client.IndexerResponseBuilder{
			Events: clientResponseEvents,
		},
	}).Build(), nil
}
