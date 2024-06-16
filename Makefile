build:
	@go build -o bin/p2pgrpc

run: build
	@./bin/p2pgrpc

proto: 
	protoc --go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/*.proto

.PHONY: proto