# kubectl-eks-viewer

A kubectl plugin that provides a convenient way to view AWS EKS cluster resources. This plugin allows you to quickly inspect various EKS resources like nodegroups, Fargate profiles, addons, and more, directly from your kubectl environment.

Read more about [Why I created this plugin](https://xingyahao.com/posts/introducing-kubectl-eks-viewer-bridge-the-gap-between-kubectl-and-aws-cli/).

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

You can add custom index(since krew only accept vendor plugins that come from the vendors) as shown below and install the plugin from there. To install kubectl-eks-viewer using Krew:

```bash
kubectl krew index add keidarcy https://github.com/keidarcy/kubectl-eks-viewer.git
kubectl krew install keidarcy/eks-viewer
```

Details:

```
$ kubectl krew index add keidarcy https://github.com/keidarcy/kubectl-eks-viewer.git
WARNING: You have added a new index from "https://github.com/keidarcy/kubectl-eks-viewer.git"
The plugins in this index are not audited for security by the Krew maintainers.
Install them at your own risk.

$ kubectl krew install  keidarcy/eks-viewer
Updated the local copy of plugin index.
Updated the local copy of plugin index "keidarcy".
Installing plugin: eks-viewer
Installed plugin: eks-viewer
\
 | Use this plugin:
 |      kubectl eks-viewer
 | Documentation:
 |      https://github.com/keidarcy/kubectl-eks-viewer
/

$ kubectl eks-viewer --help
View EKS cluster resources.
Without arguments, shows all resource types.
Optionally specify a resource type to show only that type.

Valid resource types:
  - access-entries
  - addons
  - cluster
  - fargate-profiles
  - insights
  - nodegroups
  - pod-identity-associations

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

Flags:
      --allow-missing-template-keys    If true, ignore any errors in templates when a field or map key is missing in the template. Only applies to golang and jsonpath output formats. (default true)
      --as string                      Username to impersonate for the operation. User could be a regular user or a service account in a namespace.
      --as-group stringArray           Group to impersonate for the operation, this flag can be repeated to specify multiple groups.
      --as-uid string                  UID to impersonate for the operation.
      --cache-dir string               Default cache directory (default "/Users/xyh/.kube/cache")
      --certificate-authority string   Path to a cert file for the certificate authority
      --client-certificate string      Path to a client certificate file for TLS
      --client-key string              Path to a client key file for TLS
      --cluster string                 The name of the kubeconfig cluster to use
      --context string                 The name of the kubeconfig context to use
      --disable-compression            If true, opt-out of response compression for all requests to the server
  -h, --help                           help for kubectl
      --insecure-skip-tls-verify       If true, the server's certificate will not be checked for validity. This will make your HTTPS connections insecure
      --kubeconfig string              Path to the kubeconfig file to use for CLI requests.
  -n, --namespace string               If present, the namespace scope for this CLI request
  -o, --output string                  Output format. One of: (json, yaml, name, go-template, go-template-file, template, templatefile, jsonpath, jsonpath-as-json, jsonpath-file).
      --request-timeout string         The length of time to wait before giving up on a single server request. Non-zero values should contain a corresponding time unit (e.g. 1s, 2m, 3h). A value of zero means don't timeout requests. (default "0")
  -s, --server string                  The address and port of the Kubernetes API server
      --show-managed-fields            If true, keep the managedFields when printing objects in JSON or YAML format.
      --template string                Template string or path to template file to use when -o=go-template, -o=go-template-file. The template format is golang templates [http://golang.org/pkg/text/template/#pkg-overview].
      --tls-server-name string         Server name to use for server certificate validation. If it is not provided, the hostname used to contact the server is used
      --token string                   Bearer token for authentication to the API server
      --user string                    The name of the kubeconfig user to use

$ kubectl eks-viewer
=== cluster ===
NAME   VERSION   STATUS   PLATFORM VERSION   AUTH MODE
<REDACTED>

=== access-entries ===
ACCESS ENTRY PRINCIPAL ARN    KUBERNETES GROUPS    ACCESS POLICIES
<REDACTED>

=== addons ===
NAME    VERSION    STATUS   ISSUES
<REDACTED>

=== nodegroups ===
NAME      STATUS   INSTANCE TYPE   DESIRED SIZE   MIN SIZE   MAX SIZE   VERSION   AMI TYPE     CAPACITY TYPE
<REDACTED>

=== fargate-profiles ===
NAME   SELECTOR NAMESPACE   SELECTOR LABELS   POD EXECUTION ROLE ARN     SUBNETS     STATUS
<REDACTED>

=== pod-identity-associations ===
ARN    AMESPACE   SERVICE ACCOUNT NAME   IAM ROLE ARN      OWNER ARN
<REDACTED>

=== insights ===
NAME       CATEGORY    STATUS
<REDACTED>

$ kubectl krew upgrade eks-viewer
Updated the local copy of plugin index.
Updated the local copy of plugin index "keidarcy".
Upgrading plugin: keidarcy/eks-viewer
Upgraded plugin: keidarcy/eks-viewer
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