FROM golang:1.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY client/go.mod ./client/
COPY client/go.sum ./client/
COPY protobuf/go.mod ./protobuf/
COPY protobuf/go.sum ./protobuf/
COPY oslayer/go.mod ./oslayer/
COPY oslayer/go.sum ./oslayer/
COPY config/go.mod ./config/
COPY config/go.sum ./config/
RUN go mod download

COPY . ./
RUN go build -o system-monitor
ENTRYPOINT [ "./system-monitor" ]
CMD []
