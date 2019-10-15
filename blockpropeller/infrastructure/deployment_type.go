package infrastructure

var (
	// DeploymentTypeBinanceNode deploys a Binance Chain node to a target Server.
	DeploymentTypeBinanceNode = NewDeploymentType("binance_node")

	// ValidDeploymentTypes that are recognized by BlockPropeller.
	ValidDeploymentTypes = []DeploymentType{DeploymentTypeBinanceNode}
)

// DeploymentType is an identifier for the specific deployment being deployed.
type DeploymentType string

// NewDeploymentType returns a new DeploymentType instance.
func NewDeploymentType(t string) DeploymentType {
	return DeploymentType(t)
}

// IsValid checks whether the DeploymentType is one of recognized values.
func (t DeploymentType) IsValid() bool {
	for _, valid := range ValidDeploymentTypes {
		if t == valid {
			return true
		}
	}

	return false
}

// String satisfies the Stringer interface.
func (t DeploymentType) String() string {
	return string(t)
}
