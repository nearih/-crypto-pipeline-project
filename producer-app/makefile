## dev: run to start server
.PHONY: dev
dev:
	go run *.go


## protoV1: create proto3V1
.PHONY: proto
proto:
	rm -f pkg/*
	protoc --proto_path=server --go_out=generated --go-grpc_out=generated pipeline.proto 

