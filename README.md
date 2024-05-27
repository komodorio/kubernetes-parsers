# Kubernetes Parsers

`kubernetes-parsers` is a reusable Go module designed to host various parsing functions for Kubernetes resources. This module aims to simplify and streamline the extraction of useful information from Kubernetes resource definitions.

## Features

Currently, `kubernetes-parsers` supports the following parsing function:

- **FindStatusForPod**: This function retrieves the status of a given Pod object.

## Installation

To install the `kubernetes-parsers` module, use the following command:

```sh
go get github.com/komodorio/kubernetes-parsers
```

## Usage

Here's a basic example of how to use the `FindStatusForPod` function:

```go
package main

import (
    "fmt"
    "github.com/komodorio/kubernetes-parsers/parsers/pods"
    v1 "k8s.io/api/core/v1"
)

func main() {
    // Example Pod object
    pod := &v1.Pod{
        Status: v1.PodStatus{
            Phase: v1.PodRunning,
        },
    }

    // Use the FindStatusForPod function to get the Pod status
    status := parsers_pods.FindStatusForPod(pod)
    fmt.Println("Pod status:", status)
}
```

This example demonstrates how to import the `kubernetes-parsers` module and use the `FindStatusForPod` function to retrieve and print the status of a Pod.

## Contribution

We welcome contributions to expand the functionality of this module. If you have a parsing function you'd like to add or an improvement to suggest, please open an issue or submit a pull request on our [GitHub repository](https://github.com/yourusername/kubernetes-parsers).

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](https://github.com/yourusername/kubernetes-parsers/blob/main/LICENSE) file for details.
