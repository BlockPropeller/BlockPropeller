package binance_test

import (
	"testing"

	"chainup.dev/chainup/binance"
	"chainup.dev/lib/test"
)

func TestNodeTypeIsValid(t *testing.T) {
	tests := []struct {
		name  string
		valid bool
	}{
		{"", false},
		{"lightnode", true},
		{"fullnode", true},
		{"supernode", false},
		{"light", false},
		{"full", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			net := binance.NewNodeType(tt.name)

			got := net.IsValid()
			test.AssertBoolEqual(t, "NodeType.IsValid()", got, tt.valid)
			test.AssertStringsEqual(t, "NodeType.String()", net.String(), tt.name)
		})
	}
}
