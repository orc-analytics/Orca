# Analytics Framework
The Analytics Framework Backend

This repository contains the back end scaffolding for the analytics framework.

## Functionality

## Development
### Spinning up a DB for local Development
When `make` is run from the root directory, a Postgres DB along with a PGAdmin4 instance will be started. 

To make available

### Documentation
Auto documentation of Protocol Buffer definitions is performed using the [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) repo. To use this it's recommended to install the Dockerised version. To do this, first install Docker. Then, pull the Protobuf autodoc tool with the command:
```bash
docker pull pseudomuto/protoc-gen-doc
```

