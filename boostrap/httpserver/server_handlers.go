// Copyright 2019 the orbs-network-go authors
// This file is part of the orbs-network-go library in the Orbs project.
//
// This source code is licensed under the MIT license found in the LICENSE file in the root directory of this source tree.
// The above notice should be included in all copies or substantial portions of the software.

package httpserver

import (
	"encoding/json"
	"github.com/orbs-network/orbs-network-events-service/config"
	"github.com/orbs-network/orbs-spec/types/go/protocol/client"
	"github.com/orbs-network/orbs-spec/types/go/services"
	"github.com/orbs-network/scribe/log"
	"net/http"
)

type IndexResponse struct {
	Status      string
	Description string
	Version     config.Version
}

// Serves both index and 404 because router is built that way
func (s *HttpServer) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	data, _ := json.MarshalIndent(IndexResponse{
		Status:      "OK",
		Description: "ORBS event indexer API",
		Version:     config.GetVersion(),
	}, "", "  ")

	_, err := w.Write(data)
	if err != nil {
		s.logger.Info("error writing index.json response", log.Error(err))
	}
}

func (s *HttpServer) robots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	_, err := w.Write([]byte("User-agent: *\nDisallow: /\n"))
	if err != nil {
		s.logger.Info("error writing robots.txt response", log.Error(err))
	}
}

func (s *HttpServer) getEventsHandler(w http.ResponseWriter, r *http.Request) {
	bytes, e := readInput(r)
	if e != nil {
		s.writeErrorResponseAndLog(w, e)
		return
	}

	result, err := s.indexer.GetEvents(r.Context(), (&services.GetEventsInputBuilder{
		ClientRequest: client.IndexerRequestBuilderFromRaw(bytes),
	}).Build())

	if result != nil && result.ClientResponse != nil {
		// FIXME http codes
		s.writeMembuffResponse(w, result, 200, err)
	} else {
		s.writeErrorResponseAndLog(w, &httpErr{http.StatusInternalServerError, log.Error(err), err.Error()})
	}
}
