.PHONY: all build_proto build_store remove_store refresh_store

all: .proto .datalayer
proto: .proto
datalayer: .datalayer

build_proto: .proto .proto_docs 
build_store: .create_ssl_cert .spin_up_datalayer
start_store: .start_datalayer
stop_store: .stop_datalayer
remove_store: .remove_datalayer
redo_store: .remove_datalayer .remove_store_cache .create_ssl_cert .spin_up_datalayer
create_ssl: .create_ssl_cert
test: .test_all

.proto:
	cd core/protobufs && protoc \
	--go_out=go \
	--go_opt=paths=source_relative \
	--go-grpc_out=go \
	--go-grpc_opt=paths=source_relative \
	*.proto vendor/*.proto
	cd core/protobufs && python -m grpc_tools.protoc \
    --proto_path=./ \
    --python_out=./python \
    --pyi_out=./python \
    --grpc_python_out=./python \
	*.proto vendor/*.proto


.datalayer:
	sqlc vet -f core/datalayer/postgresql/sqlc.yaml
	sqlc generate -f core/datalayer/postgresql/sqlc.yaml

.proto_docs:
	cd core/protobufs && docker run --rm \
	-v ./../docs/:/out \
	-v ./:/protos \
	pseudomuto/protoc-gen-doc \
	--doc_opt=markdown,ProtocolBuffers.md

.stop_datalayer:
	cd local_storage && docker-compose stop

.start_datalayer:
	cd local_storage && docker-compose start

.remove_datalayer:
	cd local_storage && docker-compose down
	docker volume remove local_storage_datalayer

.spin_up_datalayer:
	@if [ ! -d "./local_storage/_datalayer" ]; then \
        sudo mkdir -p ./local_storage/_datalayer; \
				sudo chmod 777 ./local_storage/_datalayer; \
	fi
	cd local_storage && docker-compose up -d

.remove_store_cache:
	sudo rm -rf local_storage/_*

.create_ssl_cert:
	@if [ ! -d "./local_storage/_ca" ]; then \
        sudo mkdir -p ./local_storage/_ca; \
				sudo chmod 777 ./local_storage/_ca; \
	fi
	cd ./local_storage/_ca && \
		sudo openssl req -new -text -passout pass:abcd -subj /CN=localhost -out server.req -keyout privkey.pem
	cd ./local_storage/_ca && \
		sudo openssl rsa -in privkey.pem -passin pass:abcd -out server.key
	cd ./local_storage/_ca && \
		sudo openssl req -x509 -in server.req -text -key server.key -out server.crt
	sudo chown 0:70 local_storage/_ca/server.key
	sudo chmod 640 local_storage/_ca/server.key

.test_all:
	cd core && go test ./internal/... -v


## CLI Build scripts
BINARY_NAME = orca

# disabled CGO to produce statically-linked binaries
export CGO_ENABLED = 0

# flags for stripping debugging info
LDFLAGS = -s -w

# build command
BUILD = go build -ldflags "$(LDFLAGS)" -o

# platform targets
.PHONY: cli_all cli_clean cli_linux cli_windows cli_mac_arm cli_mac_intel

cli_all: cli_linux cli_windows cli_mac_arm cli_mac_intel

cli_linux:
	cd cli && GOOS=linux GOARCH=amd64 $(BUILD) ../build/$(BINARY_NAME)-linux .

cli_windows:
	cd cli && GOOS=windows GOARCH=amd64 $(BUILD) ../build/$(BINARY_NAME)-windows.exe .

cli_mac_arm:
	cd cli && GOOS=darwin GOARCH=arm64 $(BUILD) ../build/$(BINARY_NAME)-mac-arm .

cli_mac_intel:
	cd cli && GOOS=darwin GOARCH=amd64 $(BUILD) ../build/$(BINARY_NAME)-mac-intel .

cli_clean:
	rm -rf build/
