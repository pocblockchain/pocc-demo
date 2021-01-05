package chainrpc

import (
	"fmt"
	"github.com/pocblockchain/pocc/x/mint"
)

func (c *Client) QueryMintParams() (mint.Params, error) {
	result := mint.Params{}

	err := c.Query(fmt.Sprintf("%s/%s/deposit", mint.QuerierRoute, mint.QueryParameters), nil, &result)
	return result, err
}
