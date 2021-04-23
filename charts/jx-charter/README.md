

# jx charter

This chart bootstraps installation of Jenkins X previews.

## Installing

- Add jx3 helm charts repo

```bash
helm repo add jx3 https://storage.googleapis.com/jenkinsxio/charts

helm repo update
```

- Install (or upgrade)

```bash
# This will install jx-charter in the jx namespace (with a jx-charter release name)

# Helm v3
helm upgrade --install jx-charter --namespace jx jx3/jx-charter
```

Look [below](#values) for the list of all available options and their corresponding description.

## Uninstalling

To uninstall the chart, simply delete the release.

```bash
# This will uninstall jx-charter in the jx-charter namespace (assuming a jx-charter release name)

# Helm v3
helm uninstall jx-charter --namespace jx
```

## Version

