module github.com/ansavin/system-monitoring

go 1.16

require (
	google.golang.org/grpc v1.42.0
	os/linux v1.0.0
	protobuf v1.0.0
)

replace os/linux => ./os/linux

replace protobuf => ./protobuf
