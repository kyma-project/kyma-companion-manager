package controller

import (
	"context"

	kctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
)

func (r *Reconciler) containsFinalizer(companion *kcmv1alpha1.Companion) bool {
	return controllerutil.ContainsFinalizer(companion, FinalizerName)
}

func (r *Reconciler) addFinalizer(ctx context.Context, companion *kcmv1alpha1.Companion) (kctrl.Result, error) {
	controllerutil.AddFinalizer(companion, FinalizerName)
	if err := r.Update(ctx, companion); err != nil {
		return kctrl.Result{}, err
	}
	return kctrl.Result{}, nil
}

func (r *Reconciler) removeFinalizer(ctx context.Context, companion *kcmv1alpha1.Companion) (kctrl.Result, error) {
	controllerutil.RemoveFinalizer(companion, FinalizerName)
	if err := r.Update(ctx, companion); err != nil {
		return kctrl.Result{}, err
	}

	return kctrl.Result{}, nil
}
