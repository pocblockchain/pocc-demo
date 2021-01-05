package chainrpc

import (
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

func (c *Client) QueryBlock(height *int64) (*ctypes.ResultBlock, error) {
	return c.rpc.Block(height)
}
