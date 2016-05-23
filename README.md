# git-controller

git-controller watches Deployments or ReplicationController resources inside kubernetes which use the `gitRepo` volume and then watches for changes in the git branch. When there are changes it'll automatically perform a rolling upgrade to use the new configuration; or rollback.

## Development

### Prerequisites

Install [go version 1.5.1](https://golang.org/doc/install)
Install [godep](https://github.com/tools/godep)


### Building

```sh
git clone git@github.com:fabric8io/git-controller.git $GOPATH/src/github.com/fabric8io/git-controller
./make
```

Make changes to *.go files, rerun `make` and run the generated binary..

e.g.

```sh
./build/git-controller help

```
