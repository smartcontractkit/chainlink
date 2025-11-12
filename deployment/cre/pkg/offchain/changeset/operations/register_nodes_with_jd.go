package operations

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/semver/v3"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink-deployments-framework/operations"
	nodev1 "github.com/smartcontractkit/chainlink-protos/job-distributor/v1/node"
	"github.com/smartcontractkit/chainlink-protos/job-distributor/v1/shared/ptypes"

	"github.com/smartcontractkit/chainlink/deployment/cre/pkg/offchain"
)

var ErrNodeAlreadyExists = errors.New("node with given public key already exists")

type JDRegisterNodeOpDeps struct {
	Env cldf.Environment
}
type JDRegisterNodeOpInput struct {
	IsBootstrap bool              `json:"isBootstrap" yaml:"is_bootstrap"`
	Domain      string            `json:"domain" yaml:"domain"`
	Name        string            `json:"name" yaml:"name"`
	CSAKey      string            `json:"csaKey" yaml:"csaKey"`
	P2PID       string            `json:"p2pID" yaml:"p2pID"`
	DONName     string            `json:"donName" yaml:"don_name"`
	Zone        string            `json:"zone" yaml:"zone"`
	Labels      map[string]string `json:"labels" yaml:"labels"` // key-value pairs; use empty string for value-less labels
}

type JDRegisterNodeOpOutput struct {
	Error string `json:"error,omitempty" yaml:"error,omitempty"` // empty if no error
	Node  JDNode `json:"node,omitempty" yaml:"node,omitempty"`   // non-nil if success
}

// JDRegisterNodeOp registers a node with Job Distributor.
// It should only be used if the node does not already exist. If the node already exists, it returns an error.
// If you are unsure if the node exists, use JDUpsertNodeOp instead.
var JDRegisterNodeOp = operations.NewOperation(
	"jd-register-node",
	semver.MustParse("1.0.0"),
	"Registers a node with Job Distributor",
	func(e operations.Bundle, deps JDRegisterNodeOpDeps, input JDRegisterNodeOpInput) (JDRegisterNodeOpOutput, error) {
		deps.Env.Logger.Infof("Registering node `%s` with JD", input.Name)
		node, err := registerNodeImpl(deps, input)
		if err != nil {
			return JDRegisterNodeOpOutput{
				Node:  NewJDNodeFromProto(node),
				Error: err.Error(),
			}, err
		}

		return JDRegisterNodeOpOutput{
			Node: NewJDNodeFromProto(node),
		}, nil
	},
)

// JDUpsertNodeOp upserts a node with Job Distributor, ensuring it has the correct DON label regardless of whether it already existed or not.
// If the node does not exist, it registers it. If it already exists, it ensures the DON label is present.
var JDUpsertNodeOp = operations.NewOperation(
	"jd-upsert-node",
	semver.MustParse("1.0.0"),
	"Upserts a node with Job Distributor, ensuring it has the correct DON label regardless of whether it already existed or not.",
	func(b operations.Bundle, deps JDRegisterNodeOpDeps, input JDRegisterNodeOpInput) (JDRegisterNodeOpOutput, error) {
		deps.Env.Logger.Infof("Upserting node `%s` with JD", input.Name)

		node, err := registerNodeImpl(deps, input)
		if err == nil {
			// the don label gets added if the node didn't exist before, so return early
			return JDRegisterNodeOpOutput{
				Node: NewJDNodeFromProto(node),
			}, nil
		}
		if err != nil && !strings.Contains(err.Error(), ErrNodeAlreadyExists.Error()) {
			return JDRegisterNodeOpOutput{
				Error: err.Error(),
				Node:  NewJDNodeFromProto(node),
			}, err
		}

		// node already exists, ensure it has the don label

		var output JDRegisterNodeOpOutput
		nodeInfo, nerr := ensureDONLabelOnNode(deps.Env, node, input.DONName)
		if nerr != nil {
			output.Error = nerr.Error()
			return output, nerr
		}
		deps.Env.Logger.Infof("Ensured node's `%s` DON label", input.Name)
		output.Node = NewJDNodeFromProto(nodeInfo)
		return output, nil
	},
)

func registerNodeImpl(deps JDRegisterNodeOpDeps, input JDRegisterNodeOpInput) (*nodev1.Node, error) {
	if input.CSAKey == "" {
		return nil, errors.New("CSAKey is required")
	}

	n, _ := deps.Env.Offchain.GetNode(deps.Env.GetContext(), &nodev1.GetNodeRequest{
		PublicKey: &input.CSAKey,
	})
	if n != nil {
		// node already exists, nothing to do
		return n.GetNode(), ErrNodeAlreadyExists
	}

	donLabel := "don-" + input.DONName
	zoneLbl := input.Zone
	labels := map[string]string{
		donLabel: "",
	}
	if zoneLbl != "" {
		labels["zone"] = zoneLbl
	}
	if input.P2PID != "" {
		labels["p2p_id"] = input.P2PID
	}
	for k, v := range input.Labels {
		labels[k] = v
	}
	nodeID, err := offchain.RegisterNode(
		deps.Env.GetContext(),
		deps.Env.Offchain,
		input.Name,
		input.CSAKey,
		input.IsBootstrap,
		input.Domain,
		deps.Env.Name,
		labels,
	)
	if err != nil {
		// We don't want to fail the entire migration if one node fails to register, so we just log the error.
		terr := fmt.Errorf("failed to register node %s for don %s: %w", input.Name, input.DONName, err)
		deps.Env.Logger.Errorw("failed to register node", "don", input.DONName, "node", input.Name, "error", err)
		return nil, terr
	}
	// must get the node again to return the full node info such as the IsEnabled field, IsConnected, etc.
	n2, err := deps.Env.Offchain.GetNode(deps.Env.GetContext(), &nodev1.GetNodeRequest{
		Id: nodeID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get node %s after registering it for don %s: %w", input.Name, input.DONName, err)
	}
	deps.Env.Logger.Infow("registered node", "name", input.Name, "id", nodeID)
	return n2.GetNode(), nil
}

func ensureDONLabelOnNode(env cldf.Environment, nodeInfo *nodev1.Node, donName string) (*nodev1.Node, error) {
	donLabel := "don-" + donName

	// ensure the node has a label for the current don name
	// if not, add the new don label
	hasDonLabel := false
	for _, l := range nodeInfo.Labels {
		if l.Key == donLabel {
			hasDonLabel = true
			break
		}
	}
	if hasDonLabel {
		env.Logger.Infow("don label already set", "name", donName)
		// nothing to do, return the existing node info
		return nodeInfo, nil
	}

	// add the don label
	var labels []*ptypes.Label
	newLabels := nodeInfo.GetLabels()
	for _, l := range newLabels {
		if l.Key == donName {
			continue
		}
		labels = append(labels, l)
	}

	value := ""
	labels = append(labels, &ptypes.Label{
		Key:   donLabel,
		Value: &value,
	})

	_, uerr := env.Offchain.UpdateNode(env.GetContext(), &nodev1.UpdateNodeRequest{
		Id:        nodeInfo.GetId(),
		Name:      nodeInfo.GetName(),
		PublicKey: nodeInfo.GetPublicKey(),
		Labels:    labels,
	})
	if uerr != nil {
		rerr := fmt.Errorf("failed to update DON for node %s for don %s: %w", nodeInfo.Name, donName, uerr)
		env.Logger.Errorw("failed to update DON for node", "don", donName, "name", nodeInfo.Name, "err", uerr)
		return nil, rerr
	}

	env.Logger.Infof("updated node %s to include label [key = %s, value = \"\"]", nodeInfo.Name, donLabel)

	return &nodev1.Node{
		Id:            nodeInfo.GetId(),
		Name:          nodeInfo.GetName(),
		PublicKey:     nodeInfo.GetPublicKey(),
		Labels:        labels,
		IsEnabled:     nodeInfo.GetIsEnabled(),
		IsConnected:   nodeInfo.GetIsConnected(),
		P2PKeyBundles: nodeInfo.P2PKeyBundles,
	}, nil
}
