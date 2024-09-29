# provider-avi

![Version: 0.0.0](https://img.shields.io/badge/Version-0.0.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: latest](https://img.shields.io/badge/AppVersion-latest-informational?style=flat-square)

A Helm chart for Crossplane provider-avi.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| deploymentRuntimeConfig | object | `{"container":{"args":["--debug"],"name":"package-runtime"},"name":"avi-config"}` | Configuration to be added to the provider deployment via the DeploymentRuntimeConfig resource |
| fullnameOverride | string | `""` |  |
| image.repository | string | `"ghcr.io/dana-team/provider-avi"` | The repository of the provider container image. |
| image.tag | string | `""` | The tag of the manager container image. |
| nameOverride | string | `""` |  |
| provider.name | string | `"provider-avi"` | Name of the provider |
| provider.runtimeConfigRef.name | string | `"avi-config"` | Name of the DeploymentRuntimeConfig object to use |
| providerConfig | object | `{"credentials":{"secretRef":{"key":"credentials","name":"avi-creds","namespace":"crossplane-system"},"source":"Secret"},"name":"avi-default"}` | Provider authentication configuration |
| secret | object | `{"controller":"127.0.0.1","name":"avi-creds","password":"passw0rd","tenant":"test","type":"Opaque","username":"dana","version":"21.1.1"}` | Secret values for the provider authentication. |
| secret.controller | string | `"127.0.0.1"` | IP of the controller to connect to. |
| secret.name | string | `"avi-creds"` | Name of the secret. |
| secret.password | string | `"passw0rd"` | Password to connect to authenticate with. |
| secret.tenant | string | `"test"` | Name of the tenant to connect to. |
| secret.type | string | `"Opaque"` | Type of the secret. |
| secret.username | string | `"dana"` | Username to connect to authenticate with. |
| secret.version | string | `"21.1.1"` | Version of the controller. |

