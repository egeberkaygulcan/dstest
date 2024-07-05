# syntax=docker/dockerfile:1

# We use a multi-stage build setup.
# (https://docs.docker.com/build/building/multi-stage/)

###############################################################################
# Stage 1 (to create a "build" image, ~850MB)                                 #
###############################################################################
# Image from https://hub.docker.com/_/golang
FROM golang:1.22.2 AS builder
# smoke test to verify if golang is available
RUN go version

ARG PROJECT_VERSION

COPY . /go/src/dstest/
WORKDIR /go/src/dstest/
RUN set -Eeux && \
    go mod download && \
    go mod verify

RUN GOOS=linux GOARCH=amd64 \
    go build \
    -trimpath \
    -ldflags="-w -s -X 'main.Version=${PROJECT_VERSION}'" \
    -o cmd/dstest/main cmd/dstest/main.go
RUN go test -cover -v ./...

###############################################################################
# Stage 2A (Ratis, Zookeeper)                                                 #
###############################################################################
FROM maven:3.9.8 AS java-builder

# Ratis
RUN git clone https://github.com/apache/ratis.git /src/ratis
WORKDIR /src/ratis
RUN mvn clean package -DskipTests

# Zookeeper
RUN git clone https://github.com/apache/zookeeper.git /src/zookeeper
WORKDIR /src/zookeeper
RUN mvn clean package -DskipTests

###############################################################################
# Stage 2B (etcd)                                                             #
###############################################################################
# Image from https://hub.docker.com/_/golang
FROM golang:1.22.2 AS go-builder

RUN git clone https://github.com/etcd-io/etcd.git /src/etcd
WORKDIR /src/etcd
RUN go mod download
RUN go mod verify
RUN ./scripts/build.sh

###############################################################################
# Stage 3 (to create a downsized "container executable", ~5MB)                #
###############################################################################
FROM ubuntu:24.04

# Install Java and Go runtime dependencies
RUN apt-get update
RUN apt-get install -y openjdk-17-jre-headless golang-go

WORKDIR /root/
COPY --from=builder /go/src/dstest /root/dstest
COPY --from=java-builder /src/ratis /root/ratis
COPY --from=java-builder /src/zookeeper /root/zookeeper
COPY --from=go-builder /src/etcd /root/etcd

WORKDIR /root/dstest/cmd/dstest/
ENTRYPOINT ["./main"]
CMD ["run", "-c", "./config/config.yml"]
