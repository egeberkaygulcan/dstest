# DSTest

DSTest is a Controlled Concurrency Testing Framework tool to test distributed systems without modifying the source code of the system under test.

## Prerequisites

### Running via Docker (recommended)
You'll need to [install Docker](https://docs.docker.com/get-docker/).

#### Generating the image
We are not publishing these at the moment, so you'll need to build the image yourself.
```shell
$ docker build -t egeberkaygulcan/dstest .
```

#### Running the image
This will run the image with the default configuration for Apache Ratis.
```shell
$ docker run egeberkaygulcan/dstest
```

### Running from source
You'll need to [install the Go runtime](https://go.dev/doc/install).

## Configuration
A sample configuration file is provided in [`config.yml`](cmd/dstest/config/config.yml).
You can copy this file and modify it to suit your needs.
Below is a brief explanation of the configuration options:

###### TestConfig
This section contains the general configuration for the test.
- `Name`: A human-readable name for the test.
- `Experiments`: The number of experiments to run.
- `Iterations`: The number of iterations to run per experiment.
- `WaitDuration`: The duration to wait between execution steps, in milliseconds.

###### SchedulerConfig
This section contains the configuration for the scheduler: which scheduler to use, and the parameters to pass to the scheduler.
- `Type`: The name of the scheduler to use.
- `Steps`: The number of steps to run in the scheduler.
- `Seed`: The seed to use for the random number generator.
- `Params`: A map of parameters to pass to the scheduler.

###### NetworkConfig
This section contains the configuration for the network, namely the ports to use for the replicas and their interceptors.
- `BaseReplicaPort`: The base port number to use for replicas. Each of the `N` replicas will be assigned a port number starting from this value (from `BaseReplicaPort` to `BaseReplicaPort + N - 1`).
- `BaseInterceptorPort`: The base port number to use for network interceptors. Each of the `M` interceptors will be assigned a port number starting from this value (from `BaseInterceptorPort` to `BaseInterceptorPort + M - 1`).

###### ProcessConfig
This section contains the configuration on how to spawn the processes of the system under test.
- `NumReplicas`: The number of replicas to spawn.
- `Timeout`: The timeout to wait for the system under test to finish, in seconds.
- `OutputDir`: The directory to store the output of the system under test.
- `ReplicaScript`: The script to run to start a single replica.
- `ClientScripts`: Additional scripts to run to start clients for the system under test.
- `CleanScript`: The script to run to clean up the system under test.
- `ReplicaParams`: A list of parameters to pass to the replica script; one for each replica.

## Usage
TODO

## License
See [LICENSE](LICENSE.md).
