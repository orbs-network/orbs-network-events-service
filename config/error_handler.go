package config

import (
	"github.com/orbs-network/govnr"
	"github.com/orbs-network/scribe/log"
)

type stdoutErrorer struct {
	logger log.Logger
}

func (s *stdoutErrorer) Error(err error) {
	println(err.Error())
}

func NewErrorHandler(logger log.Logger) govnr.Errorer {
	return &stdoutErrorer{
		logger: logger,
	}
}
