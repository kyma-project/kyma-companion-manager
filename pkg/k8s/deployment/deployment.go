package deployment

import (
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kyma-project/kyma-companion-manager/pkg/utils"
)

type Opt func(deployment *kappsv1.Deployment)

const secretVolumeSourceDefaultMode = 420

func NewDeployment(name, namespace string, opts ...Opt) *kappsv1.Deployment {
	newDeployment := &kappsv1.Deployment{
		TypeMeta: kmetav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: kappsv1.DeploymentSpec{
			Template: kcorev1.PodTemplateSpec{
				ObjectMeta: kmetav1.ObjectMeta{
					Name: name,
				},
				Spec: kcorev1.PodSpec{},
			},
		},
		Status: kappsv1.DeploymentStatus{},
	}
	// apply options.
	for _, o := range opts {
		o(newDeployment)
	}
	return newDeployment
}

func WithLabels(labels map[string]string) Opt {
	return func(d *kappsv1.Deployment) {
		d.ObjectMeta.Labels = labels
		d.Spec.Template.ObjectMeta.Labels = labels
	}
}

func WithSelectorLabels(labels map[string]string) Opt {
	return func(d *kappsv1.Deployment) {
		d.Spec.Selector = kmetav1.SetAsLabelSelector(labels)
	}
}

func WithPriorityClassName(name string) Opt {
	return func(deployment *kappsv1.Deployment) {
		deployment.Spec.Template.Spec.PriorityClassName = name
	}
}

func WithContainers(containers []kcorev1.Container) Opt {
	return func(d *kappsv1.Deployment) {
		d.Spec.Template.Spec.Containers = containers
	}
}

func WithRestartPolicyAlways() Opt {
	return func(deployment *kappsv1.Deployment) {
		deployment.Spec.Template.Spec.RestartPolicy = kcorev1.RestartPolicyAlways
	}
}

func WithTerminationGracePeriodSeconds(terminationGracePeriodSeconds int64) Opt {
	return func(deployment *kappsv1.Deployment) {
		deployment.Spec.Template.Spec.TerminationGracePeriodSeconds = &terminationGracePeriodSeconds
	}
}

func WithReplicas(replicas int32) Opt {
	return func(deployment *kappsv1.Deployment) {
		deployment.Spec.Replicas = utils.Int32Ptr(replicas)
	}
}

func WithSecurityContext(securityContext *kcorev1.PodSecurityContext) Opt {
	return func(deployment *kappsv1.Deployment) {
		deployment.Spec.Template.Spec.SecurityContext = securityContext
	}
}

func WithOwnerReferences(ownerReferences []kmetav1.OwnerReference) Opt {
	return func(deployment *kappsv1.Deployment) {
		deployment.OwnerReferences = ownerReferences
	}
}

func WithVolumeMountedSecret(secretName string) Opt {
	return func(deployment *kappsv1.Deployment) {
		volume := kcorev1.Volume{
			Name: secretName,
			VolumeSource: kcorev1.VolumeSource{
				Secret: &kcorev1.SecretVolumeSource{
					SecretName:  secretName,
					DefaultMode: utils.Int32Ptr(secretVolumeSourceDefaultMode),
				},
			},
		}

		if deployment.Spec.Template.Spec.Volumes == nil {
			deployment.Spec.Template.Spec.Volumes = []kcorev1.Volume{}
		}

		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes, volume)
	}
}
