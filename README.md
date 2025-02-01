# kubectl-eks-viewer

A kubectl plugin that provides a convenient way to view AWS EKS cluster resources. This plugin allows you to quickly inspect various EKS resources like nodegroups, Fargate profiles, addons, and more, directly from your kubectl environment.

## Features

- Follow kubectl output format
- View multiple EKS resource types in one command
- View specific resource types individually
- Automatic EKS cluster detection from current kubectl context
- Supported resources:
  - Access Entries
  - Addons
  - Cluster Information
  - Fargate Profiles
  - Insights
  - Nodegroups
  - Pod Identity Associations

## Prerequisites

- kubectl installed and configured
- AWS credentials configured with appropriate EKS permissions
- [Krew](https://krew.sigs.k8s.io/docs/user-guide/setup/install/) (for installing via Krew)

## Installation

### Using Krew (Recommended)

[Krew](https://krew.sigs.k8s.io) is the plugin manager for kubectl command-line tool. To install kubectl-eks-viewer using Krew:

```bash
kubectl krew install eks-viewer
```

After installation, verify it was successful:
```bash
kubectl eks-viewer --help
```

## Usage

```bash
Usage:
  kubectl eks-viewer [resource-type] [flags]

Examples:
  # List all EKS resources
  kubectl eks-viewer
  kubectl eks-viewer -o json

  # List specific resources
  kubectl eks-viewer addons
  kubectl eks-viewer -o json nodegroups
  kubectl eks-viewer nodegroups --output=jsonpath='{.items.nodegroups[*].NodegroupName}

  # Use with a specific context
  kubectl eks-viewer --context=my-context
```

## Available Resource Types

- `access-entries`: View EKS cluster access entries
- `addons`: List installed EKS addons
- `cluster`: Show cluster information
- `fargate-profiles`: Display Fargate profiles
- `insights`: View cluster insights
- `nodegroups`: List managed node groups
- `pod-identity-associations`: Show pod identity associations

## Feature requests & bug reports

If you have any feature requests or bug reports, please submit them through GitHub [Issues](https://github.com/keidarcy/kubectl-eks-viewer/issues).

## Contact & Author

Xing Yahao(https://github.com/keidarcy)