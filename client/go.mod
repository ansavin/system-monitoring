module github.com/ansavin/system-monitoring/client

go 1.16

require (
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1 // indirect
	protobuf v1.0.0
)

replace protobuf => ../protobuf
