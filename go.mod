module github.com/orbs-network/orbs-network-events-service

go 1.12

require (
	github.com/orbs-network/boyarin v0.20.0 // indirect
	github.com/orbs-network/govnr v0.2.0
	github.com/orbs-network/orbs-client-sdk-go v0.12.0
	github.com/orbs-network/orbs-spec v0.0.0-20191106111111-a9a3678f9401
	github.com/orbs-network/scribe v0.2.2
	github.com/pkg/errors v0.8.1
	github.com/stretchr/testify v1.4.0
	go.etcd.io/bbolt v1.3.3
)

replace github.com/orbs-network/orbs-spec => ../orbs-spec
