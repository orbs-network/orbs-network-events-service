package config

import (
	"github.com/orbs-network/scribe/log"
	"os"
)

func GetLogger() log.Logger {
	logger := log.GetLogger().
		WithTags(log.String("app", "signer")).
		WithOutput(log.NewFormattingOutput(os.Stdout, log.NewHumanReadableFormatter()))

	return logger
}
