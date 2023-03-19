//go:build !tinygo.wasm

package host

import (
	"context"
	"encoding/json"
	"os"
	"path"

	"github.com/tetratelabs/wazero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/dmvolod/wasm-k8s-host-proxy/internal/host/kubernetes"
)

var _ kubernetes.Proxy = (*kubernetesProxy)(nil)

type kubernetesProxy struct {
	dynamicClient dynamic.Interface
}

type ConfigOption func() (dynamic.Interface, error)

func WithDefaultKubeConfig() ConfigOption {
	home, err := os.UserHomeDir()
	if err != nil {
		return func() (dynamic.Interface, error) {
			return nil, err
		}

	}

	return WithKubeConfig(path.Join(home, ".kube", "config"))
}

func WithKubeConfig(kubeConfig string) ConfigOption {
	restConfig, err := restConfig(kubeConfig)
	if err != nil {
		return func() (dynamic.Interface, error) {
			return nil, err
		}
	}

	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return func() (dynamic.Interface, error) {
			return nil, err
		}
	}

	return func() (dynamic.Interface, error) {
		return dynamicClient, nil
	}
}

func Instantiate(ctx context.Context, runtime wazero.Runtime, option ConfigOption) error {
	kubernetesClient, err := option()
	if err != nil {
		return err
	}
	return kubernetes.Instantiate(ctx, runtime, kubernetesProxy{kubernetesClient})
}

func restConfig(kubeConfig string) (*rest.Config, error) {
	if kubeConfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeConfig)
	}

	return rest.InClusterConfig()
}

func (p kubernetesProxy) Get(ctx context.Context, request *kubernetes.GetRequest) (*kubernetes.GetReply, error) {
	gvr := schema.GroupVersionResource{
		Group:    request.Gvr.Group,
		Version:  request.Gvr.Version,
		Resource: request.Gvr.Resource,
	}
	res, err := p.dynamicClient.Resource(gvr).Namespace(request.Name.Namespace).Get(ctx, request.Name.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	payload, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return &kubernetes.GetReply{Payload: payload}, nil
}
