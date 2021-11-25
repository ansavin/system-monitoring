# system-monitoring

A small monitoring daemon that sends info about server's health via protobuf

## Requirements

* one of the supported OS
* docker
* internet access (for pulling docker images)

## Functionality & Features

* Calculating:
  * Load average
  * CPU load
  * Disks load
* Unit-tests
* Integration tests
* Simple client for daemon
* Some statistics works under different OS (MacOS, Linux)

## Tested on

* Ubuntu 18.04
* Ubuntu 20.04
* Fedora 34
* MacOS 11

## Running

* To run server locally, type
  `make grpc-server`

  And then start it with root privileges (we need them to examine FS Utilization)
  `./grpc-sever`

* To run a client locally, type
  `make grpc-server`

  And then start it:
  `./grpc-sever 3 5`
  where 3 - is time between messages with statistics,
  5 - averaging statistics time

* To run a service in docker, type
  `make docker-run`
  where 3 - is time between messages with statistics,
  5 - averaging statistics time

* To run a client in docker, after server startup type
  `docker exec -it system-monitor ./client/client 3 5`
  where 3 - is time between messages with statistics,
  5 - averaging statistics time

## Developing

* To regenerate GRPC implementation, type
  `make grpc-autogen`

* To run integration test locally, type
  `make integration-test`

* To run unit tests locally, type
  `make test`

* To run unit tests in docker, type
  `make docker-test`
