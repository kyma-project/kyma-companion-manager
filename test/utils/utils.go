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

package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	kappsv1 "k8s.io/api/apps/v1"
	kcorev1 "k8s.io/api/core/v1"
	kmetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kcmv1alpha1 "github.com/kyma-project/kyma-companion-manager/api/v1alpha1"
	kcmlabel "github.com/kyma-project/kyma-companion-manager/internal/label"
	kcmk8sdeployment "github.com/kyma-project/kyma-companion-manager/pkg/k8s/deployment"

	. "github.com/onsi/ginkgo/v2"
)

const (
	prometheusOperatorVersion = "v0.72.0"
	prometheusOperatorURL     = "https://github.com/prometheus-operator/prometheus-operator/" +
		"releases/download/%s/bundle.yaml"

	certmanagerVersion = "v1.14.4"
	certmanagerURLTmpl = "https://github.com/jetstack/cert-manager/releases/download/%s/cert-manager.yaml"

	randomNameLen = 5
	charset       = "abcdefghijklmnopqrstuvwxyz0123456789"

	NameFormat      = "name-%s"
	NamespaceFormat = "namespace-%s"
)

var (
	seededRand            = rand.New(rand.NewSource(time.Now().UnixNano())) //nolint:gosec,gochecknoglobals,lll // used in tests
	ErrPortCreationFailed = errors.New("failed to get port")
)

func warnError(err error) {
	fmt.Fprintf(GinkgoWriter, "warning: %v\n", err)
}

// InstallPrometheusOperator installs the prometheus Operator to be used to export the enabled metrics.
func InstallPrometheusOperator() error {
	url := fmt.Sprintf(prometheusOperatorURL, prometheusOperatorVersion)
	cmd := exec.Command("kubectl", "create", "-f", url)
	_, err := Run(cmd)
	return err
}

// Run executes the provided command within this context.
func Run(cmd *exec.Cmd) ([]byte, error) {
	dir, _ := GetProjectDir()
	cmd.Dir = dir

	if err := os.Chdir(cmd.Dir); err != nil {
		fmt.Fprintf(GinkgoWriter, "chdir dir: %s\n", err)
	}

	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	command := strings.Join(cmd.Args, " ")
	fmt.Fprintf(GinkgoWriter, "running: %s\n", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		//nolint:errorlint,goerr113 // used for testing purposes.
		return output, fmt.Errorf("%s failed with error: (%v) %s", command, err, string(output))
	}

	return output, nil
}

// UninstallPrometheusOperator uninstalls the prometheus.
func UninstallPrometheusOperator() {
	url := fmt.Sprintf(prometheusOperatorURL, prometheusOperatorVersion)
	cmd := exec.Command("kubectl", "delete", "-f", url)
	if _, err := Run(cmd); err != nil {
		warnError(err)
	}
}

// UninstallCertManager uninstalls the cert manager.
func UninstallCertManager() {
	url := fmt.Sprintf(certmanagerURLTmpl, certmanagerVersion)
	cmd := exec.Command("kubectl", "delete", "-f", url)
	if _, err := Run(cmd); err != nil {
		warnError(err)
	}
}

// InstallCertManager installs the cert manager bundle.
func InstallCertManager() error {
	url := fmt.Sprintf(certmanagerURLTmpl, certmanagerVersion)
	cmd := exec.Command("kubectl", "apply", "-f", url)
	if _, err := Run(cmd); err != nil {
		return err
	}
	// Wait for cert-manager-webhook to be ready, which can take time if cert-manager
	// was re-installed after uninstalling on a cluster.
	cmd = exec.Command("kubectl", "wait", "deployment.apps/cert-manager-webhook",
		"--for", "condition=Available",
		"--namespace", "cert-manager",
		"--timeout", "5m",
	)

	_, err := Run(cmd)
	return err
}

// LoadImageToKindClusterWithName loads a local docker image to the kind cluster.
func LoadImageToKindClusterWithName(name string) error {
	cluster := "kind"
	if v, ok := os.LookupEnv("KIND_CLUSTER"); ok {
		cluster = v
	}
	kindOptions := []string{"load", "docker-image", name, "--name", cluster}
	cmd := exec.Command("kind", kindOptions...)
	_, err := Run(cmd)
	return err
}

// GetNonEmptyLines converts given command output string into individual objects
// according to line breakers, and ignores the empty elements in it.
func GetNonEmptyLines(output string) []string {
	var res []string
	elements := strings.Split(output, "\n")
	for _, element := range elements {
		if element != "" {
			res = append(res, element)
		}
	}

	return res
}

// GetProjectDir will return the directory where the project is.
func GetProjectDir() (string, error) {
	wdir, err := os.Getwd()
	if err != nil {
		return wdir, err
	}
	//nolint:gocritic // we need to remove the last part of the path.
	wdir = strings.Replace(wdir, "/test/e2e", "", -1)
	return wdir, nil
}

func NewLogger() (*zap.Logger, error) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.Encoding = "json"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("Jan 02 15:04:05.000000000")

	return loggerConfig.Build()
}

func NewSugaredLogger() (*zap.SugaredLogger, error) {
	logger, err := NewLogger()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}

func GetRandString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// GetFreePort determines a free port on the host. It does so by delegating the job to net.ListenTCP.
// Then providing a port of 0 to net.ListenTCP, it will automatically choose a port for us.
func GetFreePort() (int, error) {
	if a, err := net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var listener *net.TCPListener
		if listener, err = net.ListenTCP("tcp", a); err == nil {
			portAddr, ok := listener.Addr().(*net.TCPAddr)
			if !ok {
				return -1, ErrPortCreationFailed
			}
			port := portAddr.Port
			err = listener.Close()
			return port, err
		}
	}
	return -1, nil
}

func GetRandomName() string {
	return fmt.Sprintf(NameFormat, GetRandString(randomNameLen))
}

func GetRandomNamespaceName() string {
	return fmt.Sprintf(NamespaceFormat, GetRandString(randomNameLen))
}

func NewNamespace(name string) *kcorev1.Namespace {
	namespace := kcorev1.Namespace{
		TypeMeta: kmetav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: kmetav1.ObjectMeta{
			Name: name,
		},
	}
	return &namespace
}

func NewCompanionCR(opts ...CompanionOption) *kcmv1alpha1.Companion {
	name := GetRandomName()
	namespace := GetRandomNamespaceName()

	eventing := &kcmv1alpha1.Companion{
		TypeMeta: kmetav1.TypeMeta{
			Kind:       "Companion",
			APIVersion: "operator.kyma-project.io/v1alpha1",
		},
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       "1234-5678-1234-5678",
		},
		Spec: kcmv1alpha1.CompanionSpec{},
	}

	for _, opt := range opts {
		if err := opt(eventing); err != nil {
			panic(err)
		}
	}

	return eventing
}

func NewDeployment(name, namespace string, annotations map[string]string) *kappsv1.Deployment {
	return &kappsv1.Deployment{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Annotations: annotations,
		},
		Spec: kappsv1.DeploymentSpec{
			Template: kcorev1.PodTemplateSpec{
				ObjectMeta: kmetav1.ObjectMeta{
					Name:        name,
					Namespace:   namespace,
					Annotations: annotations,
				},
				Spec: kcorev1.PodSpec{
					Containers: []kcorev1.Container{
						{
							Name:  "companion",
							Image: "test-image",
						},
					},
				},
			},
		},
	}
}

func NewCompanionDeployment(name, namespace string) *kappsv1.Deployment {
	// define labels.
	labels := kcmlabel.GetCommonLabels(name)

	// define containers.
	containers := []kcorev1.Container{
		{
			Name:            kcmlabel.ValueCompanionBackend,
			Image:           "test-image:latest",
			ImagePullPolicy: kcorev1.PullAlways,
		},
	}

	// define deployment object.
	deployment := kcmk8sdeployment.NewDeployment(
		name,
		namespace,
		kcmk8sdeployment.WithLabels(labels),
		kcmk8sdeployment.WithRestartPolicyAlways(),
		kcmk8sdeployment.WithSelectorLabels(labels),
		kcmk8sdeployment.WithContainers(containers),
	)

	return deployment
}

func NewConfigMap(name, namespace string) *kcorev1.ConfigMap {
	configMap := &kcorev1.ConfigMap{
		ObjectMeta: kmetav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}
	return configMap
}
