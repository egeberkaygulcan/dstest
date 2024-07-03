# Example Configuration Files

This directory contains example configuration files for running dstest with different distributed systems.

## Apache Ratis

Apache Ratis is a Java library that implements the Raft consensus algorithm.
It is used to build fault-tolerant, replicated state machines to ensure strong
consistency and reliability across the replicated servers in a distributed system.

To configure dstest to test Apache Ratis, you can use the provided configuration file [`config.yml`](cmd/dstest/config/config.yml).

This file contains the configuration for the Apache Ratis test, including the number of replicas, the number of interceptors, and the ports to use for the replicas and interceptors.

[ratis.yml](ratis.yml) showcases an example configuration file for testing Apache Ratis with 3 replicas.

One can run the test with the following command:
```shell
$ dstest -c configs/ratis.yml
```
