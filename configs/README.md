# Example Configuration Files

This directory contains example configuration files for running dstest with different distributed systems.

## Apache Ratis

Apache Ratis is a Java library that implements the Raft consensus algorithm.
It is used to build fault-tolerant, replicated state machines to ensure strong
consistency and reliability across the replicated servers in a distributed system.

[ratis.yml](ratis.yml) showcases an example configuration file for testing Apache Ratis with 3 replicas.

After building the docker image, one can run the test with the following command from the root directory of the project:
```shell
$ docker run -v ./configs:/configs egeberkaygulcan/dstest run -c /configs/ratis.yml
```

