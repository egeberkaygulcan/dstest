#!/bin/bash
set -e

pkill etcd
rm -rf /tmp/etcd
echo "Killed all etcd instances."
