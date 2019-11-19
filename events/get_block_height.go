package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
)

func GetBlockHeight(client *orbs.OrbsClient, account *orbs.OrbsAccount) (primitives.BlockHeight, error) {
	query, err := client.CreateQuery(account.PublicKey, "_Info", "isAlive")
	res, err := client.SendQuery(query)
	if err != nil {
		return 0, err
	}

	return primitives.BlockHeight(res.BlockHeight), err
}
