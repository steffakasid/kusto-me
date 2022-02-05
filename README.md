# kusto-me 
[![Go Build](https://github.com/steffakasid/kusto-me/actions/workflows/go-build.yml/badge.svg)](https://github.com/steffakasid/kusto-me/actions/workflows/go-build.yml) [![Go Test](https://github.com/steffakasid/kusto-me/actions/workflows/go-test.yml/badge.svg)](https://github.com/steffakasid/kusto-me/actions/workflows/go-test.yml)


kusto-me (kustomize me) can be used to initalize folders with k8s objects with a kustomize.yaml and optional with a overlay folder structure.

## Usage of kusto-me:
```
  -d, --dir stringArray     Define overlay directories to be used.
  -f, --folder string       Set the folder to create kustomize project. If not set current dir is used.
  -l, --label stringArray   Add default labels to kustomization.yml. Format: <key>:<value>
  -n, --name string         Set the projectname. If not set the current directoryname is used
  -o, --overlay             Defines if overlay structure should be created or just a simple project
```

## Install

On OSX/ Linux you may just use brew to install kusto-me:
```
 brew tap steffakasid/kusto-me
 brew install kustome
```