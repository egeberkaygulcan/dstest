# Example Configuration Files

This directory contains example configuration files for running dstest with different distributed systems.

> [!IMPORTANT]
> - Commands in this document are run from the root directory of the project.
> - Make sure to build the docker image before running the tests.

## Apache Ratis

Apache Ratis is a Java library that implements the Raft consensus algorithm.
It is used to build fault-tolerant, replicated state machines to ensure strong
consistency and reliability across the replicated servers in a distributed system.

[ratis.yml](ratis.yml) showcases an example configuration file for testing Apache Ratis with 3 replicas.
This configuration references these [startup](../scripts/start-ratis.sh) and [cleanup](../scripts/cleanup-ratis.sh) scripts used in the configuration.

After building the docker image, one can run the test with the following command from the root directory of the project:
```shell
$ docker run -v ./configs:/configs egeberkaygulcan/dstest run -c /configs/ratis.yml
```

## Apache ZooKeeper

Apache ZooKeeper is a distributed coordination service that provides a hierarchical key-value store for managing
configuration information, naming, providing distributed synchronization, and providing group services.
It is a core component of many distributed systems and is used to ensure consistency and reliability across the
replicated servers in a distributed system.

[zookeeper.yml](zookeeper.yml) showcases an example configuration file for testing Apache ZooKeeper with 3 replicas.
This configuration references these [startup](../scripts/start-zookeeper.sh) and [cleanup](../scripts/cleanup-zookeeper.sh) scripts used in the configuration.

After building the docker image, one can run the test with the following command from the root directory of the project:
```shell
$ docker run -v ./configs:/configs egeberkaygulcan/dstest run -c /configs/zookeeper.yml
```

## etcd

etcd is a distributed key-value store that provides a reliable way to store data across a cluster of machines.

[etcd.yml](etcd.yml) showcases an example configuration file for testing etcd with 3 replicas.
This configuration references these [startup](../scripts/start-etcd.sh) and [cleanup](../scripts/cleanup-etcd.sh) scripts used in the configuration.

After building the docker image, one can run the test with the following command from the root directory of the project:
```shell
$ docker run -v ./configs:/configs egeberkaygulcan/dstest run -c /configs/etcd.yml
```
