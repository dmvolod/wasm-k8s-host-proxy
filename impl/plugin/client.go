//go:build tinygo.wasm

package kubernetes

import (
	"context"

	"github.com/mailru/easyjson"

	"github.com/dmvolod/wasm-k8s-host-proxy/internal/host/kubernetes"
)

var proxy = kubernetes.NewProxy()

type proxyClient struct {
	namespace string
	resource  GroupVersionResource
}

func NewProxyClient() Interface {
	return &proxyClient{}
}

func toGVR(gvr GroupVersionResource) *kubernetes.GVR {
	return &kubernetes.GVR{
		Group:    gvr.Group,
		Version:  gvr.Version,
		Resource: gvr.Resource,
	}
}

func (p *proxyClient) Namespace(namespace string) ResourceInterface {
	ret := *p
	ret.namespace = namespace
	return &ret
}

func (p *proxyClient) Resource(resource GroupVersionResource) NamespaceableResourceInterface {
	return &proxyClient{
		resource: resource,
	}
}

func (p *proxyClient) Get(ctx context.Context, name string, options GetOptions, object easyjson.Unmarshaler, subresources ...string) error {
	reply, err := proxy.Get(ctx, &kubernetes.GetRequest{
		Gvr: toGVR(p.resource),
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
