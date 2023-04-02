//go:build tinygo.wasm

package client

import (
	"context"
	"strings"

	"github.com/kubewarden/k8s-objects/apimachinery/pkg/runtime/schema"
	"github.com/mailru/easyjson"

	"github.com/dmvolod/wasm-k8s-host-proxy/internal/host/kubernetes"
)

var proxy = kubernetes.NewProxy()

type proxyClient struct {
	namespace string
}

func NewProxyClient() NamespaceableResourceInterface {
	return &proxyClient{}
}

func toGVK(gvr schema.GroupVersionKind) *kubernetes.GVR {
	return &kubernetes.GVR{
		Group:    gvr.Group,
		Version:  gvr.Version,
		Resource: strings.ToLower(gvr.Kind) + "s",
	}
}

func (p *proxyClient) Namespace(namespace string) ResourceInterface {
	ret := *p
	ret.namespace = namespace
	return &ret
}

func (p *proxyClient) Get(ctx context.Context, name string, options GetOptions, object Object, subresources ...string) error {
	reply, err := proxy.Get(ctx, &kubernetes.GetRequest{
		Gvr: toGVK(object.GroupVersionKind()),
		Name: &kubernetes.Namespaced{
			Name:      name,
			Namespace: p.namespace,
		},
	})
	if err != nil {
		return err
	}

	return easyjson.Unmarshal(reply.Payload, object)
}
