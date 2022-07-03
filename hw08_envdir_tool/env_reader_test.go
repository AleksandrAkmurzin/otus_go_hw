package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("Wrong dir name", func(t *testing.T) {
		_, err := ReadDir("testdata")
		if err != nil {
			fmt.Println(err)
		}
		require.Error(t, err)
	})

	t.Run("Wrong file name", func(t *testing.T) {
		info, err := os.Stat("testdata/wrong=name")
		require.NoError(t, err)

		_, err = fileToEnvValue("testdata", info)
		require.ErrorIs(t, err, ErrUnsupportedFileName)
	})

	t.Run("Positive case", func(t *testing.T) {
		actualEnv, err := ReadDir("testdata/env")
		require.NoError(t, err)

		expectedEnv := Environment{
			"BAR":   EnvValue{Value: "bar"},
			"EMPTY": EnvValue{},
			"FOO":   EnvValue{Value: "   foo\nwith new line"},
			"HELLO": EnvValue{Value: "\"hello\""},
			"UNSET": EnvValue{NeedRemove: true},
		}

		for name, value := range expectedEnv {
			require.Equal(t, value, actualEnv[name], "Value for %s env param is wrong", name)
		}
	})
}
