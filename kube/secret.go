package kube

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	core "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// @see https://github.com/kubernetes/kubectl/blob/c1d804c03fefc4de986592e9db203a306467e65d/pkg/generate/versioned/secret_for_docker_registry.go
// @see https://github.com/kubernetes/kubernetes/blob/b2ecd1b3a3192fbbe2b9e348e095326f51dc43dd/pkg/apis/core/types.go#L4938-L4961
// @see https://github.com/kubernetes/client-go/blob/master/kubernetes/typed/core/v1/secret.go
// @see https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/cri-api/pkg/errors/errors.go

// DockerConfigJSON represents a local docker auth config file for pulling images.
type DockerConfigJSON struct {
	Auths DockerConfig `json:"auths" datapolicy:"token"`
}

// DockerConfig represents the config file used by the docker CLI.
// This config that represents the credentials that should be used
// when pulling images from specific image repositories.
type DockerConfig map[string]DockerConfigEntry

// DockerConfigEntry is a entry for a config of docker registry
type DockerConfigEntry struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty" datapolicy:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth,omitempty" datapolicy:"token"`
}

// Client is a client for Kubernetes APIs
type Client struct {
	set       kubernetes.Interface
	namespace string
}

// NewClient is a constructor
func NewClient(masterURL, kubeconfig, namespace string) (*Client, error) {
	if namespace == "" {
		return nil, fmt.Errorf("empty namespace was given")
	}

	cfg, err := buildConfig(masterURL, kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("Error building kubeconfig: %w", err)
	}

	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("Error building kubernetes clientset: %w", err)
	}

	return &Client{set: cli, namespace: namespace}, nil
}

func buildConfig(masterURL, kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" {
		return rest.InClusterConfig()
	}

	return clientcmd.BuildConfigFromFlags(masterURL, kubeconfig)
}

// UpdateSecret is a updater for Kubernetes secret of docker-registry type
func (c *Client) UpdateSecret(name, server, user, password, email string) error {
	if err := c.deleteSecret(name); err != nil {
		return err
	}

	if err := c.createSecret(name, server, user, password, email); err != nil {
		return err
	}

	return nil
}

func (c *Client) deleteSecret(name string) error {
	if err := c.set.CoreV1().Secrets(c.namespace).Delete(context.TODO(), name, meta.DeleteOptions{}); apierrors.IsNotFound(err) == false {
		return err
	}

	return nil
}

func (c *Client) createSecret(name, server, user, password, email string) error {
	secret := &core.Secret{}
	secret.Name = name
	secret.Type = core.SecretTypeDockerConfigJson
	secret.Data = map[string][]byte{}

	data, err := encodeSecretData(server, user, password, email)
	if err != nil {
		return err
	}

	secret.Data[core.DockerConfigJsonKey] = data

	if _, err = c.set.CoreV1().Secrets(c.namespace).Create(context.TODO(), secret, meta.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func encodeSecretData(server, user, password, email string) ([]byte, error) {
	field := fmt.Sprintf("%s:%s", user, password)
	auth := base64.StdEncoding.EncodeToString([]byte(field))

	entry := DockerConfigEntry{
		Username: user,
		Password: password,
		Email:    email,
		Auth:     auth,
	}

	body := DockerConfigJSON{
		Auths: DockerConfig{server: entry},
	}

	return json.Marshal(body)
}
