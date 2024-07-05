#!/bin/bash
BIN=../../../ratis/ratis-examples/src/main/bin
PEERS=n0:127.0.0.1:6000,n1:127.0.0.1:6001,n2:127.0.0.1:6002
${BIN}/client.sh arithmetic assign --name c --value a+b --peers ${PEERS}