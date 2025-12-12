package config

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/smartcontractkit/chainlink/system-tests/lib/infra"
)

// createOrUpdateResource is a generic helper for creating or updating Kubernetes resources
func createOrUpdateResource(ctx context.Context, clientset *kubernetes.Clientset, namespace, name string, obj metav1.Object, logger zerolog.Logger, resourceType string) error {
	obj.SetNamespace(namespace)
	obj.SetName(name)

	// Try to create first (most common case)
	switch o := obj.(type) {
	case *v1.Secret:
		if _, err := clientset.CoreV1().Secrets(namespace).Create(ctx, o, metav1.CreateOptions{}); err != nil {
			if !errors.IsAlreadyExists(err) {
				return fmt.Errorf("failed to create %s %s: %w", resourceType, name, err)
			}
			// Already exists, try to update
			if _, err := clientset.CoreV1().Secrets(namespace).Update(ctx, o, metav1.UpdateOptions{}); err != nil {
				return fmt.Errorf("failed to update %s %s: %w", resourceType, name, err)
			}
			logger.Info().Str("name", name).Str("namespace", namespace).Msgf("%s updated", resourceType)
		} else {
			logger.Info().Str("name", name).Str("namespace", namespace).Msgf("%s created", resourceType)
		}
	case *v1.ConfigMap:
		if _, err := clientset.CoreV1().ConfigMaps(namespace).Create(ctx, o, metav1.CreateOptions{}); err != nil {
			if !errors.IsAlreadyExists(err) {
				return fmt.Errorf("failed to create %s %s: %w", resourceType, name, err)
			}
			// Already exists, try to update
			if _, err := clientset.CoreV1().ConfigMaps(namespace).Update(ctx, o, metav1.UpdateOptions{}); err != nil {
				return fmt.Errorf("failed to update %s %s: %w", resourceType, name, err)
			}
			logger.Info().Str("name", name).Str("namespace", namespace).Msgf("%s updated", resourceType)
		} else {
			logger.Info().Str("name", name).Str("namespace", namespace).Msgf("%s created", resourceType)
		}
	default:
		return fmt.Errorf("unsupported resource type: %T", obj)
	}

	return nil
}

// CreateSecretsOverrideInput contains input for creating a secrets override K8s Secret
type CreateSecretsOverrideInput struct {
	Namespace    string            // K8s namespace
	InstanceName string            // Node instance name (e.g., "workflow-3")
	SecretsToml  string            // TOML content with secrets
	Labels       map[string]string // Optional labels for the secret
}

// CreateSecretsOverride creates a Kubernetes Secret for node secrets override
// Secret name format: <instance-name>-overrides-v2 (e.g., workflow-3-overrides-v2)
// This is separate from the main secret (<instance-name>-v2) which has .api and secret.toml
// Secret data key: "99-secrets-override.toml"
func CreateSecretsOverride(ctx context.Context, logger zerolog.Logger, input CreateSecretsOverrideInput) error {
	clientset, err := infra.GetKubernetesClient()
	if err != nil {
		return err
	}

	secretName := input.InstanceName + "-overrides-v2"
	secret := &v1.Secret{
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{
			"99-secrets-override.toml": []byte(input.SecretsToml),
		},
	}
	if input.Labels != nil {
		secret.Labels = input.Labels
	}

	return createOrUpdateResource(ctx, clientset, input.Namespace, secretName, secret, logger, "Secret")
}

// CreateConfigOverrideInput contains input for creating a config override K8s ConfigMap
type CreateConfigOverrideInput struct {
	Namespace    string // K8s namespace
	InstanceName string // Node instance name (e.g., "workflow-2")
	ConfigToml   string // TOML content with config overrides
}

// CreateConfigOverride creates a Kubernetes ConfigMap for node config override
// ConfigMap name format: <instance-name>-overrides-cm-v2 (e.g., workflow-2-overrides-cm-v2)
// ConfigMap data key: "99-overrides.toml"
func CreateConfigOverride(ctx context.Context, logger zerolog.Logger, input CreateConfigOverrideInput) error {
	clientset, err := infra.GetKubernetesClient()
	if err != nil {
		return err
	}

	configMapName := input.InstanceName + "-overrides-cm-v2"
	configMap := &v1.ConfigMap{
		Data: map[string]string{
			"99-overrides.toml": input.ConfigToml,
		},
	}

	return createOrUpdateResource(ctx, clientset, input.Namespace, configMapName, configMap, logger, "ConfigMap")
}
