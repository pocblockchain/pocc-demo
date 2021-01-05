package chainrpc

import (
	"encoding/hex"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (c *Client) QueryTxByHash(hash string) (*ctypes.ResultTx, error) {
	hashBytes, err := hex.DecodeString(hash)
	if err != nil {
		return nil, err
	}
	return c.rpc.Tx(hashBytes, false)

}
