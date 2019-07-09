package binance

var (
	// TypeLightNode is the light node configuration of a Binance Chain node
	TypeLightNode NodeType = "lightnode"
	// TypeFullNode is the full node that can participate as a validation in a Binance Chain network.
	TypeFullNode NodeType = "fullnode"

	// ValidNodeTypes that are recognized by ChainUP.
	ValidNodeTypes = []NodeType{TypeLightNode, TypeFullNode}
)

// NodeType of the node to be deployed.
type NodeType string

// NewNodeType returns a new NodeType instance.
func NewNodeType(node string) NodeType {
	return NodeType(node)
}

// IsValid checks whether the NodeType is one of recognized values.
func (t NodeType) IsValid() bool {
	for _, valid := range ValidNodeTypes {
		if t == valid {
			return true
		}
	}

	return false
}

// String satisfies the Stringer interface.
func (t NodeType) String() string {
	return string(t)
}
