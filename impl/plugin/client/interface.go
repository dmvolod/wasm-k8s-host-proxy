// Kubernetes client-go interface and functions ported for TinyGo.
// Copyright 2016 The Kubernetes Authors.

//go:build tinygo.wasm

package client

import (
	"context"

	"github.com/kubewarden/k8s-objects/apimachinery/pkg/runtime/schema"
	"github.com/mailru/easyjson"
)

type Object interface {
	schema.ObjectKind
	easyjson.Unmarshaler
}

type TypeMeta struct {
	// Kind is a string value representing the REST resource this object represents.
	// Servers may infer this from the endpoint the client submits requests to.
	// Cannot be updated.
	// In CamelCase.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`

	// APIVersion defines the versioned schema of this representation of an object.
	// Servers should convert recognized schemas to the latest internal value, and
	// may reject unrecognized values.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#resources
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`
}

type GetOptions struct {
	TypeMeta `json:",inline"`
	// resourceVersion sets a constraint on what resource versions a request may be served from.
	// See https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions for
	// details.
	//
	// Defaults to unset
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,1,opt,name=resourceVersion"`
	// +k8s:deprecated=includeUninitialized,protobuf=2
}

type NamespaceableResourceInterface interface {
	Namespace(string) ResourceInterface
	ResourceInterface
}

type ResourceInterface interface {
	Get(ctx context.Context, name string, options GetOptions, object Object, subresources ...string) error
}
