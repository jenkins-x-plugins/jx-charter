# jx-charter

[![Documentation](https://godoc.org/github.com/jenkins-x-plugins/jx-charter?status.svg)](https://pkg.go.dev/mod/github.com/jenkins-x-plugins/jx-charter)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenkins-x-plugins/jx-charter)](https://goreportcard.com/report/github.com/jenkins-x-plugins/jx-charter)
[![Releases](https://img.shields.io/github/release-pre/jenkins-x-plugins/jx-charter.svg)](https://github.com/jenkins-x-plugins/jx-charter/releases)
[![Apache](https://img.shields.io/badge/license-Apache-blue.svg)](https://github.com/jenkins-x-plugins/jx-charter/blob/master/LICENSE)
[![Slack Status](https://img.shields.io/badge/slack-join_chat-white.svg?logo=slack&style=social)](https://slack.k8s.io/)

`jx-charter` is a small command line tool for creating Helm `Chart` CRDs from helm releases for better metadata and reporting of what is running inside kubernetes

## Commands

See the [jx-charter command reference](https://github.com/jenkins-x-plugins/jx-charter/blob/master/docs/cmd/jx-charter.md)


## Installation

To install the chart use the following:


- Add jx3 helm charts repo

```bash
helm repo add jx3 https://storage.googleapis.com/jenkinsxio/charts

helm repo update
```

- Install (or upgrade)

```bash
# This will install or upgrade the jx-charter chart in the current namespace (with a jx-charter release name)

helm upgrade --install jx-charter jx3/jx-charter
```

## Uninstalling

To uninstall the chart, simply delete the release.

```bash
# This will uninstall jx-charter in the current namespace (assuming a jx-charter release name)

# Helm v3
helm uninstall jx-charter
```
