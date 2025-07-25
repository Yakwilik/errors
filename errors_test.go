package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	err := New("something went wrong")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "[")
	assert.Contains(t, err.Error(), "something went wrong")

	t.Logf("error: %v", err)
}

func TestWrap(t *testing.T) {
	base := errors.New("root cause")
	wrapped := Wrap(base, "context added")

	assert.True(t, errors.Is(wrapped, base))

	t.Logf("error: %v", wrapped)
}
