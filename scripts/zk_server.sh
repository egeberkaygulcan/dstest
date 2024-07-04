#!/bin/bash
BIN=../zookeeper
printf "tickTime=2000\ndataDir=$4\nclientPort=2181\ninitLimit=5\nsyncLimit=2\nserver.1=localhost:2888:$1\nserver.2=localhost:2889:$2\nserver.3=localhost:2890:$3" >> ${BIN}/conf/zoo.cfg
bash ${BIN}/bin/zkServer.sh start
