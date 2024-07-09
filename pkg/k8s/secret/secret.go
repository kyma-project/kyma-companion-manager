package secret

import (
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Opt func(secret *kcorev1.Secret)

func NewSecret(name, namespace string, opts ...Opt) *kcorev1.Secret {
	newDeployment := &kcorev1.Secret{
		TypeMeta: kmetav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Data: map[string][]byte{},
		Type: kcorev1.SecretTypeOpaque,
	}
	// apply options.
	for _, o := range opts {
		o(newDeployment)
	}
	return newDeployment
}

func WithLabels(labels map[string]string) Opt {
	return func(s *kcorev1.Secret) {
		s.ObjectMeta.Labels = labels
	}
}

func WithOwnerReferences(ownerReferences []kmetav1.OwnerReference) Opt {
	return func(s *kcorev1.Secret) {
		s.OwnerReferences = ownerReferences
	}
}

func WithDataKeyKey(key string, value []byte) Opt {
	return func(s *kcorev1.Secret) {
		s.Data[key] = value
	}
}
