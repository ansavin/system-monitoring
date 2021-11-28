GO := $(shell which go)
SRC_SERVER := $(shell find . -type f -name '*.go' -not -path "./client")
SRC_CLIENT := $(shell find ./client -type f -name '*.go')

.PHONY: test
test:
	sudo $(GO) test oslayer -race -count 10

.PHONY: integration-test
integration-test:
	./integration-test.sh

grpc-server: $(SRC_SERVER)
	echo $(GO)
	$(GO) build -o grpc-server

grpc-client: $(SRC_CLIENT)
	cd client && \
	$(GO) build -o ../grpc-client

.PHONY: grpc-autogen
grpc-autogen: protobuf/system-monitor.proto
	protoc --go_out=. --go-grpc_out=. protobuf/system-monitor.proto

.PHONY: docker
docker:
	docker build -t system-monitor .

.PHONY: docker-server
docker-server: docker
	docker run \
	--rm \
	--net="host" \
	--pid="host" \
	-v "/:/host:ro,rslave" \
	--name system-monitor-server \
	system-monitor -r /host

.PHONY: docker-client
docker-client: docker
	docker run \
	--rm \
	--net="host" \
	--pid="host" \
	--entrypoint="go" \
	-v "/:/host:ro,rslave" \
	--name system-monitor-client \
	system-monitor run client/main.go

.PHONY: docker-test
docker-test: docker
	docker run \
	--rm \
	--net="host" \
	--pid="host" \
	--entrypoint="go" \
	-v "/:/host:ro,rslave" \
	--name system-monitor-test \
	system-monitor test oslayer

.PHONY: clean
clean:
	docker rm system-monitor-server 1>/dev/null 2>&1 || true
	docker rm system-monitor-client 1>/dev/null 2>&1 || true
	docker rm system-monitor-test 1>/dev/null 2>&1 || true
	docker rmi system-monitor 1>/dev/null 2>&1 || true
	rm -f grpc-server 1>/dev/null 2>&1
	rm -f grpc-client 1>/dev/null 2>&1