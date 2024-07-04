package utils

import (
	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type CompanionOption func(*kcmv1alpha1.Companion) error

func WithCompanionCRFinalizer(finalizer string) CompanionOption {
	return func(e *kcmv1alpha1.Companion) error {
		controllerutil.AddFinalizer(e, finalizer)
		return nil
	}
}
