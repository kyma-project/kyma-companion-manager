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

	"go.uber.org/zap"
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	kapierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	"github.com/kyma-project/kyma-companion-manager/internal/backendmanager"
	"github.com/kyma-project/kyma-companion-manager/pkg/env"
	"github.com/kyma-project/kyma-companion-manager/pkg/equality"
	kcmk8s "github.com/kyma-project/kyma-companion-manager/pkg/k8s"
)

const (
	FinalizerName  = "companion.operator.kyma-project.io/finalizer"
	ControllerName = "kyma-companion-manager"
)

// Reconciler reconciles a Companion object.
type Reconciler struct {
	client.Client
	Scheme         *runtime.Scheme
	kubeClient     kcmk8s.Client
	logger         *zap.SugaredLogger
	backendManager backendmanager.Manager
	config         env.Config
}

func NewReconciler(
	client client.Client,
	kubeClient kcmk8s.Client,
	backendManager backendmanager.Manager,
	scheme *runtime.Scheme,
	logger *zap.SugaredLogger,
	config env.Config,
) *Reconciler {
	return &Reconciler{
		Client:         client,
		Scheme:         scheme,
		logger:         logger,
		backendManager: backendManager,
		config:         config,
		kubeClient:     kubeClient,
	}
}

// RBAC permissions.
//nolint:lll // ignore long line length due to kubebuilder markers.
// +kubebuilder:rbac:groups=operator.kyma-project.io,resources=companions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.kyma-project.io,resources=companions/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=operator.kyma-project.io,resources=companions/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
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

	//	reconcile secret of kyma-companion-backend.
	log.Info("reconciling secret...")
	err := r.reconcileSecret(ctx, companion, log)
	if err != nil {
		return kctrl.Result{}, err
	}

	//	reconcile deployment of kyma-companion-backend.
	log.Info("reconciling deployment...")
	err = r.reconcileDeployment(ctx, companion, log)
	if err != nil {
		return kctrl.Result{}, err
	}

	log.Info("companion reconciliation completed!")
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
		Owns(&kappsv1.Deployment{}). // watch for Deployments.
		Owns(&kcorev1.Secret{}).     // watch for Secrets.
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

func (r *Reconciler) reconcileDeployment(ctx context.Context, companion *kcmv1alpha1.Companion,
	log *zap.SugaredLogger,
) error {
	// define deployment object.
	expectedDeployment, err := r.backendManager.GenerateNewDeployment(companion, r.config.KymaCompanionBackendImage)
	if err != nil {
		return err
	}

	// fetch existing deployment.
	existingDeployment, err := r.kubeClient.GetDeployment(ctx, expectedDeployment.GetName(),
		expectedDeployment.GetNamespace())
	if err != nil && !kapierrors.IsNotFound(err) {
		return err
	}

	// compare if the deployment needs to be updated.
	if equality.Semantic.DeepEqual(existingDeployment, expectedDeployment) {
		log.Infof("deployment %s/%s already exists with expected configurations.",
			expectedDeployment.Namespace, expectedDeployment.Name)
		return nil
	}

	log.Infof("updating deployment %s/%s...", expectedDeployment.Namespace, expectedDeployment.Name)
	return r.kubeClient.PatchApply(ctx, expectedDeployment)
}

func (r *Reconciler) reconcileSecret(ctx context.Context, companion *kcmv1alpha1.Companion,
	log *zap.SugaredLogger,
) error {
	// get backend config.
	backendConfig, err := r.backendManager.GetBackendConfig(ctx)
	if err != nil {
		return err
	}

	// define secret.
	expectedSecret, err := r.backendManager.GenerateNewSecret(companion, *backendConfig)
	if err != nil {
		return err
	}

	// fetch existing secret.
	existingSecret, err := r.kubeClient.GetSecret(ctx, expectedSecret.GetName(),
		expectedSecret.GetNamespace())
	if err != nil && !kapierrors.IsNotFound(err) {
		return err
	}

	// compare if the secret needs to be updated.
	if equality.Semantic.DeepEqual(existingSecret, expectedSecret) {
		log.Infof("secret %s/%s already exists with expected data.",
			expectedSecret.Namespace, expectedSecret.Name)
		return nil
	}

	log.Infof("updating secret %s/%s...", expectedSecret.Namespace, expectedSecret.Name)
	return r.kubeClient.PatchApply(ctx, expectedSecret)
}
