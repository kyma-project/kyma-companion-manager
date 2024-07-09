package equality

import (
	"reflect"

	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
)

// Semantic can do semantic deep equality checks for API objects. Fields which
// are not relevant for the reconciliation logic are intentionally omitted.
//
//nolint:gochecknoglobals // same pattern as in apimachinery
var Semantic = conversion.EqualitiesOrDie(
	deploymentEqual,
	serviceEqual,
	secretEqual,
)

func serviceEqual(a, b *kcorev1.Service) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if !ownerReferencesDeepEqual(a.OwnerReferences, b.OwnerReferences) {
		return false
	}
	if a.Name != b.Name || a.Namespace != b.Namespace {
		return false
	}
	if !reflect.DeepEqual(a.Spec.Selector, b.Spec.Selector) {
		return false
	}
	return reflect.DeepEqual(a.Spec.Ports, b.Spec.Ports)
}

func secretEqual(a, b *kcorev1.Secret) bool {
	if a == b {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if !reflect.DeepEqual(a.Labels, b.Labels) {
		return false
	}

	if !reflect.DeepEqual(a.Type, b.Type) {
		return false
	}

	if !ownerReferencesDeepEqual(a.OwnerReferences, b.OwnerReferences) {
		return false
	}

	if a.Name != b.Name || a.Namespace != b.Namespace {
		return false
	}

	return reflect.DeepEqual(a.Data, b.Data)
}

func ownerReferencesDeepEqual(ors1, ors2 []kmetav1.OwnerReference) bool {
	if len(ors1) != len(ors2) {
		return false
	}

	var found bool
	for _, or1 := range ors1 {
		found = false
		for _, or2 := range ors2 {
			if reflect.DeepEqual(or1, or2) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	return true
}

// deploymentEqual asserts the equality of two Deployment objects.
func deploymentEqual(a, b *kappsv1.Deployment) bool {
	if a == nil || b == nil {
		return false
	}
	if a == b {
		return true
	}

	// compare objectMeta.
	objectMeta1 := a.ObjectMeta.DeepCopy()
	objectMeta2 := b.ObjectMeta.DeepCopy()
	if objectMeta1.Name != objectMeta2.Name ||
		objectMeta1.Namespace != objectMeta2.Namespace {
		return false
	}

	if !ownerReferencesDeepEqual(objectMeta1.OwnerReferences, objectMeta2.OwnerReferences) {
		return false
	}

	if !reflect.DeepEqual(a.Labels, b.Labels) {
		return false
	}

	// compare spec
	if !reflect.DeepEqual(a.Spec.Replicas, b.Spec.Replicas) ||
		!mapDeepEqual(a.Spec.Selector.MatchLabels, b.Spec.Selector.MatchLabels) ||
		a.Spec.MinReadySeconds != b.Spec.MinReadySeconds {
		return false
	}

	// compare pod template.
	cst1 := a.Spec.Template
	cst2 := b.Spec.Template
	if !mapDeepEqual(cst1.Annotations, cst2.Annotations) {
		return false
	}
	if !mapDeepEqual(cst1.Labels, cst2.Labels) {
		return false
	}

	ps1 := &cst1.Spec
	ps2 := &cst2.Spec
	return podSpecEqual(ps1, ps2)
}

// mapDeepEqual returns true if two non-empty maps are equal, otherwise returns false.
// If length of both maps evaluates to zero, it returns true.
func mapDeepEqual(m1, m2 map[string]string) bool {
	if len(m1) == 0 && len(m2) == 0 {
		return true
	}

	return reflect.DeepEqual(m1, m2)
}

// podSpecEqual asserts the equality of two PodSpec objects.
func podSpecEqual(ps1, ps2 *kcorev1.PodSpec) bool {
	if ps1 == nil || ps2 == nil {
		return false
	}
	if ps1 == ps2 {
		return true
	}

	cs1, cs2 := ps1.Containers, ps2.Containers
	if len(cs1) != len(cs2) {
		return false
	}
	for i := range cs1 {
		if !containerEqual(&cs1[i], &cs2[i]) {
			return false
		}
	}

	// compare volumeMounts.
	if !volumesEqual(ps1.Volumes, ps2.Volumes) {
		return false
	}

	return ps1.ServiceAccountName == ps2.ServiceAccountName
}

func volumesEqual(a, b []kcorev1.Volume) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func volumeMountsEqual(a, b []kcorev1.VolumeMount) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !reflect.DeepEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

// containerEqual asserts the equality of two Container objects.
func containerEqual(a, b *kcorev1.Container) bool {
	if a == nil || b == nil {
		return false
	}
	if a.Image != b.Image {
		return false
	}

	ps1, ps2 := a.Ports, b.Ports
	if len(ps1) != len(ps2) {
		return false
	}
	for i := range ps1 {
		isFound := false
		for j := range ps2 {
			p1, p2 := &ps1[i], &ps2[j]
			if p1.Name == p2.Name &&
				p1.ContainerPort == p2.ContainerPort &&
				realProto(p1.Protocol) == realProto(p2.Protocol) {
				isFound = true
				break
			}
		}
		if !isFound {
			return isFound
		}
	}

	if !envEqual(a.Env, b.Env) {
		return false
	}

	if !reflect.DeepEqual(a.Resources, b.Resources) {
		return false
	}

	if !volumeMountsEqual(a.VolumeMounts, b.VolumeMounts) {
		return false
	}

	return probeEqual(a.ReadinessProbe, b.ReadinessProbe)
}

// envEqual asserts the equality of two core environment slices. It's used
// by containerEqual.
func envEqual(a, b []kcorev1.EnvVar) bool {
	if len(a) != len(b) {
		return false
	}

	if len(a) == 0 {
		return true
	}
	isFound := false
	for _, ev1 := range a {
		isFound = false
		for _, ev2 := range b {
			if reflect.DeepEqual(ev1, ev2) {
				isFound = true
				break
			}
		}
		if !isFound {
			break
		}
	}
	return isFound
}

// probeEqual asserts the equality of two Probe objects. It's used by
// containerEqual.
func probeEqual(a, b *kcorev1.Probe) bool {
	if a == b {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	isInitialDelaySecondsEqual := a.InitialDelaySeconds != b.InitialDelaySeconds
	isTimeoutSecondsEqual := a.TimeoutSeconds != b.TimeoutSeconds && a.TimeoutSeconds != 0 && b.TimeoutSeconds != 0
	isPeriodSecondsEqual := a.PeriodSeconds != b.PeriodSeconds && a.PeriodSeconds != 0 && b.PeriodSeconds != 0
	isSuccessThresholdEqual := a.SuccessThreshold != b.SuccessThreshold &&
		a.SuccessThreshold != 0 && b.SuccessThreshold != 0
	isFailureThresholdEqual := a.FailureThreshold != b.FailureThreshold &&
		a.FailureThreshold != 0 && b.FailureThreshold != 0

	if isInitialDelaySecondsEqual || isTimeoutSecondsEqual || isPeriodSecondsEqual ||
		isSuccessThresholdEqual || isFailureThresholdEqual {
		return false
	}

	return handlerEqual(&a.ProbeHandler, &b.ProbeHandler)
}

// handlerEqual asserts the equality of two Handler objects. It's used
// by probeEqual.
func handlerEqual(a, b *kcorev1.ProbeHandler) bool {
	if a == b {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	hg1, hg2 := a.HTTPGet, b.HTTPGet
	if hg1 == nil && hg2 != nil {
		return false
	}
	if hg1 != nil {
		if hg2 == nil {
			return false
		}

		if hg1.Path != hg2.Path {
			return false
		}
	}

	return true
}

// realProto ensures the Protocol, which by default is TCP. We assume empty equals TCP.
// https://godoc.org/k8s.io/api/core/v1#ServicePort
func realProto(pr kcorev1.Protocol) kcorev1.Protocol {
	if pr == "" {
		return kcorev1.ProtocolTCP
	}
	return pr
}
