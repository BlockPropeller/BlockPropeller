package test

import (
	"io"
	"testing"
)

// AssertStringsEqual assets that the returned string matches the expected one,
// reports error in case of mismatch and immediately fails the running test.
func AssertStringsEqual(t *testing.T, name string, got string, want string) {
	if got == want {
		return
	}

	t.Errorf("%s: got '%s', want '%s'", name, got, want)
	t.FailNow()
}

// AssertBoolEqual assets that the returned bool matches the expected one,
// reports error in case of mismatch and immediately fails the running test.
func AssertBoolEqual(t *testing.T, name string, got bool, want bool) {
	if got == want {
		return
	}

	t.Errorf("%s: got '%t', want '%t'", name, got, want)
	t.FailNow()
}

// CheckErr fails the currently running test if the provided error is not nil.
func CheckErr(t *testing.T, name string, err error) {
	if err == nil {
		return
	}

	t.Errorf("%s: got err '%s'", name, err)
	t.FailNow()
}

// Close reports any errors that occurred while closing the provided io.Closer
// to the testing framework.
//
// This is a helper function for one line defers to be used in test methods.
func Close(t *testing.T, closer io.Closer) {
	err := closer.Close()
	if err == nil {
		return
	}

	t.Errorf("failed closing %T: %s", closer, err)
	t.FailNow()
}
