//nolint:dupl // these comparison functions all look very similar
package equality

import (
	"testing"

	"github.com/stretchr/testify/require"
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"

	"github.com/kyma-project/kyma-companion-manager/pkg/utils"
	testutils "github.com/kyma-project/kyma-companion-manager/test/utils"
)

func Test_deploymentEqual(t *testing.T) {
	defaultDeployment := testutils.NewCompanionDeployment("test-companion", "test-namespace")
	ownerReference := func(version, kind, name, uid string, controller, block *bool) kmetav1.OwnerReference {
		return kmetav1.OwnerReference{
			APIVersion:         version,
			Kind:               kind,
			Name:               name,
			UID:                types.UID(uid),
			Controller:         controller,
			BlockOwnerDeletion: block,
		}
	}

	testCases := map[string]struct {
		getDeployment1 func() *kappsv1.Deployment
		getDeployment2 func() *kappsv1.Deployment
		expectedResult bool
	}{
		"should be equal if same default deployments": {
			getDeployment1: func() *kappsv1.Deployment {
				p := defaultDeployment.DeepCopy()
				return p
			},
			getDeployment2: func() *kappsv1.Deployment {
				p := defaultDeployment.DeepCopy()
				return p
			},
			expectedResult: true,
		},
		"should be unequal if container image changes": {
			getDeployment1: func() *kappsv1.Deployment {
				p := defaultDeployment.DeepCopy()
				p.Spec.Template.Spec.Containers[0].Image = "new-publisher-img"
				return p
			},
			getDeployment2: func() *kappsv1.Deployment {
				return defaultDeployment.DeepCopy()
			},
			expectedResult: false,
		},
		"should be unequal if env var changes": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Containers[0].Env = []kcorev1.EnvVar{
					{
						Name:  "key",
						Value: "value1",
					},
				}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Containers[0].Env = []kcorev1.EnvVar{
					{
						Name:  "key",
						Value: "value2",
					},
				}
				return deploy
			},
			expectedResult: false,
		},
		"should be equal if env var are same": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Containers[0].Env = []kcorev1.EnvVar{
					{
						Name:  "key",
						Value: "value1",
					},
				}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Containers[0].Env = []kcorev1.EnvVar{
					{
						Name:  "key",
						Value: "value1",
					},
				}
				return deploy
			},
			expectedResult: true,
		},
		"should not be equal if replicas changes": {
			getDeployment1: func() *kappsv1.Deployment {
				replicas := int32(1)
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Replicas = &replicas
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				replicas := int32(2)
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Replicas = &replicas
				return deploy
			},
			expectedResult: false,
		},
		"should be equal if replicas are the same": {
			getDeployment1: func() *kappsv1.Deployment {
				replicas := int32(2)
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Replicas = &replicas
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				replicas := int32(2)
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Replicas = &replicas
				return deploy
			},
			expectedResult: true,
		},
		"should be equal if spec annotations are nil and empty": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Annotations = nil
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Annotations = map[string]string{}
				return deploy
			},
			expectedResult: true,
		},
		"should be unequal if spec annotations changes": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Annotations = map[string]string{"key": "value1"}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Annotations = map[string]string{"key": "value2"}
				return deploy
			},
			expectedResult: false,
		},
		"should be equal if spec Labels are nil and empty": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Labels = nil
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Labels = map[string]string{}
				return deploy
			},
			expectedResult: true,
		},
		"should be unequal if spec Labels changes": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Labels = map[string]string{"key": "value1"}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Labels = map[string]string{"key": "value2"}
				return deploy
			},
			expectedResult: false,
		},
		"should be unequal if Labels changes": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Labels = map[string]string{"key": "value1"}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Labels = map[string]string{"key": "value2"}
				return deploy
			},
			expectedResult: false,
		},
		"should be unequal if owner reference changes": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.OwnerReferences = []kmetav1.OwnerReference{
					ownerReference("v-0", "k", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.OwnerReferences = []kmetav1.OwnerReference{
					ownerReference("v-2", "k", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				}
				return deploy
			},
			expectedResult: false,
		},
		"should be equal if owner references are same": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.OwnerReferences = []kmetav1.OwnerReference{
					ownerReference("v-0", "k", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.OwnerReferences = []kmetav1.OwnerReference{
					ownerReference("v-0", "k", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				}
				return deploy
			},
			expectedResult: true,
		},
		"should be unequal if volumes are different": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Volumes = []kcorev1.Volume{
					{
						Name: "test1",
						VolumeSource: kcorev1.VolumeSource{
							Secret: &kcorev1.SecretVolumeSource{
								SecretName:  "test1",
								DefaultMode: utils.Int32Ptr(420),
							},
						},
					},
				}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Volumes = []kcorev1.Volume{
					{
						Name: "test2",
						VolumeSource: kcorev1.VolumeSource{
							Secret: &kcorev1.SecretVolumeSource{
								SecretName:  "test1",
								DefaultMode: utils.Int32Ptr(420),
							},
						},
					},
				}
				return deploy
			},
			expectedResult: false,
		},
		"should be equal if volumes are same": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Volumes = []kcorev1.Volume{
					{
						Name: "test1",
						VolumeSource: kcorev1.VolumeSource{
							Secret: &kcorev1.SecretVolumeSource{
								SecretName:  "test1",
								DefaultMode: utils.Int32Ptr(420),
							},
						},
					},
				}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Volumes = []kcorev1.Volume{
					{
						Name: "test1",
						VolumeSource: kcorev1.VolumeSource{
							Secret: &kcorev1.SecretVolumeSource{
								SecretName:  "test1",
								DefaultMode: utils.Int32Ptr(420),
							},
						},
					},
				}
				return deploy
			},
			expectedResult: true,
		},
		"should be unequal if volumeMounts are different": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Containers[0].VolumeMounts = []kcorev1.VolumeMount{
					{
						Name:      "test1",
						ReadOnly:  true,
						MountPath: "/mnt/secrets",
					},
				}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Containers[0].VolumeMounts = []kcorev1.VolumeMount{
					{
						Name:      "test1",
						ReadOnly:  true,
						MountPath: "/mnt/changed",
					},
				}
				return deploy
			},
			expectedResult: false,
		},
		"should be equal if volumeMounts are same": {
			getDeployment1: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Containers[0].VolumeMounts = []kcorev1.VolumeMount{
					{
						Name:      "test1",
						ReadOnly:  true,
						MountPath: "/mnt/secrets",
					},
				}
				return deploy
			},
			getDeployment2: func() *kappsv1.Deployment {
				deploy := defaultDeployment.DeepCopy()
				deploy.Spec.Template.Spec.Containers[0].VolumeMounts = []kcorev1.VolumeMount{
					{
						Name:      "test1",
						ReadOnly:  true,
						MountPath: "/mnt/secrets",
					},
				}
				return deploy
			},
			expectedResult: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			if deploymentEqual(tc.getDeployment1(), tc.getDeployment2()) != tc.expectedResult {
				t.Errorf("expected output to be %t", tc.expectedResult)
			}
		})
	}
}

func Test_ownerReferencesDeepEqual(t *testing.T) {
	ownerReference := func(version, kind, name, uid string, controller, block *bool) kmetav1.OwnerReference {
		return kmetav1.OwnerReference{
			APIVersion:         version,
			Kind:               kind,
			Name:               name,
			UID:                types.UID(uid),
			Controller:         controller,
			BlockOwnerDeletion: block,
		}
	}

	testCases := []struct {
		name                  string
		givenOwnerReferences1 []kmetav1.OwnerReference
		givenOwnerReferences2 []kmetav1.OwnerReference
		wantEqual             bool
	}{
		{
			name:                  "both OwnerReferences are nil",
			givenOwnerReferences1: nil,
			givenOwnerReferences2: nil,
			wantEqual:             true,
		},
		{
			name:                  "both OwnerReferences are empty",
			givenOwnerReferences1: []kmetav1.OwnerReference{},
			givenOwnerReferences2: []kmetav1.OwnerReference{},
			wantEqual:             true,
		},
		{
			name: "same OwnerReferences and same order",
			givenOwnerReferences1: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				ownerReference("v-1", "k-1", "n-1", "u-1", ptr.To(false), ptr.To(false)),
				ownerReference("v-2", "k-2", "n-2", "u-2", ptr.To(false), ptr.To(false)),
			},
			givenOwnerReferences2: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				ownerReference("v-1", "k-1", "n-1", "u-1", ptr.To(false), ptr.To(false)),
				ownerReference("v-2", "k-2", "n-2", "u-2", ptr.To(false), ptr.To(false)),
			},
			wantEqual: true,
		},
		{
			name: "same OwnerReferences but different order",
			givenOwnerReferences1: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				ownerReference("v-1", "k-1", "n-1", "u-1", ptr.To(false), ptr.To(false)),
				ownerReference("v-2", "k-2", "n-2", "u-2", ptr.To(false), ptr.To(false)),
			},
			givenOwnerReferences2: []kmetav1.OwnerReference{
				ownerReference("v-2", "k-2", "n-2", "u-2", ptr.To(false), ptr.To(false)),
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				ownerReference("v-1", "k-1", "n-1", "u-1", ptr.To(false), ptr.To(false)),
			},
			wantEqual: true,
		},
		{
			name: "different OwnerReference APIVersion",
			givenOwnerReferences1: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
			},
			givenOwnerReferences2: []kmetav1.OwnerReference{
				ownerReference("v-1", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
			},
			wantEqual: false,
		},
		{
			name: "different OwnerReference Kind",
			givenOwnerReferences1: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
			},
			givenOwnerReferences2: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-1", "n-0", "u-0", ptr.To(false), ptr.To(false)),
			},
			wantEqual: false,
		},
		{
			name: "different OwnerReference Name",
			givenOwnerReferences1: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
			},
			givenOwnerReferences2: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-1", "u-0", ptr.To(false), ptr.To(false)),
			},
			wantEqual: false,
		},
		{
			name: "different OwnerReference UID",
			givenOwnerReferences1: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
			},
			givenOwnerReferences2: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-1", ptr.To(false), ptr.To(false)),
			},
			wantEqual: false,
		},
		{
			name: "different OwnerReference Controller",
			givenOwnerReferences1: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
			},
			givenOwnerReferences2: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(true), ptr.To(false)),
			},
			wantEqual: false,
		},
		{
			name: "different OwnerReference BlockOwnerDeletion",
			givenOwnerReferences1: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(false)),
			},
			givenOwnerReferences2: []kmetav1.OwnerReference{
				ownerReference("v-0", "k-0", "n-0", "u-0", ptr.To(false), ptr.To(true)),
			},
			wantEqual: false,
		},
	}

	for _, tc := range testCases {
		testcase := tc
		t.Run(testcase.name, func(t *testing.T) {
			require.Equal(t, testcase.wantEqual, ownerReferencesDeepEqual(testcase.givenOwnerReferences1,
				testcase.givenOwnerReferences2))
		})
	}
}

func Test_containerEqual(t *testing.T) {
	quantityA, _ := resource.ParseQuantity("5m")
	quantityB, _ := resource.ParseQuantity("10k")

	type args struct {
		c1 *kcorev1.Container
		c2 *kcorev1.Container
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "container are equal",
			args: args{
				c1: &kcorev1.Container{
					Name:       "test",
					Image:      "bla",
					Command:    []string{"1", "2"},
					Args:       []string{"a", "b"},
					WorkingDir: "foodir",
					Ports: []kcorev1.ContainerPort{{
						Name:          "testport",
						HostPort:      1,
						ContainerPort: 2,
						Protocol:      "http",
						HostIP:        "192.168.1.1",
					}},
					Resources: kcorev1.ResourceRequirements{
						Limits: map[kcorev1.ResourceName]resource.Quantity{
							"cpu": quantityA,
						},
						Requests: map[kcorev1.ResourceName]resource.Quantity{
							"mem": quantityA,
						},
					},
					ReadinessProbe: &kcorev1.Probe{
						ProbeHandler:        kcorev1.ProbeHandler{},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      0,
						PeriodSeconds:       0,
						SuccessThreshold:    0,
						FailureThreshold:    0,
					},
				},
				c2: &kcorev1.Container{
					Name:       "test",
					Image:      "bla",
					Command:    []string{"1", "2"},
					Args:       []string{"a", "b"},
					WorkingDir: "foodir",
					Ports: []kcorev1.ContainerPort{{
						Name:          "testport",
						HostPort:      1,
						ContainerPort: 2,
						Protocol:      "http",
						HostIP:        "192.168.1.1",
					}},
					Resources: kcorev1.ResourceRequirements{
						Limits: map[kcorev1.ResourceName]resource.Quantity{
							"cpu": quantityA,
						},
						Requests: map[kcorev1.ResourceName]resource.Quantity{
							"mem": quantityA,
						},
					},
					ReadinessProbe: &kcorev1.Probe{
						ProbeHandler:        kcorev1.ProbeHandler{},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      0,
						PeriodSeconds:       0,
						SuccessThreshold:    0,
						FailureThreshold:    0,
					},
				},
			},
			want: true,
		},
		{
			name: "ContainerPort are not equal",
			args: args{
				c1: &kcorev1.Container{
					Name:       "test",
					Image:      "bla",
					Command:    []string{"1", "2"},
					Args:       []string{"a", "b"},
					WorkingDir: "foodir",
					Ports: []kcorev1.ContainerPort{{
						Name:          "testport",
						HostPort:      1,
						ContainerPort: 2,
						Protocol:      "http",
						HostIP:        "192.168.1.1",
					}},
					Resources: kcorev1.ResourceRequirements{
						Limits: map[kcorev1.ResourceName]resource.Quantity{
							"cpu": quantityA,
						},
						Requests: map[kcorev1.ResourceName]resource.Quantity{
							"mem": quantityA,
						},
					},
					ReadinessProbe: &kcorev1.Probe{
						ProbeHandler:        kcorev1.ProbeHandler{},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      0,
						PeriodSeconds:       0,
						SuccessThreshold:    0,
						FailureThreshold:    0,
					},
				},
				c2: &kcorev1.Container{
					Name:       "test",
					Image:      "bla",
					Command:    []string{"1", "2"},
					Args:       []string{"a", "b"},
					WorkingDir: "foodir",
					Ports: []kcorev1.ContainerPort{{
						Name:          "testport",
						HostPort:      1,
						ContainerPort: 3,
						Protocol:      "http",
						HostIP:        "192.168.1.1",
					}},
					Resources: kcorev1.ResourceRequirements{
						Limits: map[kcorev1.ResourceName]resource.Quantity{
							"cpu": quantityA,
						},
						Requests: map[kcorev1.ResourceName]resource.Quantity{
							"mem": quantityA,
						},
					},
					ReadinessProbe: &kcorev1.Probe{
						ProbeHandler:        kcorev1.ProbeHandler{},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      0,
						PeriodSeconds:       0,
						SuccessThreshold:    0,
						FailureThreshold:    0,
					},
				},
			},
			want: false,
		},
		{
			name: "resources are not equal",
			args: args{
				c1: &kcorev1.Container{
					Name:       "test",
					Image:      "bla",
					Command:    []string{"1", "2"},
					Args:       []string{"a", "b"},
					WorkingDir: "foodir",
					Ports: []kcorev1.ContainerPort{{
						Name:          "testport",
						HostPort:      1,
						ContainerPort: 2,
						Protocol:      "http",
						HostIP:        "192.168.1.1",
					}},
					Resources: kcorev1.ResourceRequirements{
						Limits: map[kcorev1.ResourceName]resource.Quantity{
							"cpu": quantityA,
						},
						Requests: map[kcorev1.ResourceName]resource.Quantity{
							"mem": quantityA,
						},
					},
					ReadinessProbe: &kcorev1.Probe{
						ProbeHandler:        kcorev1.ProbeHandler{},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      0,
						PeriodSeconds:       0,
						SuccessThreshold:    0,
						FailureThreshold:    0,
					},
				},
				c2: &kcorev1.Container{
					Name:       "test",
					Image:      "bla",
					Command:    []string{"1", "2"},
					Args:       []string{"a", "b"},
					WorkingDir: "foodir",
					Ports: []kcorev1.ContainerPort{{
						Name:          "testport",
						HostPort:      1,
						ContainerPort: 2,
						Protocol:      "http",
						HostIP:        "192.168.1.1",
					}},
					Resources: kcorev1.ResourceRequirements{
						Limits: map[kcorev1.ResourceName]resource.Quantity{
							"cpu": quantityB,
						},
						Requests: map[kcorev1.ResourceName]resource.Quantity{
							"mem": quantityB,
						},
					},
					ReadinessProbe: &kcorev1.Probe{
						ProbeHandler:        kcorev1.ProbeHandler{},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      0,
						PeriodSeconds:       0,
						SuccessThreshold:    0,
						FailureThreshold:    0,
					},
				},
			},
			want: false,
		},
		{
			name: "ports are not equal",
			args: args{
				c1: &kcorev1.Container{
					Name:       "test",
					Image:      "bla",
					Command:    []string{"1", "2"},
					Args:       []string{"a", "b"},
					WorkingDir: "foodir",
					Ports: []kcorev1.ContainerPort{
						{
							Name:          "testport-0",
							HostPort:      1,
							ContainerPort: 2,
							Protocol:      "http",
							HostIP:        "192.168.1.1",
						},
					},
					Resources: kcorev1.ResourceRequirements{
						Limits: map[kcorev1.ResourceName]resource.Quantity{
							"cpu": quantityA,
						},
						Requests: map[kcorev1.ResourceName]resource.Quantity{
							"mem": quantityA,
						},
					},
					ReadinessProbe: &kcorev1.Probe{
						ProbeHandler:        kcorev1.ProbeHandler{},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      0,
						PeriodSeconds:       0,
						SuccessThreshold:    0,
						FailureThreshold:    0,
					},
				},
				c2: &kcorev1.Container{
					Name:       "test",
					Image:      "bla",
					Command:    []string{"1", "2"},
					Args:       []string{"a", "b"},
					WorkingDir: "foodir",
					Ports: []kcorev1.ContainerPort{
						{
							Name:          "testport-0",
							HostPort:      1,
							ContainerPort: 2,
							Protocol:      "http",
							HostIP:        "192.168.1.1",
						},
						{
							Name:          "testport-1",
							HostPort:      1,
							ContainerPort: 2,
							Protocol:      "http",
							HostIP:        "192.168.1.1",
						},
					},
					Resources: kcorev1.ResourceRequirements{
						Limits: map[kcorev1.ResourceName]resource.Quantity{
							"cpu": quantityA,
						},
						Requests: map[kcorev1.ResourceName]resource.Quantity{
							"mem": quantityA,
						},
					},
					ReadinessProbe: &kcorev1.Probe{
						ProbeHandler:        kcorev1.ProbeHandler{},
						InitialDelaySeconds: 0,
						TimeoutSeconds:      0,
						PeriodSeconds:       0,
						SuccessThreshold:    0,
						FailureThreshold:    0,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := containerEqual(tt.args.c1, tt.args.c2); got != tt.want {
				t.Errorf("containerEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_serviceEqual(t *testing.T) {
	const (
		name0 = "name-0"
		name1 = "name-1"

		namespace0 = "namespace0"
		namespace1 = "namespace1"
	)

	var (
		ownerReferences0 = []kmetav1.OwnerReference{
			{
				Name:               "name-0",
				APIVersion:         "version-0",
				Kind:               "kind-0",
				UID:                "000000",
				Controller:         ptr.To(false),
				BlockOwnerDeletion: ptr.To(false),
			},
		}
		ownerReferences1 = []kmetav1.OwnerReference{
			{
				Name:               "name-1",
				APIVersion:         "version-0",
				Kind:               "kind-0",
				UID:                "000000",
				Controller:         ptr.To(false),
				BlockOwnerDeletion: ptr.To(false),
			},
		}

		ports0 = []kcorev1.ServicePort{
			{
				Name:        "name-0",
				Protocol:    "protocol-0",
				AppProtocol: nil,
				Port:        0,
				TargetPort: intstr.IntOrString{
					Type:   0,
					IntVal: 0,
					StrVal: "val-0",
				},
				NodePort: 0,
			},
		}
		ports1 = []kcorev1.ServicePort{
			{
				Name:        "name-1",
				Protocol:    "protocol-0",
				AppProtocol: nil,
				Port:        0,
				TargetPort: intstr.IntOrString{
					Type:   0,
					IntVal: 0,
					StrVal: "val-0",
				},
				NodePort: 0,
			},
		}

		selector0 = map[string]string{
			"key": "val-0",
		}
		selector1 = map[string]string{
			"key": "val-1",
		}
	)

	type args struct {
		a *kcorev1.Service
		b *kcorev1.Service
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Services are equal",
			args: args{
				a: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
				b: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
			},
			want: true,
		},
		{
			name: "Service names are not equal",
			args: args{
				a: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
				b: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name1,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
			},
			want: false,
		},
		{
			name: "Service namespaces are not equal",
			args: args{
				a: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
				b: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace1,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
			},
			want: false,
		},
		{
			name: "Service OwnerReferences are not equal",
			args: args{
				a: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
				b: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences1,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
			},
			want: false,
		},
		{
			name: "Service ports are not equal",
			args: args{
				a: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
				b: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports1,
						Selector: selector0,
					},
				},
			},
			want: false,
		},
		{
			name: "Service selectors are not equal",
			args: args{
				a: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector0,
					},
				},
				b: &kcorev1.Service{
					ObjectMeta: kmetav1.ObjectMeta{
						Name:            name0,
						Namespace:       namespace0,
						OwnerReferences: ownerReferences0,
					},
					Spec: kcorev1.ServiceSpec{
						Ports:    ports0,
						Selector: selector1,
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := serviceEqual(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("serviceEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_envEqual(t *testing.T) {
	type args struct {
		e1 []kcorev1.EnvVar
		e2 []kcorev1.EnvVar
	}

	var11 := kcorev1.EnvVar{
		Name:  "var1",
		Value: "var1",
	}
	var12 := kcorev1.EnvVar{
		Name:  "var1",
		Value: "var2",
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "envs equal, order equals",
			args: args{
				e1: []kcorev1.EnvVar{var11, var12},
				e2: []kcorev1.EnvVar{var11, var12},
			},
			want: true,
		},
		{
			name: "envs equal, different order",
			args: args{
				e1: []kcorev1.EnvVar{var11, var12},
				e2: []kcorev1.EnvVar{var12, var11},
			},
			want: true,
		},
		{
			name: "different length",
			args: args{
				e1: []kcorev1.EnvVar{var11, var11},
				e2: []kcorev1.EnvVar{var11},
			},
			want: false,
		},
		{
			name: "envs different",
			args: args{
				e1: []kcorev1.EnvVar{var11, var12},
				e2: []kcorev1.EnvVar{var11, var11},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := envEqual(tt.args.e1, tt.args.e2); got != tt.want {
				t.Errorf("envEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_probeEqual(t *testing.T) {
	probe := &kcorev1.Probe{}

	type args struct {
		p1 *kcorev1.Probe
		p2 *kcorev1.Probe
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Probe refs are equal",
			args: args{
				p1: probe,
				p2: probe,
			},
			want: true,
		},
		{
			name: "one Probe is Nil",
			args: args{
				p1: nil,
				p2: probe,
			},
			want: false,
		},
		{
			name: "both Probes are Nil",
			args: args{
				p1: nil,
				p2: nil,
			},
			want: true,
		},
		{
			name: "Probes are not equal",
			args: args{
				p1: &kcorev1.Probe{
					InitialDelaySeconds: 1,
				},
				p2: &kcorev1.Probe{
					InitialDelaySeconds: 2,
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := probeEqual(tt.args.p1, tt.args.p2); got != tt.want {
				t.Errorf("probeEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_secretEqual(t *testing.T) {
	defaultSecret := &kcorev1.Secret{
		TypeMeta: kmetav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      "test",
			Namespace: "testNamespace",
		},
		Data: map[string][]byte{
			"hana-db-secret": []byte("hana-db-secret"),
			"redis-secret":   []byte("redis-secret"),
			"ai-core-config": []byte("ai-core-config"),
			"ai-core-secret": []byte("ai-core-secret"),
		},
		Type: kcorev1.SecretTypeOpaque,
	}

	ownerReference := func(version, kind, name, uid string, controller, block *bool) kmetav1.OwnerReference {
		return kmetav1.OwnerReference{
			APIVersion:         version,
			Kind:               kind,
			Name:               name,
			UID:                types.UID(uid),
			Controller:         controller,
			BlockOwnerDeletion: block,
		}
	}

	tests := []struct {
		name           string
		getSecret1     func() *kcorev1.Secret
		getSecret2     func() *kcorev1.Secret
		expectedResult bool
	}{
		{
			name: "should be equal when secrets are equal",
			getSecret1: func() *kcorev1.Secret {
				return defaultSecret.DeepCopy()
			},
			getSecret2: func() *kcorev1.Secret {
				return defaultSecret.DeepCopy()
			},
			expectedResult: true,
		},
		{
			name: "should be unequal when owner references are different",
			getSecret1: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.OwnerReferences = []kmetav1.OwnerReference{
					ownerReference("v-0", "k", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				}
				return secret
			},
			getSecret2: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.OwnerReferences = []kmetav1.OwnerReference{
					ownerReference("v-1", "k", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				}
				return secret
			},
			expectedResult: false,
		},
		{
			name: "should be equal when owner references are same",
			getSecret1: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.OwnerReferences = []kmetav1.OwnerReference{
					ownerReference("v-0", "k", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				}
				return secret
			},
			getSecret2: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.OwnerReferences = []kmetav1.OwnerReference{
					ownerReference("v-0", "k", "n-0", "u-0", ptr.To(false), ptr.To(false)),
				}
				return secret
			},
			expectedResult: true,
		},
		{
			name: "should be unequal when data is different",
			getSecret1: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.Data = map[string][]byte{
					"hana-db-secret": []byte("hana-db-secret"),
					"redis-secret":   []byte("redis-secret"),
					"ai-core-config": []byte("ai-core-config"),
					"ai-core-secret": []byte("ai-core-secret"),
				}
				return secret
			},
			getSecret2: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.Data = map[string][]byte{
					"hana-db-secret": []byte("changed"),
					"redis-secret":   []byte("redis-secret"),
					"ai-core-config": []byte("ai-core-config"),
					"ai-core-secret": []byte("ai-core-secret"),
				}
				return secret
			},
			expectedResult: false,
		},
		{
			name: "should be unequal when name is different",
			getSecret1: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.Name = "test1"
				return secret
			},
			getSecret2: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.Name = "test2"
				return secret
			},
			expectedResult: false,
		},
		{
			name: "should be unequal when labels are different",
			getSecret1: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.Labels = map[string]string{
					"key": "val",
				}
				return secret
			},
			getSecret2: func() *kcorev1.Secret {
				secret := defaultSecret.DeepCopy()
				secret.Labels = map[string]string{
					"key": "val-changed",
				}
				return secret
			},
			expectedResult: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if secretEqual(tc.getSecret1(), tc.getSecret2()) != tc.expectedResult {
				t.Errorf("expected output to be %t", tc.expectedResult)
			}
		})
	}
}
