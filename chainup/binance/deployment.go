package binance

import (
	"chainup.dev/chainup/infrastructure"
	"github.com/blang/semver"
)

// NewNodeDeployment returns a new Binance Chain configuration in to form of a Deployment instance.
func NewNodeDeployment(network Network, nodeType NodeType, version semver.Version) *infrastructure.Deployment {
	return infrastructure.NewDeployment(
		infrastructure.DeploymentTypeBinanceNode,
		map[string]string{
			"binance_node_network": network.String(),
			"binance_node_type":    nodeType.String(),
			"binance_node_version": version.String(),
		},
	)
}
