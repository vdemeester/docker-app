// +build !ignore_autogenerated

/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was autogenerated by deepcopy-gen. Do not edit it manually!

package runtime

import (
	conversion "k8s.io/apimachinery/pkg/conversion"
	reflect "reflect"
)

// GetGeneratedDeepCopyFuncs returns the generated funcs, since we aren't registering them.
//
// Deprecated: deepcopy registration will go away when static deepcopy is fully implemented.
func GetGeneratedDeepCopyFuncs() []conversion.GeneratedDeepCopyFunc {
	return []conversion.GeneratedDeepCopyFunc{
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*RawExtension).DeepCopyInto(out.(*RawExtension))
			return nil
		}, InType: reflect.TypeOf(&RawExtension{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*Unknown).DeepCopyInto(out.(*Unknown))
			return nil
		}, InType: reflect.TypeOf(&Unknown{})},
		{Fn: func(in interface{}, out interface{}, c *conversion.Cloner) error {
			in.(*VersionedObjects).DeepCopyInto(out.(*VersionedObjects))
			return nil
		}, InType: reflect.TypeOf(&VersionedObjects{})},
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RawExtension) DeepCopyInto(out *RawExtension) {
	*out = *in
	if in.Raw != nil {
		in, out := &in.Raw, &out.Raw
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	if in.Object == nil {
		out.Object = nil
	} else {
		out.Object = in.Object.DeepCopyObject()
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RawExtension.
func (in *RawExtension) DeepCopy() *RawExtension {
	if in == nil {
		return nil
	}
	out := new(RawExtension)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Unknown) DeepCopyInto(out *Unknown) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Raw != nil {
		in, out := &in.Raw, &out.Raw
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Unknown.
func (in *Unknown) DeepCopy() *Unknown {
	if in == nil {
		return nil
	}
	out := new(Unknown)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new Object.
func (in *Unknown) DeepCopyObject() Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VersionedObjects) DeepCopyInto(out *VersionedObjects) {
	*out = *in
	if in.Objects != nil {
		in, out := &in.Objects, &out.Objects
		*out = make([]Object, len(*in))
		for i := range *in {
			if (*in)[i] == nil {
				(*out)[i] = nil
			} else {
				(*out)[i] = (*in)[i].DeepCopyObject()
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VersionedObjects.
func (in *VersionedObjects) DeepCopy() *VersionedObjects {
	if in == nil {
		return nil
	}
	out := new(VersionedObjects)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new Object.
func (in *VersionedObjects) DeepCopyObject() Object {
	if c := in.DeepCopy(); c != nil {
		return c
	} else {
		return nil
	}
}
