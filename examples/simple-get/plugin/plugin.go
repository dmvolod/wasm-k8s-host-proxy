//go:build tinygo.wasm

package main

import (
	"context"
	"fmt"

	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"

	"github.com/dmvolod/wasm-k8s-host-proxy/examples/simple-get/getter"
	"github.com/dmvolod/wasm-k8s-host-proxy/impl/plugin/client"
)

// main is required for TinyGo to compile to Wasm.
func main() {
	getter.RegisterGetter(GetterPlugin{})
}

type GetterPlugin struct{}

var _ getter.Getter = (*GetterPlugin)(nil)

func (g GetterPlugin) GetConfigMap(ctx context.Context, request *getter.GetRequest) (*getter.GetReply, error) {
	hostFunctions := getter.NewHostFunctions()
	kubeClientProxy := client.NewProxyClient()
	hostFunctions.Log(ctx, &getter.LogRequest{
		Message: fmt.Sprintf("Getting ConfigMap object '%s/%s' from the Kubernetes cluster...", request.Name, request.Namespace),
	})

	cm := &corev1.ConfigMap{}
	err := kubeClientProxy.Namespace(request.Namespace).Get(ctx, request.Name, client.GetOptions{}, cm)
	if err != nil {
		return nil, err
	}

	return &getter.GetReply{
		Data: cm.Data["my-data"],
	}, nil
}
