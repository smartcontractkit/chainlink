package infra

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"
)

// K8sClient wraps the Kubernetes client for discovering node runtime info.
type K8sClient struct {
	client     *kubernetes.Clientset
	dynamic    dynamic.Interface
	restConfig *rest.Config
	namespace  string
	log        zerolog.Logger
}

// NewK8sClient creates a Kubernetes client. If kubeconfigPath is empty,
// it falls back to in-cluster config, then to KUBECONFIG env or ~/.kube/config.
func NewK8sClient(kubeconfigPath, namespace string, log zerolog.Logger) (*K8sClient, error) {
	config, err := buildK8sConfig(kubeconfigPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build k8s config")
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s clientset")
	}

	dyn, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dynamic client")
	}

	return &K8sClient{
		client:     client,
		dynamic:    dyn,
		restConfig: config,
		namespace:  namespace,
		log:        log,
	}, nil
}

func buildK8sConfig(kubeconfigPath string) (*rest.Config, error) {
	if kubeconfigPath != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	// Try in-cluster first
	if config, err := rest.InClusterConfig(); err == nil {
		return config, nil
	}

	// Fall back to default kubeconfig
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
}

// NodeAPIInfo holds a node's API URL and credentials discovered from K8s.
type NodeAPIInfo struct {
	URL      string
	Email    string
	Password string
}

// GetNodeAPIInfo discovers a node's API URL and credentials from K8s resources.
// The API URL comes from an HTTPRoute, and the credentials come from a
// CNKO-managed Secret named <node>-secret with key ".api".
func (k *K8sClient) GetNodeAPIInfo(ctx context.Context, nodeName, namespace string) (*NodeAPIInfo, error) {
	if namespace == "" {
		namespace = k.namespace
	}
	apiURL, err := k.getNodeAPIURL(ctx, nodeName, namespace)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get API URL for node %s", nodeName)
	}

	email, password, err := k.getNodeAPICreds(ctx, nodeName, namespace)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get API creds for node %s", nodeName)
	}

	return &NodeAPIInfo{
		URL:      apiURL,
		Email:    email,
		Password: password,
	}, nil
}

func (k *K8sClient) getNodeAPIURL(ctx context.Context, nodeName, namespace string) (string, error) {
	// HTTPRoute is a Gateway API CRD (gateway-networking.k8s.io/v1)
	gvr := schema.GroupVersionResource{
		Group:    "gateway.networking.k8s.io",
		Version:  "v1",
		Resource: "httproutes",
	}

	routes, err := k.dynamic.Resource(gvr).Namespace(namespace).List(ctx, v1.ListOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "failed to list HTTPRoutes in namespace %s", namespace)
	}

	for _, route := range routes.Items {
		name := route.GetName()
		// Match by name prefix — HTTPRoutes are typically named <node-name> or similar
		if strings.HasPrefix(name, nodeName) || strings.Contains(name, nodeName) {
			hostnames, _, _ := unstructuredNestedStringSlice(route.Object, "spec", "hostnames")
			if len(hostnames) > 0 {
				return "https://" + hostnames[0], nil
			}
		}
	}

	return "", fmt.Errorf("no HTTPRoute found matching node %s in namespace %s", nodeName, namespace)
}

func (k *K8sClient) getNodeAPICreds(ctx context.Context, nodeName, namespace string) (string, string, error) {
	secretName := nodeName + "-secret"
	secret, err := k.client.CoreV1().Secrets(namespace).Get(ctx, secretName, v1.GetOptions{})
	if err != nil {
		return "", "", errors.Wrapf(err, "failed to get secret %s/%s", namespace, secretName)
	}

	apiData, ok := secret.Data[".api"]
	if !ok {
		return "", "", fmt.Errorf("secret %s/%s has no .api key", namespace, secretName)
	}

	decoded, err := base64.StdEncoding.DecodeString(string(apiData))
	if err != nil {
		// Might already be decoded (raw bytes)
		decoded = apiData
	}

	lines := strings.SplitN(string(decoded), "\n", 2)
	if len(lines) < 2 {
		return "", "", fmt.Errorf("unexpected .api format in secret %s/%s: expected 2 lines (email, password)", namespace, secretName)
	}

	return strings.TrimSpace(lines[0]), strings.TrimSpace(lines[1]), nil
}

// GetNodeSecretsToml reads the 00-secrets.toml key from a node's K8s secret.
// This contains the P2P encrypted key JSON needed for NewDonMetadata ImportedSecrets.
func (k *K8sClient) GetNodeSecretsToml(ctx context.Context, nodeName, namespace string) (string, error) {
	if namespace == "" {
		namespace = k.namespace
	}
	secretName := nodeName + "-secret"
	secret, err := k.client.CoreV1().Secrets(namespace).Get(ctx, secretName, v1.GetOptions{})
	if err != nil {
		return "", errors.Wrapf(err, "failed to get secret %s/%s", namespace, secretName)
	}
	data, ok := secret.Data["00-secrets.toml"]
	if !ok {
		return "", fmt.Errorf("secret %s/%s has no 00-secrets.toml key", namespace, secretName)
	}
	return string(data), nil
}

// GetNodePeerID reads a node's P2P PeerID from the node's admin API.
// This is preferred over reading from K8s secrets because it works regardless
// of how the P2P key was delivered to the node.
func (k *K8sClient) GetNodePeerID(ctx context.Context, nodeName, namespace string) (string, error) {
	apiInfo, err := k.GetNodeAPIInfo(ctx, nodeName, namespace)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get API info for node %s", nodeName)
	}

	client, err := clclient.NewChainlinkClient(&clclient.Config{
		URL:      apiInfo.URL,
		Email:    apiInfo.Email,
		Password: apiInfo.Password,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to connect to node %s API at %s", nodeName, apiInfo.URL)
	}

	p2pKeys, err := client.MustReadP2PKeys()
	if err != nil {
		return "", errors.Wrapf(err, "failed to read P2P keys from node %s", nodeName)
	}

	if len(p2pKeys.Data) == 0 {
		return "", fmt.Errorf("no P2P keys found on node %s", nodeName)
	}

	// MustReadP2PKeys already strips the "p2p_" prefix from PeerID
	return p2pKeys.Data[0].Attributes.PeerID, nil
}

// GetNodeCSAKey reads a node's CSA key by authenticating to the node's admin API
// using the chainlink-testing-framework clclient. The node API URL and credentials
// are discovered from K8s (HTTPRoute + CNKO secret).
func (k *K8sClient) GetNodeCSAKey(ctx context.Context, nodeName, namespace string) (string, error) {
	apiInfo, err := k.GetNodeAPIInfo(ctx, nodeName, namespace)
	if err != nil {
		return "", errors.Wrapf(err, "failed to get API info for node %s", nodeName)
	}

	client, err := clclient.NewChainlinkClient(&clclient.Config{
		URL:      apiInfo.URL,
		Email:    apiInfo.Email,
		Password: apiInfo.Password,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to connect to node %s API at %s", nodeName, apiInfo.URL)
	}

	csaKeys, _, err := client.ReadCSAKeys()
	if err != nil {
		return "", errors.Wrapf(err, "failed to read CSA keys from node %s", nodeName)
	}

	if len(csaKeys.Data) == 0 {
		return "", fmt.Errorf("no CSA keys found on node %s", nodeName)
	}

	key := csaKeys.Data[0].Attributes.PublicKey
	// The node API returns the CSA key with a "csa_" prefix (e.g. "csa_2471a9a3...").
	// JD expects the raw hex key without the prefix.
	key = strings.TrimPrefix(key, "csa_")

	return key, nil
}

// WaitForRollout polls until the node's pod has restarted with new config.
// It checks the pod's restart count or creation timestamp to detect re-roll.
func (k *K8sClient) WaitForRollout(ctx context.Context, nodeName, namespace string, timeout time.Duration) error {
	if namespace == "" {
		namespace = k.namespace
	}
	// Use kubectl rollout status equivalent — check deployment or pod status
	// CNKO creates Deployments named <node-name> or similar
	labelSelector := "app.kubernetes.io/instance=" + nodeName

	pods, err := k.client.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to list pods for node %s", nodeName)
	}

	if len(pods.Items) == 0 {
		// Try a broader selector
		pods, err = k.client.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{
			LabelSelector: "app=" + nodeName,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to list pods for node %s (broad selector)", nodeName)
		}
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for node %s in namespace %s", nodeName, namespace)
	}

	deadline := time.Now().Add(timeout)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for rollout of node %s", nodeName)
		}

		// Check if pods are ready
		pods, err = k.client.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{
			LabelSelector: labelSelector,
		})
		if err != nil {
			k.log.Warn().Err(err).Msgf("failed to list pods for %s, retrying", nodeName)
		} else {
			allReady := true
			for _, pod := range pods.Items {
				for _, cond := range pod.Status.Conditions {
					if cond.Type == "Ready" && cond.Status != "True" {
						allReady = false
						break
					}
				}
			}
			if allReady && len(pods.Items) > 0 {
				return nil
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

// RestartNodePods deletes the pod(s) backing nodeName so its owning
// controller recreates them fresh. Used as an ad-hoc workaround for
// capabilities that don't clean up state correctly after job cancellation.
func (k *K8sClient) RestartNodePods(ctx context.Context, nodeName, namespace string) error {
	if namespace == "" {
		namespace = k.namespace
	}

	labelSelector := "app.kubernetes.io/instance=" + nodeName
	pods, err := k.client.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return errors.Wrapf(err, "failed to list pods for node %s", nodeName)
	}

	if len(pods.Items) == 0 {
		pods, err = k.client.CoreV1().Pods(namespace).List(ctx, v1.ListOptions{
			LabelSelector: "app=" + nodeName,
		})
		if err != nil {
			return errors.Wrapf(err, "failed to list pods for node %s (broad selector)", nodeName)
		}
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("no pods found for node %s in namespace %s", nodeName, namespace)
	}

	for _, pod := range pods.Items {
		if err := k.client.CoreV1().Pods(namespace).Delete(ctx, pod.Name, v1.DeleteOptions{}); err != nil {
			return errors.Wrapf(err, "failed to delete pod %s for node %s", pod.Name, nodeName)
		}
	}

	return nil
}

// CopyFilesToPod streams the given local files into destDir inside the named
// container of a running pod, the same way `kubectl cp` does internally: it
// execs a shell in the container that creates destDir (tar itself won't
// create a missing -C target) and then runs `tar xf -` into it, piping a tar
// archive of the files into its stdin over a SPDY exec stream. No pod
// restart is triggered — the files just appear on the container's filesystem.
func (k *K8sClient) CopyFilesToPod(ctx context.Context, namespace, podName, container, destDir string, localPaths []string) error {
	if namespace == "" {
		namespace = k.namespace
	}
	if k.restConfig == nil {
		return errors.New("k8s client has no rest.Config — cannot exec into pods")
	}

	var tarBuf bytes.Buffer
	if err := writeFilesAsTar(&tarBuf, localPaths); err != nil {
		return errors.Wrap(err, "failed to build tar archive")
	}

	req := k.client.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")
	req.VersionedParams(&corev1.PodExecOptions{
		Container: container,
		// destDir is passed as $1 rather than interpolated into the script
		// string, so it's safe even if it contains spaces or shell metachars.
		Command: []string{"sh", "-c", `mkdir -p "$1" && exec tar xf - -C "$1"`, "sh", destDir},
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
	}, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(k.restConfig, "POST", req.URL())
	if err != nil {
		return errors.Wrapf(err, "failed to create exec stream for pod %s/%s", namespace, podName)
	}

	var stdout, stderr bytes.Buffer
	err = executor.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdin:  &tarBuf,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		return errors.Wrapf(err, "failed to stream files into pod %s/%s: %s", namespace, podName, stderr.String())
	}

	return nil
}

// writeFilesAsTar writes each local file as a flat entry (basename only, no
// directory structure) into a tar archive, so the remote `tar xf -C destDir`
// extracts them directly into destDir.
func writeFilesAsTar(w io.Writer, localPaths []string) error {
	tw := tar.NewWriter(w)

	for _, path := range localPaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrapf(err, "failed to read local file %s", path)
		}

		hdr := &tar.Header{
			Name: filepath.Base(path),
			Mode: 0o644,
			Size: int64(len(data)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			return errors.Wrapf(err, "failed to write tar header for %s", path)
		}
		if _, err := tw.Write(data); err != nil {
			return errors.Wrapf(err, "failed to write tar contents for %s", path)
		}
	}

	return tw.Close()
}

// unstructuredNestedStringSlice navigates a nested map[string]interface{} (from
// unstructured.Unstructured.Object) to find a string slice at the given path.
func unstructuredNestedStringSlice(obj map[string]any, keys ...string) ([]string, bool, error) {
	current := any(obj)
	for _, key := range keys {
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false, fmt.Errorf("expected map at key %s", key)
		}
		current, ok = m[key]
		if !ok {
			return nil, false, nil
		}
	}
	slice, ok := current.([]any)
	if !ok {
		return nil, false, fmt.Errorf("expected slice at path %v", keys)
	}
	result := make([]string, 0, len(slice))
	for _, v := range slice {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}
	return result, true, nil
}
