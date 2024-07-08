package integration

import (
	"context"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/avast/retry-go/v3"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	kcorev1 "k8s.io/api/core/v1"
	kapierrors "k8s.io/apimachinery/pkg/api/errors"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	ktypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/dynamic"
	kkubernetesscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	kctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	kctrllogzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	kcmctrl "github.com/kyma-project/kyma-companion-manager/internal/controller"
	testutils "github.com/kyma-project/kyma-companion-manager/test/utils"
)

const (
	useExistingCluster       = false
	attachControlPlaneOutput = false
	testEnvStartDelay        = time.Minute
	testEnvStartAttempts     = 10
	namespacePrefixLength    = 5
	BigPollingInterval       = 3 * time.Second
	BigTimeOut               = 40 * time.Second
	SmallTimeOut             = 5 * time.Second
	SmallPollingInterval     = 1 * time.Second
)

// TestEnvironment provides mocked resources for integration tests.
type TestEnvironment struct {
	EnvTestInstance  *envtest.Environment
	k8sClient        client.Client
	K8sDynamicClient *dynamic.DynamicClient
	Reconciler       *kcmctrl.Reconciler
	Logger           *zap.SugaredLogger
	Recorder         *record.EventRecorder
	TestCancelFn     context.CancelFunc
}

//nolint:funlen // Used in testing
func NewTestEnvironment(projectRootDir string, allowedCompanionCR *kcmv1alpha1.Companion,
) (*TestEnvironment, error) {
	var err error

	// setup logger
	sugaredLogger, err := testutils.NewSugaredLogger()
	if err != nil {
		return nil, err
	}
	kctrl.SetLogger(kctrllogzap.New())

	testEnv, envTestKubeCfg, err := StartEnvTest(projectRootDir)
	if err != nil {
		return nil, err
	}

	// add to Scheme
	err = kcmv1alpha1.AddToScheme(kkubernetesscheme.Scheme)
	if err != nil {
		return nil, err
	}

	// +kubebuilder:scaffold:scheme

	k8sClient, err := client.New(envTestKubeCfg, client.Options{Scheme: kkubernetesscheme.Scheme})
	if err != nil {
		return nil, err
	}

	dynamicClient, err := dynamic.NewForConfig(envTestKubeCfg)
	if err != nil {
		return nil, err
	}

	// setup ctrl manager
	metricsPort, err := testutils.GetFreePort()
	if err != nil {
		return nil, err
	}

	ctrlMgr, err := kctrl.NewManager(envTestKubeCfg, kctrl.Options{
		Scheme:                 kkubernetesscheme.Scheme,
		HealthProbeBindAddress: "0",                              // disable
		PprofBindAddress:       "0",                              // disable
		Metrics:                server.Options{BindAddress: "0"}, // disable
		WebhookServer:          webhook.NewServer(webhook.Options{Port: metricsPort}),
	})
	if err != nil {
		return nil, err
	}
	recorder := ctrlMgr.GetEventRecorderFor("kyma-companion-manager")

	// setup reconciler
	kcmReconciler := kcmctrl.NewReconciler(
		ctrlMgr.GetClient(),
		ctrlMgr.GetScheme(),
		sugaredLogger,
	)
	if err = (kcmReconciler).SetupWithManager(ctrlMgr); err != nil {
		return nil, err
	}

	// start manager
	var cancelCtx context.CancelFunc
	go func() {
		var mgrCtx context.Context
		mgrCtx, cancelCtx = context.WithCancel(kctrl.SetupSignalHandler())
		err = ctrlMgr.Start(mgrCtx)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return &TestEnvironment{
		k8sClient:        k8sClient,
		K8sDynamicClient: dynamicClient,
		Reconciler:       kcmReconciler,
		Logger:           sugaredLogger,
		Recorder:         &recorder,
		EnvTestInstance:  testEnv,
		TestCancelFn:     cancelCtx,
	}, nil
}

func (env TestEnvironment) TearDown() error {
	if env.TestCancelFn != nil {
		env.TestCancelFn()
	}

	// retry to stop the api-server
	sleepTime := 1 * time.Second
	var err error
	const retries = 20
	for range retries {
		if err = env.EnvTestInstance.Stop(); err == nil {
			break
		}
		time.Sleep(sleepTime)
	}
	return err
}

func (env TestEnvironment) CreateK8sResource(ctx context.Context, obj client.Object) error {
	return env.k8sClient.Create(ctx, obj)
}

func (env TestEnvironment) EnsureNamespaceCreation(t *testing.T, ctx context.Context, namespace string) {
	t.Helper()
	if namespace == "default" {
		return
	}
	// create namespace
	ns := testutils.NewNamespace(namespace)
	require.NoError(t, client.IgnoreAlreadyExists(env.k8sClient.Create(ctx, ns)))
}

func (env TestEnvironment) EnsureK8sResourceCreated(t *testing.T, ctx context.Context, obj client.Object) {
	t.Helper()
	require.NoError(t, env.k8sClient.Create(ctx, obj))
}

func (env TestEnvironment) EnsureK8sUnStructResourceCreated(t *testing.T, ctx context.Context,
	obj *unstructured.Unstructured,
) {
	t.Helper()
	require.NoError(t, env.k8sClient.Create(ctx, obj))
}

func (env TestEnvironment) CreateUnstructuredK8sResource(ctx context.Context, obj *unstructured.Unstructured) error {
	return env.k8sClient.Create(ctx, obj)
}

func (env TestEnvironment) EnsureK8sResourceUpdated(t *testing.T, ctx context.Context, obj client.Object) {
	t.Helper()
	require.NoError(t, env.k8sClient.Update(ctx, obj))
}

func (env TestEnvironment) EnsureK8sResourceDeleted(t *testing.T, ctx context.Context, obj client.Object) {
	t.Helper()
	require.NoError(t, env.k8sClient.Delete(ctx, obj))
}

func (env TestEnvironment) EnsureK8sConfigMapExists(t *testing.T, ctx context.Context, name, namespace string) {
	t.Helper()
	require.Eventually(t, func() bool {
		result, err := env.GetConfigMapFromK8s(ctx, name, namespace)
		return err == nil && result != nil
	}, SmallTimeOut, SmallPollingInterval, "failed to ensure existence of ConfigMap")
}

func (env TestEnvironment) EnsureK8sSecretExists(t *testing.T, ctx context.Context, name, namespace string) {
	t.Helper()
	require.Eventually(t, func() bool {
		result, err := env.GetSecretFromK8s(ctx, name, namespace)
		return err == nil && result != nil
	}, SmallTimeOut, SmallPollingInterval, "failed to ensure existence of Secret")
}

func (env TestEnvironment) EnsureK8sServiceExists(t *testing.T, ctx context.Context, name, namespace string) {
	t.Helper()
	require.Eventually(t, func() bool {
		result, err := env.GetServiceFromK8s(ctx, name, namespace)
		return err == nil && result != nil
	}, SmallTimeOut, SmallPollingInterval, "failed to ensure existence of Service")
}

func (env TestEnvironment) EnsureK8sConfigMapNotFound(t *testing.T, ctx context.Context, name, namespace string) {
	t.Helper()
	require.Eventually(t, func() bool {
		_, err := env.GetConfigMapFromK8s(ctx, name, namespace)
		return err != nil && kapierrors.IsNotFound(err)
	}, SmallTimeOut, SmallPollingInterval, "failed to ensure non-existence of ConfigMap")
}

func (env TestEnvironment) EnsureK8sSecretNotFound(t *testing.T, ctx context.Context, name, namespace string) {
	t.Helper()
	require.Eventually(t, func() bool {
		_, err := env.GetSecretFromK8s(ctx, name, namespace)
		return err != nil && kapierrors.IsNotFound(err)
	}, SmallTimeOut, SmallPollingInterval, "failed to ensure non-existence of Secret")
}

func (env TestEnvironment) EnsureK8sServiceNotFound(t *testing.T, ctx context.Context, name, namespace string) {
	t.Helper()
	require.Eventually(t, func() bool {
		_, err := env.GetServiceFromK8s(ctx, name, namespace)
		return err != nil && kapierrors.IsNotFound(err)
	}, SmallTimeOut, SmallPollingInterval, "failed to ensure non-existence of Service")
}

func (env TestEnvironment) GetConfigMapFromK8s(ctx context.Context,
	name, namespace string,
) (*kcorev1.ConfigMap, error) {
	nn := ktypes.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	result := &kcorev1.ConfigMap{}
	if err := env.k8sClient.Get(ctx, nn, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (env TestEnvironment) GetSecretFromK8s(ctx context.Context, name, namespace string) (*kcorev1.Secret, error) {
	nn := ktypes.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	result := &kcorev1.Secret{}
	if err := env.k8sClient.Get(ctx, nn, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (env TestEnvironment) GetServiceFromK8s(ctx context.Context, name, namespace string) (*kcorev1.Service, error) {
	nn := ktypes.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}
	result := &kcorev1.Service{}
	if err := env.k8sClient.Get(ctx, nn, result); err != nil {
		return nil, err
	}
	return result, nil
}

func (env TestEnvironment) GetCompanionCRFromK8s(ctx context.Context,
	name, namespace string,
) (kcmv1alpha1.Companion, error) {
	var companion kcmv1alpha1.Companion
	err := env.k8sClient.Get(ctx, ktypes.NamespacedName{
		Name:      name,
		Namespace: namespace,
	}, &companion)
	return companion, err
}

func (env TestEnvironment) DeleteServiceFromK8s(ctx context.Context, name, namespace string) error {
	return env.k8sClient.Delete(ctx, &kcorev1.Service{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	})
}

func (env TestEnvironment) DeleteConfigMapFromK8s(ctx context.Context, name, namespace string) error {
	return env.k8sClient.Delete(ctx, &kcorev1.ConfigMap{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	})
}

func (env TestEnvironment) DeleteSecretFromK8s(ctx context.Context, name, namespace string) error {
	return env.k8sClient.Delete(ctx, &kcorev1.Secret{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	})
}

func StartEnvTest(projectRootDir string) (*envtest.Environment, *rest.Config, error) {
	// Reference: https://book.kubebuilder.io/reference/envtest.html
	useExistingCluster := useExistingCluster
	testEnv := &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join(projectRootDir, "config", "crd", "bases"),
		},
		ErrorIfCRDPathMissing:    false,
		AttachControlPlaneOutput: attachControlPlaneOutput,
		UseExistingCluster:       &useExistingCluster,
	}

	var cfg *rest.Config
	err := retry.Do(func() error {
		defer func() {
			if r := recover(); r != nil {
				log.Println("panic recovered:", r)
			}
		}()
		cfgLocal, startErr := testEnv.Start()
		cfg = cfgLocal
		return startErr
	},
		retry.Delay(testEnvStartDelay),
		retry.DelayType(retry.FixedDelay),
		retry.Attempts(testEnvStartAttempts),
		retry.OnRetry(func(n uint, err error) {
			log.Printf("[%v] try failed to start testenv: %s", n, err)
			if stopErr := testEnv.Stop(); stopErr != nil {
				log.Printf("failed to stop testenv: %s", stopErr)
			}
		}),
	)
	return testEnv, cfg, err
}
