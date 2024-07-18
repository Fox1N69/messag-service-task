package k8s

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type KubernetesDeployer interface {
	CreatePod(name string, image string) error
	DeletePod(name string) error
	GetPodList() ([]string, error)
}

type kubernetesDeployer struct{}

func NewKubernetesDeployer() KubernetesDeployer {
	return &kubernetesDeployer{}
}

// CreatePod created pod by name and with the image of the client
func (k *kubernetesDeployer) CreatePod(name string, image string) error {
	cmd := exec.Command("kubectl", "run", name, "--image="+image)
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

// DeletePod deleted pod by name
func (k *kubernetesDeployer) DeletePod(name string) error {
	cmd := exec.Command("kubectl", "delete", "pod", name)
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

// GetAllPodList returns all list pods in the kubernetes
func (k *kubernetesDeployer) GetPodList() ([]string, error) {
	cmd := exec.Command("kubectl", "get", "pods", "-o", "json")
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
