# wasm-k8s-host-proxy

### Generating proxy SDK stubs

```bash
protoc --go-plugin_out=. --go-plugin_opt=paths=source_relative internal/host/kubernetes/kubernetes.proto
```