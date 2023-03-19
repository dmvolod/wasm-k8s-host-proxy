package main

import (
	"context"
	"fmt"
	"log"

	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/tetratelabs/wazero"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"

	"github.com/dmvolod/wasm-k8s-host-proxy/examples/simple-get/getter"
	"github.com/dmvolod/wasm-k8s-host-proxy/impl/host"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// Instantiate Fake Kubernetes client for demo
	fakeClient := fake.NewSimpleDynamicClient(runtime.NewScheme())
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo",
			Namespace: "default",
		},
		Data: map[string]string{
			"my-data": "This is test data",
		},
	}

	cmunstr, err := runtime.DefaultUnstructuredConverter.ToUnstructured(cm)
	if err != nil {
		return err
	}

	_, err = fakeClient.Resource(schema.GroupVersionResource{
		Version:  "v1",
		Resource: "configmaps",
	}).Namespace(cm.Namespace).Create(ctx, &unstructured.Unstructured{
		Object: cmunstr,
	}, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	p, err := getter.NewGetterPlugin(ctx, getter.WazeroRuntime(func(ctx context.Context) (wazero.Runtime, error) {
		r, err := getter.DefaultWazeroRuntime()(ctx)
		if err != nil {
			return nil, err
		}

		return r, host.Instantiate(ctx, r, func() (dynamic.Interface, error) {
			return fakeClient, nil
		})
	}))
	if err != nil {
		return err
	}

	// Pass my host functions that are embedded into the plugin.
	getterPlugin, err := p.Load(ctx, "plugin/plugin.wasm", getterHostFunctions{})
	if err != nil {
		return err
	}
	defer getterPlugin.Close(ctx)

	reply, err := getterPlugin.GetConfigMap(ctx, &getter.GetRequest{
		Name:      "demo",
		Namespace: "default",
	})
	if err != nil {
		return err
	}

	fmt.Println(reply.Data)

	return nil
}

// getterHostFunctions implements getter.HostFunctions
type getterHostFunctions struct{}

var _ getter.HostFunctions = (*getterHostFunctions)(nil)

// Log is embedded into the plugin and can be called by the plugin.
func (getterHostFunctions) Log(_ context.Context, request *getter.LogRequest) (*emptypb.Empty, error) {
	// Use the host logger
	log.Println(request.GetMessage())
	return &emptypb.Empty{}, nil
}
