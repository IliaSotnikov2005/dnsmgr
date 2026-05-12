PROTO_DIR = proto
PROTO_FILE = $(PROTO_DIR)/dns.proto
OUT_DIR = proto

all: generate

generate:
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILE)

clean:
	rm -f $(PROTO_DIR)/*.pb.go

lint:
	golangci-lint run ./...

test:
	go test -v -race ./...

.PHONY: all generate clean lint test
