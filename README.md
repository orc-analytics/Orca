# Analytics Framework
The Analytics Framework Backend

This repository contains the back end scaffolding for the analytics framework.

## Functionality

## Development

### Spinning up a DB for local Development
When `make build_store` is run from the root directory, a Postgres DB along with a PGAdmin4 instance will be started. The neccessary SSL certificates will also be created. The following make commands exist to perform various manipulations on the local store:


```
- `make build_store` - To create the store from scratch, and start it
- `make start_store` - To start the store
- `stop_store` - To stop a running instance of the store
- `remove_store` - To remove an existing store and the data contents
- `redo_store` - Tear down a store, delete the data and spin it back up
- `create_ssl` - Create a fresh SSL certificate
```

To make available

### Documentation
Auto documentation of Protocol Buffer definitions is performed using the [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) repo. To use this it's recommended to install the Dockerised version. To do this, first install Docker. Then, pull the Protobuf autodoc tool with the command:


```bash
docker pull pseudomuto/protoc-gen-doc
```

