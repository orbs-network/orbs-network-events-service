package bootstrap

import (
	"context"
	"github.com/orbs-network/orbs-network-events-service/bootstrap/httpserver"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/scribe/log"
)

func NewCluster(ctx context.Context, cfg *config.Config, logger log.Logger) (*httpserver.HttpServer, error) {
	var nodes []*Node
	for _, vcid := range cfg.VirtualChains {
		node, err := NewNode(ctx, cfg, logger, vcid)
		if err != nil {
			return nil, err
		}

		nodes = append(nodes, node)
	}

	apis := make(map[uint32]types.Indexer)
	for _, node := range nodes {
		apis[node.vcid] = node.api
	}

	server := httpserver.NewHttpServer(ctx, httpserver.NewServerConfig("0.0.0.0:9201", false), logger, apis)
	return server, nil
}
