.PHONY: all build_proto build_store remove_store refresh_store

all: build_proto build_store

build_proto: .proto .proto_docs 
build_store: .spin_up_datalayer
stop_store: .stop_datalayer
remove_store: .remove_datalayer
refresh_store: .shut_down_datalayer .spin_up_datalayer
create_ssl: .create_ssl_cert

.proto:
	cd protobufs && protoc \
	--go_out=go \
	--go_opt=paths=source_relative \
	--go-grpc_out=go \
	--go-grpc_opt=paths=source_relative \
	*.proto

.proto_docs:
	cd protobufs && docker run --rm \
	-v ./../docs/:/out \
	-v ./:/protos \
	pseudomuto/protoc-gen-doc \
	--doc_opt=markdown,ProtocolBuffers.md

.stop_datalayer:
	cd storage && docker-compose stop

.remove_datalayer:
	cd storage && docker-compose down

.spin_up_datalayer:
	@if [ ! -d "./storage/_datalayer" ]; then \
        sudo mkdir -p ./storage/_datalayer; \
	fi
	cd storage && docker-compose up -d

.create_ssl_cert:
	@if [ ! -d "./storage/_ca" ]; then \
        sudo mkdir -p ./storage/_ca; \
	fi
	cd ./storage/_ca && sudo openssl req \
		-x509 -newkey rsa:4096 -sha256 -days 3650 \
  	-nodes -keyout key.pem -out cert.pem -subj "/CN=example.com" \
  	-addext "subjectAltName=DNS:example.com,DNS:*.example.com,IP:10.0.0.1"
