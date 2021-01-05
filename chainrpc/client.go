package chainrpc

import (
	"fmt"
	"github.com/pocblockchain/pocc-demo/user"
	"github.com/pocblockchain/pocc/client/flags"
	"github.com/pocblockchain/pocc/codec"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto"
	rpcclient "github.com/tendermint/tendermint/rpc/client"
)

type Signer struct {
	*user.User
	AccountNumAndNonceManager
}

func NewSigner(privateKey crypto.PrivKey, m AccountNumAndNonceManager) *Signer {
	return &Signer{User: user.NewUser(privateKey), AccountNumAndNonceManager: m}
}

func NewSyncSignerFromPrivateKey(privateKey crypto.PrivKey) *Signer {
	return &Signer{
		User:                      user.NewUser(privateKey),
		AccountNumAndNonceManager: &SyncAccountNumAndNonceManager{},
	}
}

func NewSyncSigner(user *user.User) *Signer {
	return &Signer{
		User:                      user,
		AccountNumAndNonceManager: &SyncAccountNumAndNonceManager{},
	}
}

func NewSyncSignerWithMnemonic(mnemonic string) (*Signer, error) {
	privateKey, err := user.PrivateKeyFromMnemonic(mnemonic, "", uint32(0), uint32(0))
	if err != nil {
		return nil, fmt.Errorf("failed to get private key from mnemonic when new signer. error details: %v", err)
	}
	return &Signer{
		User:                      user.NewUser(privateKey),
		AccountNumAndNonceManager: &SyncAccountNumAndNonceManager{},
	}, nil
}

type AccountNumAndNonceManager interface {
	accountNumAndNonce(address sdk.AccAddress, c *Client) (uint64, uint64, error)
}

type SyncAccountNumAndNonceManager struct{}

func (s *SyncAccountNumAndNonceManager) accountNumAndNonce(address sdk.AccAddress, c *Client) (uint64, uint64, error) {
	acc, err := c.QueryAccInfo(address)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to query cu info when get nonce. error details: %v", err)
	}

	return acc.AccountNumber, acc.Sequence, nil
}

type Client struct {
	rpc           *rpcclient.HTTP
	cdc           *codec.Codec
	chainId       string
	gas           uint64
	minGasPrice   string
	broadcastMode string
}

func NewClient(codec *codec.Codec, rpcUrl, chainID, minGasPrice string, gas uint64) *Client {
	return &Client{
		rpc:           rpcclient.NewHTTP(rpcUrl, "/websocket"),
		cdc:           codec,
		chainId:       chainID,
		gas:           gas,
		minGasPrice:   minGasPrice,
		broadcastMode: flags.BroadcastSync,
	}
}

func (c *Client) Query(path string, params interface{}, result interface{}) error {

	paramdata, err := c.cdc.MarshalJSON(params)

	resultABCIQuery, err := c.rpc.ABCIQuery("custom/"+path, paramdata)
	if err != nil {
		return errors.Wrap(err, "failed to query "+path)
	}

	resp := resultABCIQuery.Response
	if !resp.IsOK() {
		return fmt.Errorf("failed to query with response log %s", resp.Log)
	}
	err = c.cdc.UnmarshalJSON(resp.Value, result)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal result ")
	}

	return nil
}

func (c *Client) SignAndBroadcastTx(signer *Signer, msg sdk.Msg, memo string) (*sdk.TxResponse, error) {
	signedMsg, err := c.SignMsg(signer, []sdk.Msg{msg}, memo)
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign msg")
	}

	var result sdk.TxResponse
	switch c.broadcastMode {
	case flags.BroadcastSync:
		res, err := c.rpc.BroadcastTxSync(signedMsg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to broadcast tx")
		}
		result = sdk.NewResponseFormatBroadcastTx(res)

	case flags.BroadcastAsync:
		res, err := c.rpc.BroadcastTxAsync(signedMsg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to broadcast tx")
		}
		result = sdk.NewResponseFormatBroadcastTx(res)
	case flags.BroadcastBlock:
		res, err := c.rpc.BroadcastTxCommit(signedMsg)
		if err != nil {
			return nil, errors.Wrap(err, "failed to broadcast tx")
		}
		result = sdk.NewResponseFormatBroadcastTxCommit(res)

	default:
		return nil, fmt.Errorf("unknow broadcast mode")
	}

	return &result, nil
}

func (c *Client) SignMsg(signer *Signer, msgs []sdk.Msg,memo string) ([]byte, error) {
	txEncoder := auth.DefaultTxEncoder(c.cdc)
	txBuilder := auth.NewTxBuilderFromCLI()

	accNum, nonce, err := signer.accountNumAndNonce(signer.AccAddress(), c)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce when sign msg: %v", err)
	}

	txBuilder = txBuilder.WithChainID(c.chainId).WithSequence(nonce).WithAccountNumber(accNum).WithGasPrices(c.minGasPrice).WithGas(c.gas).WithMemo(memo)
	stdSignMsg, err := txBuilder.BuildSignMsg(msgs)

	if err != nil {
		return nil, errors.Wrap(err, "failed to build msg")
	}

	sigBytes, err := signer.Sign(stdSignMsg.Bytes())
	if err != nil {
		return nil, errors.Wrap(err, "failed to sign msg")
	}

	sig := auth.StdSignature{
		PubKey:    signer.PubKey(),
		Signature: sigBytes,
	}

	signedMsg, err := txEncoder(auth.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, []auth.StdSignature{sig}, stdSignMsg.Memo))

	if err != nil {
		return nil, errors.Wrap(err, "failed to encode signed msg")
	}

	fmt.Printf("stdSignMsg:%+v\n", stdSignMsg)
	fmt.Printf("signedMsg:%+v\n", auth.NewStdTx(stdSignMsg.Msgs, stdSignMsg.Fee, []auth.StdSignature{sig}, stdSignMsg.Memo))
	return signedMsg, nil
}
