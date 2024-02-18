build_proto: .proto .proto_docs

.proto:
	protoc --proto_path=${PROTOPATH} --proto_path=protobufs/ protobufs/*.proto  --go_out=protobufs/go/ --go_opt=paths=source_relative

.proto_docs:
	cd protobufs && docker run --rm \
	-v ./docs:/out \
	-v ./:/protos \
	pseudomuto/protoc-gen-doc \
	--doc_opt=markdown,ProtocolBuffers.md