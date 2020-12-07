FROM golang:buster

RUN go get -u github.com/gorilla/mux

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /Key-Value-Store

COPY . .

RUN go build .

EXPOSE 8080

CMD ["./Key-Value-Store"]