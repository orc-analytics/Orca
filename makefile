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
	cd storage && docker-compose up -d



