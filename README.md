## Lesson Structure

1. Ask about how many people code on a day-to-day basis
   a) For those that do ask which language they use
2. Explain why Golang is a important language to learn nowadays especially for CNCF tools
   a) The tools are all written in Go
   b) Because of the former point, you can hop into OS codebases to gain understanding of how the software works
   c) If there are any bugs, you can even fork the project yourself and fix the bugs yourself. That also can be contributed upstream (brownie points)
3. Explain how learning how to write code as a DevOps gives one so much power, and that you highly recommend they should learn
4. Give brief overview of k8s operators and how they work
   a) Kubernetes way of extending their API and functionality to developers
   b) You can specify `CustomResources`, and controllers to manage them for reconciling whatever state you want
   c) Whenever you create/delete something on k8s API server, controllers try to reconcile the state of the objects to the "reality of the world" (can be anything) _Give an example of the aws-lb-contrller_
5. Show that operators are actually a leveraged pattern in native k8s
   a) Show this example of the `Deployment` controller (https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/deployment_controller.go)
6. Show the `operator-sdk` docs and how it makes it pretty easy to create a custom controller with your custom resource
7. Start to show the `git-repo-operator` and how I wrote a custom one using the `operator-sdk`
   a) Interfaces with the GitHub API, and creates/destroys repositories

## Getting Started

### Prerequisites

- go version v1.21.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/git-repo-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/git-repo-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
> privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

> **NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall

**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/git-repo-operator:tag
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/git-repo-operator/<tag or branch>/dist/install.yaml
```

## Contributing

// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
