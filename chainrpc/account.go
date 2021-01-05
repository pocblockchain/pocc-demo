package chainrpc

import (
	"fmt"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth"
)

func (c *Client) QueryAccInfo(cuadress sdk.AccAddress) (*auth.BaseAccount, error) {
	params := auth.QueryAccountParams{
		Address: cuadress,
	}

	var acc auth.BaseAccount
	err := c.Query(fmt.Sprintf("%s/%s", auth.QuerierRoute, auth.QueryAccount), &params, &acc)
	if err != nil {
		return nil, err
	}
	return &acc, err
}

