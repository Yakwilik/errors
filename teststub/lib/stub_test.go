package lib

import (
	"testing"
)

func TestWithDependency(t *testing.T) {
	dependency := New()

	err := dependency.DoWork(true)

	t.Logf("error: %v", err)
}
