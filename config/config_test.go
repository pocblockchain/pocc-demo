package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetConfig(t *testing.T) {
	path := "config.yml"
	conf, err := GetConfig(path)
	require.Nil(t, err)
	require.Equal(t, "18.182.40.136:26657", conf.RPCURL)
	require.Equal(t, "poc-mainnet", conf.ChainID)
	require.EqualValues(t, 200000, conf.Gas)
	require.Equal(t, "10000000000poc", conf.MinGasPrice)
}
