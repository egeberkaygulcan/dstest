#!/bin/bash
set -e
BIN=/Users/berkay/Documents/Research/ratis/ratis-examples/src/main/bin
PEERS=n0:127.0.0.1:6000,n1:127.0.0.1:6001,n2:127.0.0.1:6002

ID=$1; ${BIN}/server.sh arithmetic server --id ${ID} --storage /tmp/ratis/${ID} --peers ${PEERS}