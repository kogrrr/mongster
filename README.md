# Mongster - Single-Page App using Golang / Vue.js / Bootstrap / MongoDB

---
* [Developer Setup](#developer-setup)
  * [Prerequisites](#prerequisites)
  * [Getting Started](#getting-started)
  * [Building](#building)
* [Running](#running)
---


## Developer Setup

### Prerequisites

You will need a few things in order to get started:

* bash 4 or higher
* GNU make 3.81+
* a sensible version of Golang with go mod support (1.12+)
* [golangci-lint](https://github.com/golangci/golangci-lint#install) to perform backend code linting
* [ginkgo](https://onsi.github.io/ginkgo/) to run backend tests
* a sane version of `npm` (if there is such a thing)
* the `parcel` bundler

You will also need a MongoDB instance running to act as backend storage.

The easiest way to do this is using Docker:
```
docker run -ti -p 27017:27017 mongo:latest
```

*Note*: without also mounting a storage volume, the above docker command will wipe the database each time the
container is restarted.


### Getting Started

#### Ensure your Go installation is correct

If you haven't already, [install Go](https://golang.org/doc/install) and decide where to
keep your `GOPATH`. A sensible location might be `~/go`

Next, clone this repository into your `GOPATH`
```
$ cd $GOPATH
$ mkdir -p src/github.com/gargath/
$ cd src/github.com/gargath/
$ git clone <this repository>
$ cd mongster
```

Next, run `./configure` to validate and set up your development. If you are missing any
of the prerequisites, the script will inform you.

Once `configure` completes successfully, you are ready to build.


### Building

Mongster can be built as a self-contained binary that will serve both backend and frontend
resources and does not rely on any external files.
It is thus ideal for packaging in containers.

To build the binary:
```
$ make clean build
```

You can also run the application in-place and serve the frontend assets from the filesytem.
This is useful during development and includes a hot-reload feature.

To start dev mode:
```
$ make dev
```

To exit dev mode, simply issue `Ctrl-C` in the same terminal.


## Running

Once built, the binary can be run with configuration from either the environment or command-line
parameters:

```
$ ./mongster --mongsterConnstr mongodb://localhost:27017
```
is equivalent to
```
$ MONGSTER_MONGOCONNSTR=mongodb://localhost:27017 ./mongster
```

The available configuration options can be displayed with `./mongster -help`

