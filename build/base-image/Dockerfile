FROM golang:1.21.6-alpine3.19

RUN apk add make git protoc

RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.25.0
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.19.1
RUN go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.32
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

WORKDIR /go/src


