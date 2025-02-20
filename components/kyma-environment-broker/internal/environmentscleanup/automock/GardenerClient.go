// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	"context"

	mock "github.com/stretchr/testify/mock"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// GardenerClient is an autogenerated mock type for the GardenerClient type
type GardenerClient struct {
	mock.Mock
}

// List provides a mock function with given fields: opts
func (_m *GardenerClient) List(ctx context.Context, opts v1.ListOptions) (*unstructured.UnstructuredList, error) {
	ret := _m.Called(opts)

	var r0 *unstructured.UnstructuredList
	if rf, ok := ret.Get(0).(func(v1.ListOptions) *unstructured.UnstructuredList); ok {
		r0 = rf(opts)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*unstructured.UnstructuredList)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(v1.ListOptions) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
