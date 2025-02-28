package signer

import (
	"fmt"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
)

type hdWallet struct {
	MasterKey      *hdkeychain.ExtendedKey
	NextChildIndex uint32
	Paths          map[string]string
}

func (s *hdWallet) InitHDMaster(params *chaincfg.Params) (*hdkeychain.ExtendedKey, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	masterKey, err := hdkeychain.NewMaster(seed, params)
	if err != nil {
		fmt.Println(err)
	}

	return masterKey, nil
}

func (s *hdWallet) DeriveFromParent(parent *hdkeychain.ExtendedKey) (*hdkeychain.ExtendedKey, error) {
	parent.Depth()
	return parent, nil
}

func (s *hdWallet) defaultBip44Path() []uint32 {
	return []uint32{
		44 + hdkeychain.HardenedKeyStart,
		0 + hdkeychain.HardenedKeyStart,
		0 + hdkeychain.HardenedKeyStart,
		0,
		0,
	}
}

func (s *hdWallet) deriveCustomBip44Path(coinType, account, change, index uint32) []uint32 {
	return []uint32{
		44 + hdkeychain.HardenedKeyStart, // BIP44 proposal
		coinType + hdkeychain.HardenedKeyStart,
		account + hdkeychain.HardenedKeyStart,
		change,
		index,
	}
}
