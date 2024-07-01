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
# Stage 1b (Ratis)                                                            #
###############################################################################
FROM maven:3.9.8 AS ratis-builder
RUN git clone https://github.com/apache/ratis.git /src/ratis
WORKDIR /src/ratis
RUN mvn clean package -DskipTests

###############################################################################
# Stage 2 (to create a downsized "container executable", ~5MB)                #
###############################################################################
# If you need SSL certificates for HTTPS, replace `FROM SCRATCH` with:
#
#   FROM alpine:3.17.1
#   RUN apk --no-cache add ca-certificates
#
#FROM scratch
FROM openjdk:17
WORKDIR /root/
COPY --from=builder /go/src/dstest /root/dstest
COPY --from=ratis-builder /src/ratis /root/ratis

#EXPOSE 8123
WORKDIR /root/dstest/cmd/dstest/
ENTRYPOINT ["./main"]
CMD ["run", "-c", "./config/config.yml"]
