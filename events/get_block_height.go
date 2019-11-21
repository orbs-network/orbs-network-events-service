package events

import (
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
)

func GetBlockHeight(client *orbs.OrbsClient, account *orbs.OrbsAccount) (uint64, error) {
	query, err := client.CreateQuery(account.PublicKey, "_Info", "isAlive")
	res, err := client.SendQuery(query)
	if err != nil {
		return 0, err
	}

	return uint64(res.BlockHeight), err
}
