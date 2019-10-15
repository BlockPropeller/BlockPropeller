package binance_test

import (
	"testing"

	"blockpropeller.dev/blockpropeller/binance"
	"blockpropeller.dev/lib/test"
)

func TestNetworkIsValid(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"testnet", true},
		{"prod", true},
		{"", false},
		{"superenv", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			net := binance.NewNetwork(tt.name)

			got := net.IsValid()
			test.AssertBoolEqual(t, "Network.IsValid()", got, tt.valid)
			test.AssertStringsEqual(t, "Network.String()", net.String(), tt.name)
		})
	}
}
