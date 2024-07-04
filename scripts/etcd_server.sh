#!/bin/bash
set -e
TOKEN=token-01
BIN=../../etcd/bin
CLUSTER_STATE=new
CLUSTER=n0=http://127.0.0.1:$1,n1=http://127.0.0.1:$2,n2=http://127.0.0.1:$3

THIS_MACHINE=${4}

# get current entry from CLUSTER based on THIS_MACHINE
# e.g. if THIS_MACHINE=n1, then get n1=http://127.0.0.1:$1
# then split the string by '=' and get the second element
THIS_HOST=$(echo $CLUSTER | grep -o "${THIS_MACHINE}=http://[^,]*" | cut -d'=' -f2)
THIS_PORT=$(echo $THIS_HOST | cut -d':' -f3)

# etcd uses a different port for client requests
# lets just add 100 to the port number to get the client port number
THIS_CLIENT_PORT=$(($THIS_PORT + 100))
# maybe we should use 0.0.0.0 instead of the IP address
THIS_CLIENT_HOST="http://0.0.0.0:"$THIS_CLIENT_PORT

# print the host and client host
echo "Host: $THIS_HOST"
echo "Port: $THIS_PORT"
echo "Client Host: $THIS_CLIENT_HOST"
echo "Client Port: $THIS_CLIENT_PORT"

# Create a temporary directory to store data
mkdir -p /tmp/etcd

ID=$4;${BIN}/etcd \
  --name ${ID} \
  --data-dir /tmp/etcd/${ID} \
  --initial-advertise-peer-urls $THIS_HOST \
  --listen-peer-urls $THIS_HOST \
  --advertise-client-urls $THIS_CLIENT_HOST \
  --listen-client-urls $THIS_CLIENT_HOST \
  --initial-cluster ${CLUSTER} \
  --initial-cluster-state ${CLUSTER_STATE} \
  --initial-cluster-token ${TOKEN}
