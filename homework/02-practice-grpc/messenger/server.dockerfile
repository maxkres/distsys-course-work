FROM golang:alpine

WORKDIR /practice-grpc/solutions
RUN mkdir proto server
COPY proto/go.* proto/
COPY server/go.* server/

RUN for d in proto server; do (cd $d && go mod tidy && go mod download -x); done || exit 1

COPY proto proto
COPY server server
RUN for d in proto server; do (cd $d && go build .); done || exit 1


CMD ./server/server
