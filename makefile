.PHONY: all build_proto build_store remove_store refresh_store

all: build_proto

build_proto: .proto .proto_docs 
build_store: .create_ssl_cert .spin_up_datalayer
start_store: .start_datalayer
stop_store: .stop_datalayer
remove_store: .remove_datalayer
redo_store: .remove_datalayer .remove_store_cache .create_ssl_cert .spin_up_datalayer
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

.start_datalayer:
	cd storage && docker-compose start

.remove_datalayer:
	cd storage && docker-compose down
	docker volume remove storage_datalayer

.spin_up_datalayer:
	@if [ ! -d "./storage/_datalayer" ]; then \
        sudo mkdir -p ./storage/_datalayer; \
	fi
	cd storage && docker-compose up -d

.remove_store_cache:
	sudo rm -rf storage/_*

.create_ssl_cert:
	@if [ ! -d "./storage/_ca" ]; then \
        sudo mkdir -p ./storage/_ca; \
				sudo chmod 777 ./storage/_ca; \
	fi
	cd ./storage/_ca && \
		sudo openssl req -new -text -passout pass:abcd -subj /CN=localhost -out server.req
	cd ./storage/_ca && \
		sudo openssl rsa -in privkey.pem -passin pass:abcd -out server.key
	cd ./storage/_ca && \
		sudo openssl req -x509 -in server.req -text -key server.key -out server.crt

	sudo chown 0:70 storage/_ca/server.key
	sudo chmod 640 storage/_ca/server.key
