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
    primitives.contract_name contract_name = 1;
    string event_name = 2;

    primitives.block_height block_height = 3;
    primitives.timestamp_nano timestamp = 4;

    primitives.sha256 txhash = 5;
    protocol.ExecutionResult execution_result = 6;
    uint32 index = 7; // suggestion from Noam to preserve order

    primitives.packed_argument_array arguments = 8;
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
    primitives.protocol_version protocol_version = 1;
    primitives.virtual_chain_id virtual_chain_id = 2;

    primitives.contract_name contract_name = 5;
    repeated string event_name = 6;

    primitives.block_height from_block = 7;
    primitives.block_height to_block = 8;

    primitives.timestamp_nano from_time = 9;
    primitives.timestamp_nano to_time = 10;
}

```


### Response

```protobuf
message IndexerResponse {
    repeated protocol.IndexedEvent events = 1;
}
```