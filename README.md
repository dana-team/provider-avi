# Provider Avi

`provider-avi` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
Avi API.

## Getting Started

### Install the provider

```yaml
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-avi
spec:
  package: ghcr.io/dana-team/provider-avi:<release>
  runtimeConfigRef:
    apiVersion: pkg.crossplane.io/v1beta1
    kind: DeploymentRuntimeConfig
    name: config
```

```yaml
apiVersion: pkg.crossplane.io/v1beta1
kind: DeploymentRuntimeConfig
metadata:
  name: config
spec:
  deploymentTemplate:
    spec:
      selector:
        matchLabels:
          pkg.crossplane.io/provider: provider-avi
      template:
        spec:
          containers:
          - args:
            - --debug
            name: package-runtime
```

## Configuration

To connect to the provider, create the following `secret`:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: example-creds
  namespace: crossplane-system
type: Opaque
stringData:
  credentials: |
    {
      "avi_username": "<AVI-USERNAME>",
      "avi_tenant": "<AVI-TENANT>",
      "avi_password": "<AVI-PASSWORD>",
      "avi_controller": "<AVI-CONTROLLER>",
      "avi_version": "<AVI-VERSION>"
    }
```

Then create the `ProviderConfig`:

```yaml
apiVersion: avi.crossplane.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  credentials:
    source: Secret
    secretRef:
      name: example-creds
      namespace: crossplane-system
      key: credentials
```

## Resources

To Install the CRDs manually, run:

```bash
$ go install golang.org/x/tools/cmd/goimports@latest
$ make generate
$ kubectl apply -f package/crds
```
