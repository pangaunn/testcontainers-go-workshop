# Testcontainers-go workshop

## Prerequisite
- installed golang 
- [ginkgo](https://onsi.github.io/ginkgo/MIGRATING_TO_V2)

## Before workshop please run
this command will pull docker image that we use in this workshop
```
$ make int-test
```

## Running on local
copy `.env.example` => `.env`
```
$ make run-docker
```

## Running a test
```
$ make int-test
```