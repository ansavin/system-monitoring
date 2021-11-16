FROM golang:1.16-alpine
WORKDIR /app
COPY . ./
RUN go build -o system-monitor
CMD [ "./system-monitor" ]
