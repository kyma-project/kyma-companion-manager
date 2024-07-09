package backendmanager

import (
	kcorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	kcmutils "github.com/kyma-project/kyma-companion-manager/pkg/utils"
)

const BackendResourceName = "kyma-companion-backend"

func getContainerPorts() []kcorev1.ContainerPort {
	return []kcorev1.ContainerPort{
		{
			Name:          backendPortName,
			ContainerPort: backendPortNum,
		},
		{
			Name:          backendMetricsPortName,
			ContainerPort: backendMetricsPortNum,
		},
	}
}

func getPodSecurityContext() *kcorev1.PodSecurityContext { //nolint: unused // It will be used in the future.
	const id = 10001
	return &kcorev1.PodSecurityContext{
		FSGroup:      kcmutils.Int64Ptr(id),
		RunAsUser:    kcmutils.Int64Ptr(id),
		RunAsGroup:   kcmutils.Int64Ptr(id),
		RunAsNonRoot: kcmutils.BoolPtr(true),
		SeccompProfile: &kcorev1.SeccompProfile{
			Type: kcorev1.SeccompProfileTypeRuntimeDefault,
		},
	}
}

func getContainerSecurityContext() *kcorev1.SecurityContext { //nolint: unused // It will be used in the future.
	return &kcorev1.SecurityContext{
		Privileged:               kcmutils.BoolPtr(false),
		AllowPrivilegeEscalation: kcmutils.BoolPtr(false),
		RunAsNonRoot:             kcmutils.BoolPtr(true),
		Capabilities: &kcorev1.Capabilities{
			Drop: []kcorev1.Capability{"ALL"},
		},
	}
}

func getReadinessProbe() *kcorev1.Probe {
	const (
		readyPath  = "/readyz"
		readyPort  = 8000
		maxFailure = 3
	)
	return &kcorev1.Probe{
		ProbeHandler: kcorev1.ProbeHandler{
			HTTPGet: &kcorev1.HTTPGetAction{
				Path:   readyPath,
				Port:   intstr.FromInt32(readyPort),
				Scheme: kcorev1.URISchemeHTTP,
			},
		},
		FailureThreshold: maxFailure,
	}
}

func getLivenessProbe() *kcorev1.Probe {
	const (
		healthPath = "/healthz"
		healthPort = 8000
		minSuccess = 1
		maxError   = 3
	)
	return &kcorev1.Probe{
		ProbeHandler: kcorev1.ProbeHandler{
			HTTPGet: &kcorev1.HTTPGetAction{
				Path:   healthPath,
				Port:   intstr.FromInt32(healthPort),
				Scheme: kcorev1.URISchemeHTTP,
			},
		},
		InitialDelaySeconds: livenessInitialDelaySecs,
		TimeoutSeconds:      livenessTimeoutSecs,
		PeriodSeconds:       livenessPeriodSecs,
		SuccessThreshold:    minSuccess,
		FailureThreshold:    maxError,
	}
}

func getResources(requestsCPU, requestsMemory, limitsCPU, limitsMemory string) kcorev1.ResourceRequirements {
	return kcorev1.ResourceRequirements{
		Requests: kcorev1.ResourceList{
			kcorev1.ResourceCPU:    resource.MustParse(requestsCPU),
			kcorev1.ResourceMemory: resource.MustParse(requestsMemory),
		},
		Limits: kcorev1.ResourceList{
			kcorev1.ResourceCPU:    resource.MustParse(limitsCPU),
			kcorev1.ResourceMemory: resource.MustParse(limitsMemory),
		},
	}
}

func getOwnerReferences(companion kcmv1alpha1.Companion) []kmetav1.OwnerReference {
	return []kmetav1.OwnerReference{
		{
			APIVersion:         companion.APIVersion,
			Kind:               companion.Kind,
			Name:               companion.Name,
			UID:                companion.UID,
			Controller:         kcmutils.BoolPtr(true),
			BlockOwnerDeletion: kcmutils.BoolPtr(true),
		},
	}
}
