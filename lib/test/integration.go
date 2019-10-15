package test

import (
	"testing"

	"blockpropeller.dev/lib/log"
)

// Integration marks the test as an integration test.
//
// In case `go test` command is executed with a `-short` flag, integration test will be skipped.
func Integration(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	log.SetGlobal(log.NewTestingLogger(t))
}
