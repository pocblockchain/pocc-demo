package tests

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/pocblockchain/pocc-demo/chainrpc"
	"github.com/pocblockchain/pocc-demo/config"
	"github.com/pocblockchain/pocc-demo/user"
	"github.com/pocblockchain/pocc/pocapp"
	sdk "github.com/pocblockchain/pocc/types"
	"github.com/pocblockchain/pocc/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"strings"
	"testing"
	"time"
)

func MyHash(bz []byte) []byte {
	h := sha256.Sum256(bz)
	return h[:]
}

func TestQueryBlock(t *testing.T) {
	cfg, _ := config.GetConfig("../config/config.yml")
	cdc := pocapp.MakeCodec()
	c := chainrpc.NewClient(pocapp.MakeCodec(), cfg.RPCURL, cfg.ChainID, cfg.MinGasPrice, cfg.Gas)

	height := int64(585380)

	block, _ := c.QueryBlock(&height)

	for _, tx := range block.Block.Txs {
		t.Logf("tx:%v\n", tx)
		hash := tx.Hash()
		myHash := MyHash(tx)
		require.Equal(t, hash, myHash)
		resp, _ := c.QueryTxByHash(hex.EncodeToString(hash))
		t.Logf("resp:%+v\n", resp)
		require.EqualValues(t, tx.Hash(), resp.Hash)

		//decode tx
		var tx types.StdTx
		err := cdc.UnmarshalBinaryLengthPrefixed(resp.Tx, &tx)
		if err != nil {
			continue
		}
		t.Logf("tx:%+v\n",tx)
		t.Logf("memo:%v", tx.GetMemo())
	}
}





func TestGenPrivateKey(t *testing.T){
	priv := secp256k1.GenPrivKey()
	t.Logf("priv:%v, addr:%s", hex.EncodeToString(priv[:]), sdk.AccAddress(priv.PubKey().Address()))
}

//c3d109930aec46cf62d39681387954be78d7c78da80548486725e92d60e315d6
//poc1jzs5vjfx2zu94lujxhq6g384nfwn9m7yma4uhn

func TestBasicSend(t *testing.T) {
	cfg, err := config.GetConfig("../config/config.yml")
    bz, err := hex.DecodeString("c3d109930aec46cf62d39681387954be78d7c78da80548486725e92d60e315d6")

    var PrivBz [32]byte
    copy(PrivBz[:], bz)
	sender := user.NewUser(secp256k1.PrivKeySecp256k1(PrivBz))
	c := chainrpc.NewClient(pocapp.MakeCodec(), cfg.RPCURL, cfg.ChainID, cfg.MinGasPrice, cfg.Gas)
	signer := chainrpc.NewSyncSigner(sender)

	to, err := sdk.AccAddressFromBech32("poc1y3dmx4un7sqetfxmljevpltkp4pqf5whtv88u9")
	require.Nil(t, err)

	resp, err := c.Send(signer, to, sdk.NewCoins(sdk.NewCoin(sdk.NativeToken, sdk.NewInt(10000))), "123")
	require.Nil(t, err)
	require.EqualValues(t, sdk.CodeOK, resp.Code)
	t.Logf("hash:%v", resp.TxHash)
	t.Logf("resp:%v", resp)

	time.Sleep(5 * time.Second)
	resp1, err := c.QueryTxByHash(resp.TxHash)
	require.EqualValues(t, resp.TxHash, strings.ToUpper(hex.EncodeToString(resp1.Hash)))

	acc, err := c.QueryAccInfo(to)
	t.Logf("acc:%v", acc)

	amount := acc.Coins.AmountOf(sdk.NativeToken)
	amountDec := sdk.NewDecFromInt(amount).MulTruncate(sdk.NewDecWithPrec(1, 18))
	t.Logf("amount:%v", amountDec)

}


func TestVerifyAddress(t *testing.T){
	addr, err := sdk.AccAddressFromBech32("poc1tfknn0jemvnrupqv04he2mzp7lsngwflh6x0gm")
	require.Nil(t, err)
	t.Logf("addr:%v", addr)
}



func TestVerifyAccAndValAddress (t *testing.T){
	data := []struct{
		accAddr string
		valAddr string

	}{
		{"poc1v3fhdccyz8qahacm68gduae4rr3z8z3q9qdscd", "pocvaloper1v3fhdccyz8qahacm68gduae4rr3z8z3q5smcx0"},
		{"poc10hdaf4y26zp9fsd43gaht0tsmsn60zf0h0hlcn", "pocvaloper10hdaf4y26zp9fsd43gaht0tsmsn60zf0xlphx3"},
		{"poc10ncxv4pcgpg7vaevdlpnp8f46exhpu8glhxajx","pocvaloper10ncxv4pcgpg7vaevdlpnp8f46exhpu8gw8s4vy"},
	}

	//bechAccPrefix = "poc"
	//bechValPrefix = "pocvaloper"
	//
	for _, d := range data {
		addr1,err :=  sdk.AccAddressFromBech32(d.accAddr)
		require.Nil(t, err)
		addr2, err := sdk.ValAddressFromBech32(d.valAddr)
		require.Nil(t, err)
		require.Equal(t, addr1.Bytes(), addr2.Bytes())
	}
}


//
//// GetFromBech32 decodes a bytestring from a Bech32 encoded string.
//func GetFromBech32(bech32str, prefix string) ([]byte, error) {
//	if len(bech32str) == 0 {
//		return nil, errors.New("decoding Bech32 address failed: must provide an address")
//	}
//
//	hrp, bz, err := bech32.DecodeAndConvert(bech32str)
//	if err != nil {
//		return nil, err
//	}
//
//	if hrp != prefix {
//		return nil, fmt.Errorf("invalid Bech32 prefix; expected %s, got %s", prefix, hrp)
//	}
//
//	return bz, nil
//}
