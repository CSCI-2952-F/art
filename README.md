# ART: Automated Resilience Testing for microservices

ART is a framework for automated resilience testing on microservices deployed on top of Istio. Our purpose is to provide meaningful insights of application behavior during failure and improve ease-of-use of microservice fault-tolerance testing. ART accomplishes these goals by identifying critical application components, automatically conducting fault-injection experiments, and presenting the results to the user.


[Full Paper](./static/paper.pdf)

## Usage

### Build
Run ```make``` at the top level directory to build the client binary.

Dependencies:
- python3/pip3
- go 1.14

If ```make``` is successful, you can run the client binary with ```./bin/client```. By default, the client library
will connect to a test cluster we already have running so all tracing and fault injection functionality
should be already set up.

Shell commands:
- ```experiment``` - Runs a full experiment as described from the paper
- ```continue``` - Continues an experiment in case of temporary failures
- ```kiali``` - Uses kiali instead of sub-graph mining to conduct an experiment

The sections below are for deploying ART on your own application cluster. You can disregard them if you are
simply demo-ing ART with our test cluster.

### Fault Service
ART requires the fault service be deployed in the application cluster.
This allows for the ART client to apply fault injections on the cluster without having the necessary
Kuberentes certificates. We have provided yaml files in ```kubernetes-manifests/``` to help you deploy
the fault service.

### Client
The client is an interactive shell, which you can start with the following command:
```shell script
$ client --jaeger <jaeger-addr> --fault <fault-addr> --kiali <kiali-addr>
```
Flags:
- jaeger: Address of jaeger component where you can query traces from
- fault: Address of fault service deployed in the cluster
- kiali: Address of kiali