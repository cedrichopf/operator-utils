# Operator SDK Utils

---

- [Operator SDK Utils](#operator-sdk-utils)
  - [Requirements](#requirements)
  - [Usage](#usage)
    - [Utils Package](#utils-package)
    - [Hash Package](#hash-package)
    - [Network Package](#network-package)
  - [License](#license)

---

This repository contains a Go module with utilities to build Kubernetes Operators using the
[Operator SDK](https://sdk.operatorframework.io/).

## Requirements

- Operator SDK
- Go (Version >=1.20)

## Usage

Install the Go module:

```sh
go get github.com/cedrichopf/operator-utils
```

### Utils Package

The utils package contains reconcile functions for many Kubernetes objects. While reconciling, the given
object will be either created or updated.

Additionally, the function will add a revision hash of the current object configuration as a label to the
object. Once the function is called again, the current active object can be easily compared with the expected
object using the provided revision hash.

Example:

```go
func() {
  service := &corev1.Service{
    ObjectMeta: metav1.ObjectMeta{
      Name:      "example-service",
      Namespace: "default",
    },
    Spec: corev1.ServiceSpec{
      Ports: []corev1.ServicePort{
        {
          Port: 3000,
        },
      },
    },
  }

  result, err := utils.ReconcileConfigMap(ctx, service, owner, r.Client, r.Scheme)
  if err != nil {
    log.Println(err)
  }

  if result.Updated {
    log.Println("Object has been updated")
  }
}
```

### Hash Package

The hash package contains functions to generate a revision hash and the related label for a given object.

### Network Package

The network package contains functions to check and verify network functionalities. E.g. checking if a host
is resolvable over DNS.

## License

This project is licensed under Apache-2.0 license
