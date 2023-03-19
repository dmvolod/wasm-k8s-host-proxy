// Copyright 2021 OnMetal authors

package unstructuredutil

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
)

// UnstructuredSliceToObjectSlice transforms the given list of unstructured.Unstructured to a list of
// runtime.Object, copying the unstructured.Unstructured and using the pointers of them for the resulting runtime.Object.
func UnstructuredSliceToObjectSlice(unstructureds []unstructured.Unstructured) []runtime.Object {
	if unstructureds == nil {
		return nil
	}
	res := make([]runtime.Object, 0, len(unstructureds))
	for _, u := range unstructureds {
		u := u
		res = append(res, &u)
	}
	return res
}
