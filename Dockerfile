FROM golang:1.16-alpine
WORKDIR /app
COPY . ./
RUN cd client && go build -o client && cd ..
RUN go build -o system-monitor
CMD [ "./system-monitor" ]
