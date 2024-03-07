.PHONY: build_proto build_store remove_store refresh_store

build_proto: .proto .proto_docs .spin_up_datalayer
build_store: .spin_up_datalayer
remove_store: .shut_down_datalayer
refresh_store: .shut_down_datalayer .spin_up_datalayer
	
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


.shut_down_datalayer:
	cd storage && docker-compose down

.spin_up_datalayer:
	@if [ ! -d "./storage/datalayer" ]; then \
        sudo mkdir -p ./storage/datalayer; \
        sudo chmod 777 ./storage/datalayer; \
    fi
	cd storage && docker-compose up -d

