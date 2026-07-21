package domain

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// ChartNodeInfo holds information about a node parsed from Griddle chart values.
type ChartNodeInfo struct {
	Name       string
	NodeType   NodeRole // standard, boot, gateway (default: standard)
	Namespace  string   // K8s namespace for this node's node-set instance
	ConfigFile string   // absolute path to the env-specific YAML to patch
	DONName    string   // chainlink-node.registerNodes.labels.don-name for this node's node-set instance
}

// ChartValues holds the parsed Griddle chart values for a nodeset.
type ChartValues struct {
	Nodes     []ChartNodeInfo
	Namespace string
}

// LoadChartValues reads and parses the chart values for a Griddle nodeset.
//
// repoRoot can be:
//   - the repo root containing griddle.yaml (preferred — we parse it to find
//     the service name, namespace, and config file paths)
//   - a directory that directly contains shared.yaml + <env>.yaml (legacy
//     fallback for testing)
//
// env is the environment name (e.g. "dev", "stage", "prod") matching the
// key in griddle.yaml's deploy section.
func LoadChartValues(repoRoot, env string) (*ChartValues, error) {
	// Try to parse griddle.yaml first
	griddlePath := filepath.Join(repoRoot, "griddle.yaml")
	if _, err := os.Stat(griddlePath); err == nil {
		return loadFromGriddleYAML(repoRoot, griddlePath, env)
	}

	// Fallback: treat repoRoot as a direct config dir
	return loadFromConfigDir(repoRoot, env)
}

// loadFromGriddleYAML parses griddle.yaml to find the node-set instance,
// its namespace, and config file paths for the given environment.
func loadFromGriddleYAML(repoRoot, griddlePath, env string) (*ChartValues, error) {
	data, err := os.ReadFile(griddlePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read griddle.yaml %s", griddlePath)
	}

	var gy griddleYAML
	if err := yaml.Unmarshal(data, &gy); err != nil {
		return nil, errors.Wrapf(err, "failed to parse griddle.yaml %s", griddlePath)
	}

	// Find the environment
	envDeploy, ok := gy.Deploy[env]
	if !ok {
		// If only one environment exists, use it regardless of the env param
		if len(gy.Deploy) != 1 {
			available := make([]string, 0, len(gy.Deploy))
			for k := range gy.Deploy {
				available = append(available, k)
			}
			sort.Strings(available)
			return nil, fmt.Errorf("environment %q not found in griddle.yaml; available: %s", env, strings.Join(available, ", "))
		}
		for _, d := range gy.Deploy {
			envDeploy = d
			break
		}
	}

	// Collect all node-set instances (skip chaos and other non-node-set stacks).
	var allInstances []griddleInstance
	for i := range envDeploy.Instances {
		inst := &envDeploy.Instances[i]
		if strings.Contains(inst.Path, "node-set") {
			allInstances = append(allInstances, *inst)
		}
	}
	if len(allInstances) == 0 {
		return nil, fmt.Errorf("no node-set instance found in griddle.yaml for environment %q", env)
	}

	cv := &ChartValues{}

	for _, inst := range allInstances {
		envConfigFile, err := resolveInstanceEnvConfigFile(repoRoot, env, inst)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to resolve env config file for instance %s", inst.Name)
		}

		merged := make(map[string]any)
		for _, cfgPath := range inst.Config {
			fullPath := cfgPath
			if !filepath.IsAbs(fullPath) {
				fullPath = filepath.Join(repoRoot, cfgPath)
			}
			fileData, readErr := os.ReadFile(fullPath)
			if readErr != nil {
				return nil, errors.Wrapf(readErr, "failed to read config file %s", fullPath)
			}
			var fileVals map[string]any
			if unmarshalErr := yaml.Unmarshal(fileData, &fileVals); unmarshalErr != nil {
				return nil, errors.Wrapf(unmarshalErr, "failed to parse config file %s", fullPath)
			}
			merged = deepMerge(merged, fileVals)
		}

		instances := getNodeInstances(merged)
		var donName string
		if len(instances) > 0 {
			donName, err = donNameLabel(merged)
			if err != nil {
				return nil, errors.Wrapf(err, "node-set instance %q (%s)", inst.Name, envConfigFile)
			}
		}
		for name, nodeType := range instances {
			cv.Nodes = append(cv.Nodes, ChartNodeInfo{
				Name:       name,
				NodeType:   nodeType,
				Namespace:  inst.Namespace,
				ConfigFile: envConfigFile,
				DONName:    donName,
			})
		}

		// Primary namespace comes from the first node-set instance.
		if cv.Namespace == "" {
			cv.Namespace = inst.Namespace
		}
	}

	sort.Slice(cv.Nodes, func(i, j int) bool { return cv.Nodes[i].Name < cv.Nodes[j].Name })
	applyInferredRoles(cv)

	return cv, nil
}

// loadFromConfigDir loads chart values directly from a directory containing
// shared.yaml and <env>.yaml (used for testing).
func loadFromConfigDir(configDir, env string) (*ChartValues, error) {
	sharedPath := filepath.Join(configDir, "shared.yaml")
	envPath := filepath.Join(configDir, env+".yaml")

	sharedData, err := os.ReadFile(sharedPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read shared values %s", sharedPath)
	}

	var shared, envVals map[string]any
	if err = yaml.Unmarshal(sharedData, &shared); err != nil {
		return nil, errors.Wrapf(err, "failed to parse shared values %s", sharedPath)
	}

	if _, statErr := os.Stat(envPath); statErr == nil {
		envData, _ := os.ReadFile(envPath)
		if err = yaml.Unmarshal(envData, &envVals); err != nil {
			return nil, errors.Wrapf(err, "failed to parse %s values %s", env, envPath)
		}
	}

	merged := deepMerge(shared, envVals)

	absEnvPath, err := filepath.Abs(envPath)
	if err != nil {
		absEnvPath = envPath
	}

	cv := &ChartValues{}

	// Parse chainlink-node instances
	instances := getNodeInstances(merged)
	var donName string
	if len(instances) > 0 {
		donName, err = donNameLabel(merged)
		if err != nil {
			return nil, errors.Wrapf(err, "node-set config %q", absEnvPath)
		}
	}
	for name, nodeType := range instances {
		cv.Nodes = append(cv.Nodes, ChartNodeInfo{
			Name:       name,
			NodeType:   nodeType,
			ConfigFile: absEnvPath,
			DONName:    donName,
		})
	}
	sort.Slice(cv.Nodes, func(i, j int) bool { return cv.Nodes[i].Name < cv.Nodes[j].Name })
	applyInferredRoles(cv)

	return cv, nil
}

// griddleYAML is the structure of griddle.yaml (only the fields we need).
type griddleYAML struct {
	Deploy map[string]griddleEnvDeploy `yaml:"deploy"`
}

type griddleEnvDeploy struct {
	Instances []griddleInstance `yaml:"instances"`
}

type griddleInstance struct {
	Name      string   `yaml:"name"`
	Namespace string   `yaml:"namespace"`
	Path      string   `yaml:"path"`
	Version   string   `yaml:"version"`
	Config    []string `yaml:"config"`
}

// applyInferredRoles upgrades node roles that were inferred by naming convention.
// If a node is named *-bt-* or *bootstrap* but its nodeType is standard (the
// chart default), we upgrade it to boot so the UI and downstream logic see the
// correct role.
func applyInferredRoles(cv *ChartValues) {
	for i, n := range cv.Nodes {
		if n.NodeType != RoleStandard {
			continue
		}
		lower := strings.ToLower(n.Name)
		if strings.Contains(lower, "-bt-") || strings.Contains(lower, "bootstrap") {
			cv.Nodes[i].NodeType = RoleBootstrap
		}
	}
}

// FindBootstrap returns the bootstrap node name, or empty if none found.
func (cv *ChartValues) FindBootstrap() string {
	// First check for nodeType: boot
	for _, n := range cv.Nodes {
		if n.NodeType == RoleBootstrap {
			return n.Name
		}
	}
	// Fallback: naming convention *-bt-* or *bootstrap*
	for _, n := range cv.Nodes {
		lower := strings.ToLower(n.Name)
		if strings.Contains(lower, "-bt-") || strings.Contains(lower, "bootstrap") {
			return n.Name
		}
	}
	return ""
}

// FindGatewayNodes returns all nodes with nodeType: gateway.
func (cv *ChartValues) FindGatewayNodes() []ChartNodeInfo {
	var gws []ChartNodeInfo
	for _, n := range cv.Nodes {
		if n.NodeType == RoleGateway {
			gws = append(gws, n)
		}
	}
	return gws
}

// HasNode checks if a node with the given name exists in the chart values.
func (cv *ChartValues) HasNode(name string) bool {
	for _, n := range cv.Nodes {
		if n.Name == name {
			return true
		}
	}
	return false
}

// NodeIdentity returns a stable namespace-aware key for a node.
func NodeIdentity(namespace, name string) string {
	return namespace + "/" + name
}

// BuildNodeConfigFileMap returns nodeKey -> absolute config YAML path for all nodes.
func (cv *ChartValues) BuildNodeConfigFileMap() (map[string]string, error) {
	mapping := make(map[string]string, len(cv.Nodes))
	for _, n := range cv.Nodes {
		key := NodeIdentity(n.Namespace, n.Name)
		if n.ConfigFile == "" {
			return nil, fmt.Errorf("node %s: missing config file metadata", key)
		}
		if existing, ok := mapping[key]; ok && existing != n.ConfigFile {
			return nil, fmt.Errorf("node %s: conflicting config file mappings (%s vs %s)", key, existing, n.ConfigFile)
		}
		mapping[key] = n.ConfigFile
	}
	return mapping, nil
}

// GetNodeInNamespace returns node info for a namespace/name pair, or nil if not found.
func (cv *ChartValues) GetNodeInNamespace(namespace, name string) *ChartNodeInfo {
	key := NodeIdentity(namespace, name)
	for i := range cv.Nodes {
		if NodeIdentity(cv.Nodes[i].Namespace, cv.Nodes[i].Name) == key {
			return &cv.Nodes[i]
		}
	}
	return nil
}

// GetNode returns the node info for a given name, or nil if not found.
func (cv *ChartValues) GetNode(name string) *ChartNodeInfo {
	for i := range cv.Nodes {
		if cv.Nodes[i].Name == name {
			return &cv.Nodes[i]
		}
	}
	return nil
}

// GetNodeNamespace returns the K8s namespace for a node, or the primary namespace.
func (cv *ChartValues) GetNodeNamespace(nodeName string) string {
	for _, n := range cv.Nodes {
		if n.Name == nodeName {
			if n.Namespace != "" {
				return n.Namespace
			}
			break
		}
	}
	return cv.Namespace
}

// NodeInternalHost returns the real in-cluster DNS name for a chart-deployed
// node, using its actual (possibly per-node) namespace. This must be used
// anywhere a node's host is needed for real network reachability (e.g. OCR
// peering) — as opposed to system-tests/lib's own NodeMetadata.Host, which is
// a synthetic "<donName>-bt-<index>" convention built for CTFv2's own
// docker-compose/k8s provisioning and has no relation to Griddle's actual
// chart-deployed node names/namespaces.
func (cv *ChartValues) NodeInternalHost(nodeName string) string {
	return fmt.Sprintf("%s.%s.svc.cluster.local", nodeName, cv.GetNodeNamespace(nodeName))
}

// NodeNamesForDONName returns all chart node names whose node-set instance was
// registered with the given DON name, sorted.
func (cv *ChartValues) NodeNamesForDONName(donName string) []string {
	var names []string
	for _, n := range cv.Nodes {
		if n.DONName == donName {
			names = append(names, n.Name)
		}
	}
	sort.Strings(names)
	return names
}

// --- internal helpers ---

// donNameLabel reads the DON identity Griddle registered these nodes with in
// JD (chainlink-node.registerNodes.labels.don-name). It's required: without it
// there is no label for a JD `don_name` job-proposal filter to ever match.
func donNameLabel(values map[string]any) (string, error) {
	clNode := getMap(values, "chainlink-node")
	registerNodes := getMap(clNode, "registerNodes")
	labels := getMap(registerNodes, "labels")
	name, _ := labels["don-name"].(string)
	if name == "" {
		return "", errors.New("chainlink-node.registerNodes.labels.don-name is required")
	}
	return name, nil
}

func getNodeInstances(values map[string]any) map[string]NodeRole {
	result := make(map[string]NodeRole)

	clNode := getMap(values, "chainlink-node")
	if clNode == nil {
		return result
	}

	// Check defaults.nodeType
	defaults := getMap(clNode, "defaults")
	defaultNodeType := RoleStandard
	if defaults != nil {
		if nt, ok := defaults["nodeType"].(string); ok {
			defaultNodeType = parseNodeRole(nt)
		}
	}

	instances := getMap(clNode, "instances")
	if instances == nil {
		return result
	}

	for name, val := range instances {
		nodeType := defaultNodeType
		if instMap, ok := val.(map[string]any); ok {
			if nt, ok := instMap["nodeType"].(string); ok {
				nodeType = parseNodeRole(nt)
			}
		}
		result[name] = nodeType
	}

	return result
}

func resolveInstanceEnvConfigFile(repoRoot, env string, inst griddleInstance) (string, error) {
	suffix := "/" + env + ".yaml"
	for _, cfgPath := range inst.Config {
		normalized := strings.ReplaceAll(cfgPath, "\\", "/")
		if strings.HasSuffix(normalized, suffix) {
			fullPath := cfgPath
			if !filepath.IsAbs(fullPath) {
				fullPath = filepath.Join(repoRoot, cfgPath)
			}
			return filepath.Clean(fullPath), nil
		}
	}
	if len(inst.Config) == 0 {
		return "", fmt.Errorf("instance %s has no config files", inst.Name)
	}
	cfgPath := inst.Config[len(inst.Config)-1]
	fullPath := cfgPath
	if !filepath.IsAbs(fullPath) {
		fullPath = filepath.Join(repoRoot, cfgPath)
	}
	return filepath.Clean(fullPath), nil
}

func parseNodeRole(s string) NodeRole {
	switch strings.ToLower(s) {
	case "boot", "bootstrap":
		return RoleBootstrap
	case "gateway":
		return RoleGateway
	default:
		return RoleStandard
	}
}

func getMap(m map[string]any, key string) map[string]any {
	if m == nil {
		return nil
	}
	v, ok := m[key]
	if !ok {
		return nil
	}
	if mp, ok := v.(map[string]any); ok {
		return mp
	}
	// yaml.v3 decodes nested maps as map[string]interface{}
	if mp, ok := v.(map[string]any); ok {
		return mp
	}
	return nil
}

func deepMerge(dst, src map[string]any) map[string]any {
	result := make(map[string]any, len(dst))
	maps.Copy(result, dst)
	for k, v := range src {
		if existing, ok := result[k]; ok {
			if existingMap, ok1 := existing.(map[string]any); ok1 {
				if srcMap, ok2 := v.(map[string]any); ok2 {
					result[k] = deepMerge(existingMap, srcMap)
					continue
				}
			}
		}
		result[k] = v
	}
	return result
}
