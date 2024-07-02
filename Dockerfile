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
# Stage 3 (to create a downsized "container executable", ~5MB)                #
###############################################################################
FROM openjdk:17
WORKDIR /root/
COPY --from=builder /go/src/dstest /root/dstest
COPY --from=java-builder /src/ratis /root/ratis
COPY --from=java-builder /src/zookeeper /root/zookeeper

WORKDIR /root/dstest/cmd/dstest/
ENTRYPOINT ["./main"]
CMD ["run", "-c", "./config/config.yml"]
