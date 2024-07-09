package k8s

import (
	"context"

	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate go run github.com/vektra/mockery/v2 --name=Client --outpkg=mocks --case=underscore
type Client interface {
	GetDeployment(ctx context.Context, name, namespace string) (*kappsv1.Deployment, error)
	UpdateDeployment(ctx context.Context, deployment *kappsv1.Deployment) error
	DeleteDeployment(ctx context.Context, name, namespace string) error
	GetSecret(ctx context.Context, name, namespace string) (*kcorev1.Secret, error)
	GetConfigMap(ctx context.Context, name, namespace string) (*kcorev1.ConfigMap, error)
	DeleteResource(ctx context.Context, object client.Object) error
	PatchApply(ctx context.Context, object client.Object) error
}

type KubeClient struct {
	fieldManager  string
	client        client.Client
	dynamicClient dynamic.Interface
}

func NewKubeClient(client client.Client, fieldManager string,
	dynamicClient dynamic.Interface,
) Client {
	return &KubeClient{
		client:        client,
		fieldManager:  fieldManager,
		dynamicClient: dynamicClient,
	}
}

func (c *KubeClient) GetDeployment(ctx context.Context, name, namespace string) (*kappsv1.Deployment, error) {
	deployment := &kappsv1.Deployment{}
	if err := c.client.Get(ctx, client.ObjectKey{Name: name, Namespace: namespace}, deployment); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return deployment, nil
}

func (c *KubeClient) UpdateDeployment(ctx context.Context, deployment *kappsv1.Deployment) error {
	return c.client.Update(ctx, deployment)
}

func (c *KubeClient) DeleteDeployment(ctx context.Context, name, namespace string) error {
	deployment := &kappsv1.Deployment{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	if err := c.client.Delete(ctx, deployment); err != nil {
		return client.IgnoreNotFound(err)
	}
	return nil
}

func (c *KubeClient) DeleteResource(ctx context.Context, object client.Object) error {
	if err := c.client.Delete(ctx, object); err != nil {
		return client.IgnoreNotFound(err)
	}
	return nil
}

// PatchApply uses the server-side apply to create/update the resource.
// The object must define `GVK` (i.e. object.TypeMeta).
func (c *KubeClient) PatchApply(ctx context.Context, object client.Object) error {
	return c.client.Patch(ctx, object, client.Apply, &client.PatchOptions{
		Force:        ptr.To(true),
		FieldManager: c.fieldManager,
	})
}

// GetSecret returns the secret with the given name.
func (c *KubeClient) GetSecret(ctx context.Context, name, namespace string) (*kcorev1.Secret, error) {
	secret := &kcorev1.Secret{}
	err := c.client.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}

// GetConfigMap returns a ConfigMap based on the given name and namespace.
func (c *KubeClient) GetConfigMap(ctx context.Context, name, namespace string) (*kcorev1.ConfigMap, error) {
	cm := &kcorev1.ConfigMap{}
	key := client.ObjectKey{Name: name, Namespace: namespace}
	if err := c.client.Get(ctx, key, cm); err != nil {
		return nil, err
	}
	return cm, nil
}
