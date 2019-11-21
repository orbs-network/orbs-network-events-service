# Event Indexer

The goal of event indexer service is to provide queryable data access to smart contract events.
The events are indexed from the different blocks and a simple API enables any app to query which events happened at what time for a specific smart contract, without scanning each block.

To make life easier for SDK users, the service should be part of node deployment, similar to the signer service.

## Interfaces

Below are the protobuf interfaces for talking to the service

### IndexedEvent

This is the response data model of a single event occurance

```protobuf
message IndexedEvent {
    string contract_name = 1;
    string event_name = 2;

    uint64 block_height = 3;
    uint64 timestamp = 4;

    bytes txhash = 5;
    ExecutionResult execution_result = 6;
    uint32 index = 7; // suggestion from Noam to preserve order

    bytes arguments = 8;
}
```

### Request

The types of queries that we currently aim to answer:

* events by contract name and a list of event names
* events by range of block height
* events by range of dates

The contract name and event names (one or more) must be supplied.

Only one of the ranges should be supplied, either the block range or the time. If both are supplied an error will be returned. 

Currently, results are not paginated and all the data available for the current filter is returned
It is also not possible to filter according to the event message/data.

```protobuf
message IndexerRequest {
    uint32 protocol_version = 1;
    uint32 virtual_chain_id = 2;

    string contract_name = 5;
    repeated string event_name = 6;

    uint64 from_block = 7;
    uint64 to_block = 8;

    uint64 from_time = 9;
    uint64 to_time = 10;
}

```


### Response

```protobuf
message IndexerResponse {
    repeated IndexedEvent events = 1;
}
```

## Development

### Testing

```
gamma-cli start-local -env experimental
go test ./... -v
```

### Recompiling protos

```
cd types && membufc -g --go-ctx *.proto
```
