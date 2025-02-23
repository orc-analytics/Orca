# Orca

Streamline the analytics and amplify the insight.

### Building

Install the proto compiler:
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

Install the proto GRPC compiler:
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

Install the python GRPC code generator:

python -m pip install grpcio grpcio-tools

## Functionality

### Orca Execution Flow

1. Processors start up and register with the Orca-core service
2. Window arrives at the core
3. Core identifies affected algorithms
4. Core creates execution plan based on DAG
5. Tasks streamed to appropriate processors
6. Results flow back to core
7. Core triggers dependent algorithms when dependencies complete

### The Datalayer

The datalayer refers to the utilities required to connect algorithms to storage. This includes, writing windows, results and algorithm definitions.

Datalayers for the following data stores are as follows:

- `Postgres`

To implement a datalayer, the following methods should be exposed:

- TODO

### Fundamental Components

The fundamanetal components to the Predixus DB solution are:

- Windows
- Algorithms

Similarly to Apache Beam, windows define regions of interest to perform analytics on. In data pipelines where there are many available streams, a window will select a region in time over that stream to perform analyics.

Algorithms, perform the actual analytics. They are triggered by windows and can depend on eachother. This dependency allows results of 1 algorithm always to be present before another, so that the result can be used.

## Development

### Spinning up a DB for local Development

When `make build_store` is run from the root directory, a Postgres DB along with a PGAdmin4 instance will be started. The neccessary SSL certificates will also be created. The following make commands exist to perform various manipulations on the local store:

```
- `make build_store`    - To create the store from scratch, and start it
- `make start_store`    - To start the store
- `stop_store`          - To stop a running instance of the store
- `remove_store`        - To remove an existing store and the data contents
- `redo_store`          - Tear down a store, delete the data and spin it back up
```

### Documentation

Auto documentation of Protocol Buffer definitions is performed using the [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) repo. To use this it's recommended to install the Dockerised version. To do this, first install Docker. Then, pull the Protobuf autodoc tool with the command:

```bash
docker pull pseudomuto/protoc-gen-doc
```
