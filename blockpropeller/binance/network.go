package binance

var (
	// NetworkTest is a Binance Chain testnet identifier.
	NetworkTest Network = "testnet"
	// NetworkProd is a Binance Chain production identifier.
	NetworkProd Network = "prod"

	// ValidNetworks that are recognized by BlockPropeller.
	ValidNetworks = []Network{NetworkTest, NetworkProd}
)

// Network that a Binance Chain node can join.
type Network string

// NewNetwork returns a new Network instance.
func NewNetwork(network string) Network {
	return Network(network)
}

// IsValid checks if the network is one of BlockPropeller recognized networks.
func (n Network) IsValid() bool {
	for _, valid := range ValidNetworks {
		if n == valid {
			return true
		}
	}

	return false
}

// String satisfies the stringer interface.
func (n Network) String() string {
	return string(n)
}
