package events

type StoredEvent struct {
	ContractName string
	EventName    string

	BlockHeight uint64
	Timestamp   uint64
	Arguments   []interface{}
}
