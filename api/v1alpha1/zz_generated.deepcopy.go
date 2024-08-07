//go:build !ignore_autogenerated

/*
Copyright 2024.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AICoreConfig) DeepCopyInto(out *AICoreConfig) {
	*out = *in
	out.Secret = in.Secret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AICoreConfig.
func (in *AICoreConfig) DeepCopy() *AICoreConfig {
	if in == nil {
		return nil
	}
	out := new(AICoreConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Companion) DeepCopyInto(out *Companion) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Companion.
func (in *Companion) DeepCopy() *Companion {
	if in == nil {
		return nil
	}
	out := new(Companion)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Companion) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CompanionConfig) DeepCopyInto(out *CompanionConfig) {
	*out = *in
	out.Secret = in.Secret
	out.Replicas = in.Replicas
	in.Resources.DeepCopyInto(&out.Resources)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CompanionConfig.
func (in *CompanionConfig) DeepCopy() *CompanionConfig {
	if in == nil {
		return nil
	}
	out := new(CompanionConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CompanionList) DeepCopyInto(out *CompanionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Companion, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CompanionList.
func (in *CompanionList) DeepCopy() *CompanionList {
	if in == nil {
		return nil
	}
	out := new(CompanionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *CompanionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CompanionSpec) DeepCopyInto(out *CompanionSpec) {
	*out = *in
	out.AICore = in.AICore
	out.HanaCloud = in.HanaCloud
	out.Redis = in.Redis
	in.Companion.DeepCopyInto(&out.Companion)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CompanionSpec.
func (in *CompanionSpec) DeepCopy() *CompanionSpec {
	if in == nil {
		return nil
	}
	out := new(CompanionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CompanionStatus) DeepCopyInto(out *CompanionStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CompanionStatus.
func (in *CompanionStatus) DeepCopy() *CompanionStatus {
	if in == nil {
		return nil
	}
	out := new(CompanionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HanaConfig) DeepCopyInto(out *HanaConfig) {
	*out = *in
	out.Secret = in.Secret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HanaConfig.
func (in *HanaConfig) DeepCopy() *HanaConfig {
	if in == nil {
		return nil
	}
	out := new(HanaConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RedisConfig) DeepCopyInto(out *RedisConfig) {
	*out = *in
	out.Secret = in.Secret
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RedisConfig.
func (in *RedisConfig) DeepCopy() *RedisConfig {
	if in == nil {
		return nil
	}
	out := new(RedisConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicasConfig) DeepCopyInto(out *ReplicasConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicasConfig.
func (in *ReplicasConfig) DeepCopy() *ReplicasConfig {
	if in == nil {
		return nil
	}
	out := new(ReplicasConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretSpec) DeepCopyInto(out *SecretSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretSpec.
func (in *SecretSpec) DeepCopy() *SecretSpec {
	if in == nil {
		return nil
	}
	out := new(SecretSpec)
	in.DeepCopyInto(out)
	return out
}
