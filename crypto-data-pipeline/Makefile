## dev: run to start server
.PHONY: dev
dev:
	go run *.go

## proto: create proto
.PHONY: proto
proto:
	protoc --proto_path=server/grpc/protos --go_out=generated --go-grpc_out=generated pipeline.proto
