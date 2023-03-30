# Simple Get Kubernetes Object Example
This example shows how to organize the simple interaction with the Kubernetes cluster from the plugin code.

## Generate Go code
A proto file is under `./getter`.

```shell
$ protoc --go-plugin_out=. --go-plugin_opt=paths=source_relative getter/getter.proto 
```

## Compile a plugin
Use TinyGo to compile the plugin to Wasm.

```shell
$ go generate ./...
```

## Run
`main.go` loads plugin and interaction with fake Kubernetes cluster.

```shell
$ go run main.go
```
