//go:build tinygo.wasm

package main

import (
	"context"
	"fmt"

	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"

	"github.com/dmvolod/wasm-k8s-host-proxy/examples/simple-get/getter"
	kubernetes "github.com/dmvolod/wasm-k8s-host-proxy/impl/plugin"
)

// main is required for TinyGo to compile to Wasm.
func main() {
	getter.RegisterGetter(GetterPlugin{})
}

type GetterPlugin struct{}

var _ getter.Getter = (*GetterPlugin)(nil)

func (g GetterPlugin) GetConfigMap(ctx context.Context, request *getter.GetRequest) (*getter.GetReply, error) {
	hostFunctions := getter.NewHostFunctions()
	kubeClientProxy := kubernetes.NewProxyClient()
	hostFunctions.Log(ctx, &getter.LogRequest{
		Message: fmt.Sprintf("Getting ConfigMap object '%s/%s' from the Kubernetes cluster...", request.Name, request.Namespace),
	})

	cm := &corev1.ConfigMap{}
	cmGVR := kubernetes.GroupVersionResource{
		Version:  "v1",
		Resource: "configmaps",
	}

	err := kubeClientProxy.Resource(cmGVR).Namespace(request.Namespace).Get(ctx, request.Name, kubernetes.GetOptions{}, cm)
	if err != nil {
		return nil, err
	}

	return &getter.GetReply{
		Data: cm.Kind,
	}, nil
}
