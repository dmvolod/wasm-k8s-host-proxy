package main

import (
	"context"
	"fmt"
	"log"

	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/onmetal/controller-utils/unstructuredutils"
	"github.com/tetratelabs/wazero"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/fake"

	"github.com/dmvolod/wasm-k8s-host-proxy/examples/simple-get/getter"
	"github.com/dmvolod/wasm-k8s-host-proxy/impl/host"
	"github.com/dmvolod/wasm-k8s-host-proxy/internal/unstructuredutil"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx := context.Background()

	// Instantiate Fake Kubernetes client and load test data.
	objs, err := unstructuredutils.ReadFile("testdata/configmap.yaml")
	fakeClient := fake.NewSimpleDynamicClient(runtime.NewScheme(), unstructuredutil.UnstructuredSliceToObjectSlice(objs)...)

	p, err := getter.NewGetterPlugin(ctx, getter.WazeroRuntime(func(ctx context.Context) (wazero.Runtime, error) {
		r, err := getter.DefaultWazeroRuntime()(ctx)
		if err != nil {
			return nil, err
		}

		// Register host functions from the Kubernetes library
		return r, host.Instantiate(ctx, r, func() (dynamic.Interface, error) {
			return fakeClient, nil
		})
	}))
	if err != nil {
		return err
	}

	// Pass the host functions that are embedded into the plugin.
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
