// Code generated by protoc-gen-go-plugin. DO NOT EDIT.
// versions:
// 	protoc-gen-go-plugin v0.7.0-dev
// 	protoc               v3.21.12
// source: internal/host/kubernetes/kubernetes.proto

package kubernetes

import (
	context "context"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type GVR struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Group    string `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
	Version  string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	Resource string `protobuf:"bytes,3,opt,name=resource,proto3" json:"resource,omitempty"`
}

func (x *GVR) ProtoReflect() protoreflect.Message {
	panic(`not implemented`)
}

func (x *GVR) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *GVR) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *GVR) GetResource() string {
	if x != nil {
		return x.Resource
	}
	return ""
}

type Namespaced struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name      string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Namespace string `protobuf:"bytes,2,opt,name=namespace,proto3" json:"namespace,omitempty"`
}

func (x *Namespaced) ProtoReflect() protoreflect.Message {
	panic(`not implemented`)
}

func (x *Namespaced) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Namespaced) GetNamespace() string {
	if x != nil {
		return x.Namespace
	}
	return ""
}

type GetRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Gvr  *GVR        `protobuf:"bytes,1,opt,name=gvr,proto3" json:"gvr,omitempty"`
	Name *Namespaced `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *GetRequest) ProtoReflect() protoreflect.Message {
	panic(`not implemented`)
}

func (x *GetRequest) GetGvr() *GVR {
	if x != nil {
		return x.Gvr
	}
	return nil
}

func (x *GetRequest) GetName() *Namespaced {
	if x != nil {
		return x.Name
	}
	return nil
}

type GetReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Payload []byte `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
}

func (x *GetReply) ProtoReflect() protoreflect.Message {
	panic(`not implemented`)
}

func (x *GetReply) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

// go:plugin type=host module=kubernetes
type Proxy interface {
	Get(context.Context, *GetRequest) (*GetReply, error)
}
