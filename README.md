# DSTest

DSTest is a Controlled Concurrency Testing Framework tool to test distributed systems without modifying the source code of the system under test.

## Prerequisites

> [!NOTE]
> Code was tested on macOS on arm64 architecture. Other platforms may require additional setup.

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
- `StartupDuration`: The duration to wait for the system under test to start up before scheduling the first step, in seconds.

###### SchedulerConfig
This section contains the configuration for the scheduler: which scheduler to use, and the parameters to pass to the scheduler.
- `Type`: The name of the scheduler to use. Possible values are `Random`, `QL`, and `PCTCP`.
- `Steps`: The number of steps to run in the scheduler.
- `ClientRequests`: The number of client requests to generate per experiment.
- `Seed`: The seed to use for the random number generator.
- `Params`: A map of parameters to pass to the scheduler. Each scheduler has its own set of parameters.

###### FaultConfig
This section contains the configuration for the fault injector.
- `Faults`: An array of faults to inject. Each fault has the following fields:
  - `Type`: The type (ID) of fault to inject.
  - `Params`: A map of parameters to pass to the fault. Each fault type has its own set of parameters.

###### NetworkConfig
This section contains the configuration for the network, namely the ports to use for the replicas and their interceptors.
- `BaseReplicaPort`: The base port number to use for replicas. Each of the `N` replicas will be assigned a port number starting from this value (from `BaseReplicaPort` to `BaseReplicaPort + N - 1`).
- `BaseInterceptorPort`: The base port number to use for network interceptors. Each of the `M` interceptors will be assigned a port number starting from this value (from `BaseInterceptorPort` to `BaseInterceptorPort + M - 1`).
- `Protocol`: The protocol to use for the network. Possible values are `http` and `tcp`.
- `MessageType`: The message type to use for the network. Just `GRPC` is supported for now.

###### ProcessConfig
This section contains the configuration on how to spawn the processes of the system under test.
- `NumReplicas`: The number of replicas to spawn.
- `Timeout`: The timeout to wait for the system under test to finish, in seconds.
- `OutputDir`: The directory to store the output of the system under test.
- `ReplicaScript`: The script to run to start a single replica.
- `ClientScripts`: Additional scripts to run to start clients for the system under test.
- `CleanScript`: The script to run to clean up the system under test.
- `ReplicaParams`: A list of parameters to pass to the replica script; one for each replica.

## Running the examples

See the [configs](configs/README.md) directory for examples of how to run DSTest with different distributed systems and sample configurations.

## License
See [LICENSE](LICENSE.md).
