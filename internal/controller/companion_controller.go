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

package controller

import (
	"context"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/runtime"
	kctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	FinalizerName  = "companion.operator.kyma-project.io/finalizer"
	ControllerName = "kyma-companion-manager-controller"
)

// Reconciler reconciles a Companion object
type Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
	logger *zap.SugaredLogger
}

func NewReconciler(
	client client.Client,
	scheme *runtime.Scheme,
	logger *zap.SugaredLogger,
) *Reconciler {
	return &Reconciler{
		Client: client,
		Scheme: scheme,
		logger: logger,
	}
}

// +kubebuilder:rbac:groups=operator.kyma-project.io,resources=companions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.kyma-project.io,resources=companions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=operator.kyma-project.io,resources=companions/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Companion object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.2/pkg/reconcile
func (r *Reconciler) Reconcile(ctx context.Context, req kctrl.Request) (kctrl.Result, error) {
	r.logger.Info("Reconciliation triggered")

	// fetch latest CR.
	currentCompanion := &kcmv1alpha1.Companion{}
	if err := r.Get(ctx, req.NamespacedName, currentCompanion); err != nil {
		return kctrl.Result{}, client.IgnoreNotFound(err)
	}

	// copy the object, so we don't modify the source object.
	companionCR := currentCompanion.DeepCopy()

	// Create a logger with NATS details.
	log := r.loggerWithCompanion(companionCR)

	// check if companion CR is in deletion state.
	if !companionCR.DeletionTimestamp.IsZero() {
		return r.handleCompanionDeletion(ctx, companionCR, log)
	}

	// handle reconciliation.
	return r.handleCompanionReconcile(ctx, companionCR, log)
}

func (r *Reconciler) handleCompanionReconcile(ctx context.Context,
	companion *kcmv1alpha1.Companion, log *zap.SugaredLogger,
) (kctrl.Result, error) {
	log.Info("handling Companion reconciliation...")

	// make sure the finalizer exists.
	if !r.containsFinalizer(companion) {
		return r.addFinalizer(ctx, companion)
	}

	log.Info("dummy Companion reconciliation completed!")
	return kctrl.Result{}, nil
}

func (r *Reconciler) handleCompanionDeletion(ctx context.Context, companion *kcmv1alpha1.Companion,
	log *zap.SugaredLogger,
) (kctrl.Result, error) {
	// skip reconciliation for deletion if the finalizer is not set.
	if !r.containsFinalizer(companion) {
		log.Info("skipped reconciliation for deletion as finalizer is not set.")
		return kctrl.Result{}, nil
	}

	log.Info("handling Companion deletion...")
	return r.removeFinalizer(ctx, companion)
}

// SetupWithManager sets up the controller with the Manager.
func (r *Reconciler) SetupWithManager(mgr kctrl.Manager) error {
	return kctrl.NewControllerManagedBy(mgr).
		For(&kcmv1alpha1.Companion{}).
		Complete(r)
}

// loggerWithCompanion returns a logger with the given Companion CR details.
func (r *Reconciler) loggerWithCompanion(companion *kcmv1alpha1.Companion) *zap.SugaredLogger {
	return r.logger.With(
		"kind", companion.GetObjectKind().GroupVersionKind().Kind,
		"resourceVersion", companion.GetResourceVersion(),
		"generation", companion.GetGeneration(),
		"namespace", companion.GetNamespace(),
		"name", companion.GetName(),
	)
}
