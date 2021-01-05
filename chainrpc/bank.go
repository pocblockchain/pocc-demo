package chainrpc

import (
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/bank"
)

func (c *Client) Send(signer *Signer, to sdk.AccAddress, coins sdk.Coins, memo string) (*sdk.TxResponse, error) {
	msg := bank.MsgSend{
		FromAddress: signer.AccAddress(),
		ToAddress:   to,
		Amount:      coins,
	}
	return c.SignAndBroadcastTx(signer, msg, memo)
}
