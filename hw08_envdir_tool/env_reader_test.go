package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadEnvValue(t *testing.T) {
	testDir := "./testdata/env"
	actual := EnvValue{}

	t.Run("BAR", func(t *testing.T) {
		envs, err := ReadDir(testDir)
		for name, env := range envs {
			if name == "BAR" {
				actual = env
				break
			}
		}

		expected := EnvValue{
			Value:      "bar",
			NeedRemove: false,
		}

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("EMPTY", func(t *testing.T) {
		envs, err := ReadDir(testDir)
		for name, env := range envs {
			if name == "EMPTY" {
				actual = env
				break
			}
		}

		expected := EnvValue{
			Value:      "",
			NeedRemove: false,
		}

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("FOO", func(t *testing.T) {
		envs, err := ReadDir(testDir)
		for name, env := range envs {
			if name == "FOO" {
				actual = env
				break
			}
		}

		expected := EnvValue{
			Value:      "   foo\nwith new line",
			NeedRemove: false,
		}

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("HELLO", func(t *testing.T) {
		envs, err := ReadDir(testDir)
		for name, env := range envs {
			if name == "HELLO" {
				actual = env
				break
			}
		}

		expected := EnvValue{
			Value:      "\"hello\"",
			NeedRemove: false,
		}

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
	t.Run("UNSET", func(t *testing.T) {
		envs, err := ReadDir(testDir)
		for name, env := range envs {
			if name == "UNSET" {
				actual = env
				break
			}
		}

		expected := EnvValue{
			Value:      "",
			NeedRemove: true,
		}

		require.Equal(t, expected, actual)
		require.NoError(t, err)
	})
}
