#!/bin/bash
BIN=../../../ratis/ratis-examples/src/main/bin
PEERS=n0:127.0.0.1:$1,n1:127.0.0.1:$2,n2:127.0.0.1:$3;

ID=$4;${BIN}/server.sh arithmetic server --id ${ID} --storage /tmp/ratis/${ID} --peers ${PEERS}