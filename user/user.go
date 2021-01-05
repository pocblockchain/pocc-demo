package user

import (
	"github.com/pocblockchain/pocc/crypto/keys/hd"
	"github.com/pocblockchain/pocc/types"
	"github.com/cosmos/go-bip39"
	tmcrypto "github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

type User struct {
	privateKey tmcrypto.PrivKey
}

func NewUser(privateKey tmcrypto.PrivKey) *User {
	return &User{privateKey: privateKey}
}

func (s *User) PrivateKey() tmcrypto.PrivKey {
	return s.privateKey
}

func (s *User) Sign(bytes []byte) ([]byte, error) {
	return s.privateKey.Sign(bytes)

}

func (s *User) PubKey() tmcrypto.PubKey {
	return s.privateKey.PubKey()
}

func (s *User) AccAddress() types.AccAddress {
	return types.AccAddressFromPubKey(s.PubKey())
}

func NewUserWithMnemonic(mnemonic string) *User {
	privateKey, err := PrivateKeyFromMnemonic(mnemonic, "", uint32(0), uint32(0))
	if err != nil {
		panic(err)
	}
	return NewUser(privateKey)
}

func PrivateKeyFromMnemonic(mnemonic string, bip39Passphrase string, account uint32, index uint32) (*secp256k1.PrivKeySecp256k1, error) {
	coinType := types.GetConfig().GetCoinType()
	hdPath := hd.NewFundraiserParams(account, coinType, index)
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, bip39Passphrase)
	if err != nil {
		return nil, err
	}

	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, hdPath.String())
	priv := secp256k1.PrivKeySecp256k1(derivedPriv)

	return &priv, nil
}
