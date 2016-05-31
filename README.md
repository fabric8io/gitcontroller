# git-controller

git-controller watches [Kubernetes Deployments](http://kubernetes.io/docs/user-guide/deployments/) which use one or more [`gitRepo` volumes](http://kubernetes.io/docs/user-guide/volumes/#gitrepo) and then watches for changes in the associated git repository and branch.

When there are changes in git the `gitcontroller` will perform a rolling upgrade of the [Kubernetes Deployments](http://kubernetes.io/docs/user-guide/deployments/) to use the new configuration git repository revision; or rollback. The rolling upgrade policy (e.g. speed and number of pods which update and so forth) is all specified by your [rolling update configuration in the Deployment specification](http://kubernetes.io/docs/user-guide/deployments/#rolling-update-deployment).

Here is an [example of how to add a `gitRepo` volume to your application](https://github.com/jstrachan/springboot-config-demo/blob/master/src/main/fabric8/deployment.yml#L5-L14); in this case a spring boot application to load the [`application.properties`](https://github.com/jstrachan/sample-springboot-config/blob/master/application.properties) file from a git repository.

You can either run `gitcontroller` as a microservice in your namespace; its particularly useful at development time. Or you can use the `gitcontroller` binary at any time or as part of your [CI / CD Pipeline](http://fabric8.io/guide/cdelivery.html) process.

**Note** we recommend using separate git based configuration only for things which truly are environment specific. Its simpler to include all other configuration data with your microservice source code and then use a more regular [CI / CD Pipeline](http://fabric8.io/guide/cdelivery.html) from a single git repository to build your code, create the configuration files and package it all into an immutable docker image.

## Using gitcontroller as a command

To use `gitcontroller` as a command, such as in a CI / CD pipeline use the following:

```sh
git-controller check

```

This will check all [Deployments](http://kubernetes.io/docs/user-guide/deployments/) in the current namespace for  [`gitRepo` volumes](http://kubernetes.io/docs/user-guide/volumes/#gitrepo) and perform rolling upgrades if they have changed.

You can specify a label selector expression via the `--selector` command line option.

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
