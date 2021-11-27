# system-monitoring

[![Go Report Card](https://goreportcard.com/badge/github.com/ansavin/system-monitoring)](https://goreportcard.com/report/github.com/ansavin/system-monitoring)
[![Golang-CI](https://github.com/ansavin/system-monitoring/actions/workflows/golang-ci.yml/badge.svg)](https://github.com/ansavin/system-monitoring/actions/workflows/golang-ci.yml)

A small monitoring daemon that sends info about server's health via protobuf

## Requirements

* one of the supported OS (MacOS, Linux)
* Golang v1.16 or Docker (only for Linux)

## Functionality & Features

* Calculating:
  * Load average (MacOS, Linux)
  * CPU load (MacOS, Linux)
  * Disks load (Linux only)
* Unit-tests
* Integration tests
* Simple client for daemon
* Some statistics works under different OS (MacOS, Linux)

## Tested on

* Ubuntu 18.04
* Ubuntu 20.04
* Fedora 34
* MacOS 11

## Building

### Locally

* To build server, type
  `make grpc-server`

* To build a client, type
  `make grpc-client`

### Docker

* To build server docker image, type
  `make docker`

## Configuring

### Server

* To disable grabbing & sending statistics, edit `config.yml`

* To run server on different port, use `-p` option:
  `sudo ./grpc-server -p 8089`

* All CLI options are available by running server with `-h` option

### Client

* To run client on different port, use `-p` option (default port is `8088`):
  `sudo ./grpc-client -p 8089`

* To get data averaged for, for example, 5 seconds, use `-a` option (default is 3 sec):
  `sudo ./grpc-client -a 5`

* To get messages from server every, for example, 4 seconds, use `-m` option (default is 3 sec):
  `sudo ./grpc-client -m 4`

* All CLI options are available by running server with `-h` option

## Running

### Directly

* To run a server, build a server binary and start it with root privileges (we need them to examine FS Utilization)
  `sudo ./grpc-server`

* To run a client, build a client binary and then start it:
  `./grpc-client`

### In Docker

* To run a service in docker, type
  `make docker-server`

* To run a client in docker, type
  `make docker-client`

## Developing

* To regenerate GRPC implementation, type
  `make grpc-autogen`

* To run integration test locally, type
  `make integration-test`

* To run unit tests locally, type
  `make test`
