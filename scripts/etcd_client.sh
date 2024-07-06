#!/bin/bash
set -e

BIN=../../etcd/bin
PEERS=127.0.0.1:6100,127.0.0.1:6101,127.0.0.1:6102

export ETCDCTL_API=3
${BIN}/etcdctl --endpoints=$ENDPOINTS get foo
