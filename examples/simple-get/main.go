package main

import (
	"context"
	"fmt"
	"log"

	"github.com/knqyf263/go-plugin/types/known/emptypb"
	"github.com/tetratelabs/wazero"

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
	p, err := getter.NewGetterPlugin(ctx, getter.WazeroRuntime(func(ctx context.Context) (wazero.Runtime, error) {
		r, err := getter.DefaultWazeroRuntime()(ctx)
		if err != nil {
			return nil, err
		}
		return r, host.Init(ctx, r)
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
		Name: "go-plugin",
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
