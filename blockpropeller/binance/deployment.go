package binance

import (
	"fmt"

	"blockpropeller.dev/blockpropeller/infrastructure"
	"github.com/blang/semver"
	"github.com/pkg/errors"
)

var (
	configNetwork  = "binance_node_network"
	configNodeType = "binance_node_type"
	configVersion  = "binance_node_version"
)

func init() {
	infrastructure.RegisterDeploymentType(infrastructure.DeploymentTypeBinanceNode, nodeDeploymentSpec{})
}

// nodeDeploymentSpec exposes methods for interacting with the Binance Node Deployment.
type nodeDeploymentSpec struct {
}

// UnmarshalConfig converts a raw map of strings into a binance.NodeConfig struct.
func (nodeDeploymentSpec) UnmarshalConfig(raw map[string]string) (infrastructure.DeploymentConfig, error) {
	network := NewNetwork(raw[configNetwork])
	if !network.IsValid() {
		return nil, errors.New("invalid binance node network")
	}

	nodeType := NewNodeType(raw[configNodeType])
	if !network.IsValid() {
		return nil, errors.New("invalid binance node type")
	}

	version, err := semver.Parse(raw[configVersion])
	if err != nil {
		return nil, errors.Wrap(err, "invalid binance node version")
	}

	return &NodeConfig{
		Network:  network,
		NodeType: nodeType,
		Version:  version,
	}, nil
}

// HealthCheck returns a HealthCheck to be used to determine Deployment health.
func (nodeDeploymentSpec) HealthCheck(srv *infrastructure.Server, deployment *infrastructure.Deployment) (infrastructure.HealthCheck, error) {
	url := fmt.Sprintf("http://%s:27147/status", srv.IPAddress.String())

	return infrastructure.NewHTTPHealthCheck("GET", url, 200), nil
}

// NodeConfig holds the configuration for the Binance Chain Node.
type NodeConfig struct {
	Network  Network
	NodeType NodeType
	Version  semver.Version
}

// MarshalMap converts a NodeConfig to a map[string]string.
func (cfg *NodeConfig) MarshalMap() map[string]string {
	return map[string]string{
		configNetwork:  cfg.Network.String(),
		configNodeType: cfg.NodeType.String(),
		configVersion:  cfg.Version.String(),
	}
}

// NewNodeDeployment returns a new Binance Chain configuration in to form of a Deployment instance.
func NewNodeDeployment(network Network, nodeType NodeType, version semver.Version) *infrastructure.Deployment {
	return infrastructure.NewDeployment(
		infrastructure.DeploymentTypeBinanceNode,
		&NodeConfig{
			Network:  network,
			NodeType: nodeType,
			Version:  version,
		},
	)
}
