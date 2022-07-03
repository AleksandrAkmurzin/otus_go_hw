package main

import (
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	environment, err := ReadDir("testdata/env")
	require.NoError(t, err)

	t.Run("OK", func(t *testing.T) {
		returnCode := RunCmd([]string{"testdata/echo.sh", "arg"}, environment)
		require.Equal(t, 0, returnCode)
	})

	t.Run("ERR", func(t *testing.T) {
		returnCode := RunCmd([]string{"testdata/error.sh"}, environment)
		require.Equal(t, 1, returnCode)
	})
}

func TestProcessEnv(t *testing.T) {
	testEnvName := "testEnvName"
	testEnvValue := "testEnvValue"

	t.Run("Add", func(t *testing.T) {
		_, exist := os.LookupEnv(testEnvName)
		require.False(t, exist)

		err := processEnv(testEnvName, EnvValue{testEnvValue, false})
		require.NoError(t, err)

		addedValue, ok := os.LookupEnv(testEnvName)
		require.True(t, ok)
		require.Equal(t, testEnvValue, addedValue)
	})

	t.Run("Replace", func(t *testing.T) {
		_, exist := os.LookupEnv(testEnvName)
		if !exist {
			err := os.Setenv(testEnvName, testEnvValue)
			require.NoError(t, err)
		}

		err := processEnv(testEnvName, EnvValue{testEnvValue, false})
		require.NoError(t, err)

		replacedValue, ok := os.LookupEnv(testEnvName)
		require.True(t, ok)
		require.Equal(t, testEnvValue, replacedValue)
	})

	t.Run("Remove existed", func(t *testing.T) {
		_, exist := os.LookupEnv(testEnvName)
		if !exist {
			err := os.Setenv(testEnvName, testEnvValue)
			require.NoError(t, err)
		}

		err := processEnv(testEnvName, EnvValue{testEnvValue, true})
		require.NoError(t, err)

		removedValue, ok := os.LookupEnv(testEnvName)
		require.False(t, ok)
		require.Empty(t, removedValue)
	})

	t.Run("Remove not existed", func(t *testing.T) {
		_, exist := os.LookupEnv(testEnvName)
		if exist {
			err := os.Unsetenv(testEnvName)
			require.NoError(t, err)
		}

		err := processEnv(testEnvName, EnvValue{testEnvValue, true})
		require.NoError(t, err)

		removedValue, ok := os.LookupEnv(testEnvName)
		require.False(t, ok)
		require.Empty(t, removedValue)
	})
}
