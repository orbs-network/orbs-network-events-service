package storage

import (
	"bytes"
	"github.com/orbs-network/orbs-network-events-service/types"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

func matchEvent(event *types.IndexedEvent, filters []*protocol.ArgumentArray) bool {
	filtersCount := len(filters)
	if filtersCount == 0 {
		return true
	}

	eventArguments := protocol.ArgumentArrayReader(event.Arguments())
	i := 0
	match := false
	for argumentsIterator := eventArguments.ArgumentsIterator(); argumentsIterator.HasNext(); i++ {
		arg := argumentsIterator.NextArguments()

		if filtersCount >= i+1 {
			filter := filters[i]
			// Only direct matches
			if filter.ArgumentsIterator().HasNext() {
				filterArg0 := filter.ArgumentsIterator().NextArguments()
				println("matching", arg.String(), "against", filterArg0.String())
				match = bytes.Equal(arg.Raw(), filterArg0.Raw())
				if match { // FIXME
					return true
				}
			}
		}
	}

	return match
}
