package k8s_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCmdExecutor is a mock implementation of command executor
type MockCmdExecutor struct {
	mock.Mock
}

func (m *MockCmdExecutor) Command(name string, arg ...string) *exec.Cmd {
	args := m.Called(name, arg)
	return args.Get(0).(*exec.Cmd)
}

// KubernetesDeployerWithMock allows to inject mock command executor
type KubernetesDeployerWithMock struct {
	cmdExecutor func(name string, arg ...string) *exec.Cmd
}

// NewKubernetesDeployerWithMock creates a new KubernetesDeployer with a mock executor
func NewKubernetesDeployerWithMock(executor func(name string, arg ...string) *exec.Cmd) *KubernetesDeployerWithMock {
	return &KubernetesDeployerWithMock{cmdExecutor: executor}
}

func (k *KubernetesDeployerWithMock) CreatePod(name string, image string) error {
	cmd := k.cmdExecutor("kubectl", "run", name, "--image="+image)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "already exists") {
			return nil
		}
		return fmt.Errorf("failed to create pod: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

func (k *KubernetesDeployerWithMock) DeletePod(name string) error {
	cmd := k.cmdExecutor("kubectl", "delete", "pod", name)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if strings.Contains(stderr.String(), "NotFound") {
			return nil
		}
		return fmt.Errorf("failed to delete pod: %w, stderr: %s", err, stderr.String())
	}
	return nil
}

func (k *KubernetesDeployerWithMock) GetAllPodList() ([]string, error) {
	cmd := k.cmdExecutor("kubectl", "get", "pods", "-o", "json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get pods: %w", err)
	}

	var result struct {
		Items []struct {
			Metadata struct {
				Name string `json:"name"`
			} `json:"metadata"`
		} `json:"items"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	podNames := make([]string, len(result.Items))
	for i, item := range result.Items {
		podNames[i] = item.Metadata.Name
	}

	return podNames, nil
}

func TestKubernetesDeployer_CreatePod(t *testing.T) {
	mockCmd := new(MockCmdExecutor)
	mockCmd.On("Command", "kubectl", []string{"run", "test-pod", "--image=test-image"}).Return(exec.Command("echo", "already exists"))

	deployer := NewKubernetesDeployerWithMock(mockCmd.Command)
	err := deployer.CreatePod("test-pod", "test-image")

	assert.NoError(t, err)
	mockCmd.AssertExpectations(t)
}

func TestKubernetesDeployer_DeletePod(t *testing.T) {
	mockCmd := new(MockCmdExecutor)
	mockCmd.On("Command", "kubectl", []string{"delete", "pod", "test-pod"}).Return(exec.Command("echo", "NotFound"))

	deployer := NewKubernetesDeployerWithMock(mockCmd.Command)
	err := deployer.DeletePod("test-pod")

	assert.NoError(t, err)
	mockCmd.AssertExpectations(t)
}

func TestKubernetesDeployer_GetAllPodList(t *testing.T) {
	mockCmd := new(MockCmdExecutor)
	podList := `{"items":[{"metadata":{"name":"pod1"}},{"metadata":{"name":"pod2"}}]}`
	mockCmd.On("Command", "kubectl", []string{"get", "pods", "-o", "json"}).Return(exec.Command("echo", podList))

	deployer := NewKubernetesDeployerWithMock(mockCmd.Command)
	pods, err := deployer.GetAllPodList()

	assert.NoError(t, err)
	assert.Equal(t, []string{"pod1", "pod2"}, pods)
	mockCmd.AssertExpectations(t)
}
