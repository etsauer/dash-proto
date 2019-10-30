# Dash Prototype

[![Actions Status](https://github.com/redhat-cop/dash/workflows/Go/badge.svg)](https://github.com/redhat-cop/dash/actions)

Dash is an automation framework for Kubernetes. The aim is to combine the power of [Kubernetes declarative management](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/declarative-config/) and popular templating engines to provide a full-featured automation strategy for Kubernetes cluster and application management.

Dash helps to enable GitOps and Infrastructure as Code workflows in enterprises with complex use cases.

## Features

### Template Support

Currently Dash supports the following templates:

* FileTemplates - not really templates, just static resource definitions
* [HelmCharts](./docs/helm.md)
* [OpenShift Templates](./docs/openshift_templates.md)

## Installation

Those with a go environment can install simply by running

```
go get github.com/redhat-cop/dash
```

If you want a stable release, we recommend downloading the appropriate binary from our [releases](https://github.com/redhat-cop/dash/releases) page.

## Quickstart

To run, first log into your kubernetes cluster, and:

```
dash -i examples/default/
```

## Contribute

We would love for you to contribut to this project! For more information on how to get started, see our [dev guide](./docs/dev_guide.md).
