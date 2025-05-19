package kubernetes

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var modelRegistryGVR = schema.GroupVersionResource{
	Group:    "modelregistry.opendatahub.io",
	Version:  "v1alpha1",
	Resource: "modelregistries",
}

// RHOAITokenKubernetesClient wraps the default TokenKubernetesClient but overrides selected methods.
// It also adds a dynamic client to the base client.
type RHOAITokenKubernetesClient struct {
	*kubernetes.TokenKubernetesClient
	DynClient dynamic.Interface
}

func NewRHOAIKubernetesClient(token string, logger *slog.Logger) (kubernetes.KubernetesClientInterface, error) {
	baseClient, err := kubernetes.NewTokenKubernetesClient(token, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create base token client: %w", err)
	}

	typed := baseClient.(*kubernetes.TokenKubernetesClient)

	dynClient, err := dynamic.NewForConfig(typed.RESTConfig())

	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	return &RHOAITokenKubernetesClient{
		TokenKubernetesClient: baseClient.(*kubernetes.TokenKubernetesClient),
		DynClient:             dynClient,
	}, nil
}

// GetModelRegistrySettings overrides the default stub logic with RHOAI-specific logic.
func (c *RHOAITokenKubernetesClient) GetModelRegistrySettings(ctx context.Context, namespace string, labelSelector string) ([]unstructured.Unstructured, error) {

	listOptions := v1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}

	list, err := c.DynClient.
		Resource(modelRegistryGVR).
		Namespace(namespace).
		List(ctx, listOptions)
	if err != nil {
		c.Logger.Error("Failed to list ModelRegistry CRs", "error", err)
		return nil, fmt.Errorf("failed to list ModelRegistry CRs in namespace %q: %w", namespace, err)
	}

	return list.Items, nil
}

func (c *RHOAITokenKubernetesClient) GetModelRegistrySettingsByName(ctx context.Context, namespace string, name string) (unstructured.Unstructured, error) {

	result, err := c.DynClient.
		Resource(modelRegistryGVR).
		Namespace(namespace).
		Get(ctx, name, v1.GetOptions{})
	if err != nil {
		c.Logger.Error("Failed to get ModelRegistry CR", "error", err, "name", name, "namespace", namespace)
		return unstructured.Unstructured{}, fmt.Errorf("failed to get ModelRegistry %q in namespace %q: %w", name, namespace, err)
	}

	return *result, nil
}

func (c *RHOAITokenKubernetesClient) CreateDatabaseSecret(ctx context.Context, name string, namespace string, database string, databaseUsername string, databasePassword string, dryRun bool) (*corev1.Secret, error) {

	if database == "" || databaseUsername == "" {
		return nil, fmt.Errorf("invalid database config: database name or username missing")
	}

	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			GenerateName: name + "-db-",
			Namespace:    namespace,
			Annotations: map[string]string{
				"template.openshift.io/expose-database_name": "{.data['database-name']}",
				"template.openshift.io/expose-username":      "{.data['database-user']}",
				"template.openshift.io/expose-password":      "{.data['database-password']}",
			},
			Labels: map[string]string{
				"modelregistry.opendatahub.io/managed": "true",
				"modelregistry.opendatahub.io/name":    name,
			},
		},
		StringData: map[string]string{
			"database-name":     database,
			"database-user":     databaseUsername,
			"database-password": databasePassword,
		},
		Type: corev1.SecretTypeOpaque,
	}

	options := v1.CreateOptions{}
	if dryRun {
		options.DryRun = []string{v1.DryRunAll}
	}

	created, err := c.Client.CoreV1().Secrets(namespace).Create(ctx, secret, options)
	if err != nil {
		return nil, fmt.Errorf("failed to create secret: %w", err)
	}

	return created, nil
}

func (kc *RHOAITokenKubernetesClient) CreateModelRegistryKind(ctx context.Context, namespace string, modelRegistryKind unstructured.Unstructured, dryRun bool) (unstructured.Unstructured, error) {

	gvr := schema.GroupVersionResource{
		Group:    "modelregistry.opendatahub.io",
		Version:  "v1alpha1",
		Resource: "modelregistries", // plural, always lowercase
	}

	opts := v1.CreateOptions{}
	if dryRun {
		opts.DryRun = []string{v1.DryRunAll}
	}

	created, err := kc.DynClient.
		Resource(gvr).
		Namespace(namespace).
		Create(ctx, &modelRegistryKind, opts)

	if err != nil {
		return unstructured.Unstructured{}, fmt.Errorf("failed to create ModelRegistryKind: %w", err)
	}

	return *created, nil
}
