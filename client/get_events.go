package client

import (
	"bytes"
	"fmt"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type GetEventsQuery struct {
	VirtualChainId uint32
	ContractName   string
	EventName      []string

	FromBlock uint64
	ToBlock   uint64

	FromTime int64
	ToTime   int64

	Filters [][]interface{}
}

func GetEvents(endpoint string, query GetEventsQuery) (events []*types.IndexedEvent, err error) {
	var filters [][]byte
	for _, nativeFilters := range query.Filters {
		filtersAsBytes, err := protocol.ArgumentArrayFromNatives(nativeFilters)
		if err != nil {
			return nil, err
		}
		filters = append(filters, filtersAsBytes.Raw())
	}

	requestBody := (&types.IndexerRequestBuilder{
		VirtualChainId: query.VirtualChainId,
		ContractName:   query.ContractName,
		EventName:      query.EventName,
		FromBlock:      uint64(query.FromBlock),
		ToBlock:        uint64(query.ToBlock),
		FromTime:       uint64(query.FromTime),
		ToTime:         uint64(query.ToTime),
		Filters:        filters,
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

	indexerResponse := types.IndexerResponseReader(data)
	for i := indexerResponse.EventsIterator(); i.HasNext(); {
		events = append(events, i.NextEvents())
	}

	return events, nil
}
