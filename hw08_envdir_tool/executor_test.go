package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Go version", func(t *testing.T) {
		require.Equal(t, 0, RunCmd([]string{"go", "version"}, Environment{}))
	})
}
