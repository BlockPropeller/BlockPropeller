package test

import (
	"testing"
)

// Integration marks the test as an integration test.
//
// In case `go test` command is executed with a `-short` flag, integration test will be skipped.
func Integration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
}
