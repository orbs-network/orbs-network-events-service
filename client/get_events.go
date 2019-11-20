package client

import (
	"bytes"
	"fmt"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type GetEventsQuery struct {
	VirtualChainId uint64
	ContractName   string
	EventName      []string

	FromBlock uint64
	ToBlock   uint64

	FromTime int64
	ToTime   int64
}

func GetEvents(endpoint string, query GetEventsQuery) (events []*protocol.IndexedEvent, err error) {
	requestBody := (&client.IndexerRequestBuilder{
		VirtualChainId: primitives.VirtualChainId(query.VirtualChainId),
		ContractName:   primitives.ContractName(query.ContractName),
		EventName:      query.EventName,
		FromBlock:      primitives.BlockHeight(query.FromBlock),
		ToBlock:        primitives.BlockHeight(query.ToBlock),
		FromTime:       primitives.TimestampNano(query.FromTime),
		ToTime:         primitives.TimestampNano(query.ToTime),
	}).Build().Raw()

	res, err := http.Post(fmt.Sprintf("%s/api/v1/get-events", endpoint), "application/membuffers", bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("%s: %s", res.Status, res.Header.Get("X-ORBS-ERROR-DETAILS"))
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	indexerResponse := client.IndexerResponseReader(data)
	for i := indexerResponse.EventsIterator(); i.HasNext(); {
		events = append(events, i.NextEvents())
	}

	return events, nil
}
