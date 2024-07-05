#!/bin/bash
BIN=../zookeeper
mkdir -p $BIN/conf/dstest/zoo_$9
printf "tickTime=2000\ndataDir=$8/zk_$9\nclientPort=$7\ninitLimit=5\nsyncLimit=2\nserver.1=localhost:$1:$2\nserver.2=localhost:$3:$4\nserver.3=localhost:$5:$6" >> ${BIN}/conf/dstest/zoo_$9/zoo.cfg
mkdir -p $8/zk_$9
printf $9 >> $8/zk_$9/myid
bash ${BIN}/bin/zkServer.sh --config ${BIN}/conf/dstest/zoo_$9 start
