module github.com/ansavin/system-monitoring

go 1.16

require (
	config v1.0.0
	google.golang.org/grpc v1.42.0
	oslayer v1.0.0
	protobuf v1.0.0
)

replace oslayer => ./oslayer

replace protobuf => ./protobuf

replace config => ./config
