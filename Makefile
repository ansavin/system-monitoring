GO := $(shell which go)

.PHONY: test
test:
	sudo $(GO) test oslayer -race -count 10

grpc-server:
	$(GO) build -o grpc-server

grpc-client:
	cd client && \
	$(GO) build -o ../grpc-client

.PHONY: all
docker:
	docker build -t system-monitor .

.PHONY: docker-run
docker-run: docker
	docker rm system-monitor; \
	docker run \
	--net="host" \
	--pid="host" \
	-v "/:/host:ro,rslave" \
	--name system-monitor \
	system-monitor

.PHONY: docker-test
docker-test: docker
	docker rm system-monitor-test; \
	docker run \
	--net="host" \
	--pid="host" \
	--entrypoint="go" \
	-v "/:/host:ro,rslave" \
	--name system-monitor-test \
	system-monitor test os/linux

.PHONY: clean
clean:
	docker rm system-monitor 1>/dev/null 2>&1
	docker rm system-monitor-test 1>/dev/null 2>&1
	docker rmi system-monitor 1>/dev/null 2>&1