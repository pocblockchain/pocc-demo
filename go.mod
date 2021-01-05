module github.com/pocblockchain/pocc-demo

go 1.13

require (
	github.com/pocblockchain/pocc v0.0.0
	github.com/btcsuite/btcd v0.0.0-20190115013929-ed77733ec07d
	github.com/cosmos/go-bip39 v0.0.0-20180618194314-52158e4697b8
	github.com/pkg/errors v0.8.1
	github.com/pocblockchain/pocc v0.0.0-20200812024435-686cbedc327d
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.4.0
	github.com/tendermint/crypto v0.0.0-20180820045704-3764759f34a5
	github.com/tendermint/tendermint v0.32.6
	gopkg.in/yaml.v2 v2.2.4
)

replace github.com/pocblockchain/pocc => ../pocc
