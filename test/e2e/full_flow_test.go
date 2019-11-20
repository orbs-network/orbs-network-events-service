package e2e

import (
	"context"
	"github.com/orbs-network/orbs-network-events-service/boostrap"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func removeDB() {
	os.RemoveAll("./vchain-42.bolt")
}

func TestFullFlow(t *testing.T) {
	removeDB()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, err := boostrap.NewCluster(ctx, &config.Config{
		Endpoint:        "http://localhost:8080",
		VirtualChains:   []uint32{42},
		DB:              "./",
		PollingInterval: 10 * time.Millisecond,
	}, config.GetLogger())
	require.NoError(t, err)

	time.Sleep(time.Minute)
}
