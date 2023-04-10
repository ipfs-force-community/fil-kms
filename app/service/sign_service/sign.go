package sign_service

import (
	"fil-kms/app/keystore"
	"github.com/filecoin-project/lotus/chain/wallet"
)

var GlobalWalletService *walletService

type walletService struct {
	*wallet.LocalWallet
}

func NewWalletService(filepath string) (*walletService, error) {
	keyStore, err := keystore.NewKeyStore(filepath)
	if err != nil {
		return nil, err
	}

	localwallet, err := wallet.NewWallet(keyStore)
	if err != nil {
		return nil, err
	}

	return &walletService{localwallet}, nil
}
