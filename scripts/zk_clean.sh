#!/bin/bash
BIN=../../../zookeeper
# for dir in $BIN/conf/dstest/*/     # list directories in the form "/tmp/dirname/"
# do
#     dir=${dir%*/}      # remove the trailing "/"
#     $BIN/bin/zkServer.sh --config $dir stop
# done

rm -rf zk_conf
rm -rf /tmp/zookeeper