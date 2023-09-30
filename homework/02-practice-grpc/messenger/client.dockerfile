FROM golang:alpine

WORKDIR /practice-grpc/solutions
RUN mkdir client proto
COPY client/go.* client/
COPY proto/go.* proto/

RUN for d in proto client; do (cd $d && go mod tidy && go mod download -x); done || exit 1

COPY client client
COPY proto proto
RUN for d in proto client; do (cd $d && go build .); done || exit 1

CMD ./client/client
