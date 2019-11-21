package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1/events"
)

var PUBLIC = sdk.Export(release)
var SYSTEM = sdk.Export(_init)
var EVENTS = sdk.Export(MovieRelease)

func MovieRelease(name string, year uint32, lead string) {}

func _init() {

}

// TODO add param to emit many events
func release(name string, year uint32, lead string) {
	events.EmitEvent(MovieRelease, name, year, lead)
}
