package sign_service

import (
	"context"
	"fil-kms/app/config"
)

func Init(storePath string, cfg *config.Config) error {
	var wapi IWallet
	var err error
	if len(cfg.WalletURL) != 0 {
		wapi, err = NewRemoteWallet(context.Background(), cfg.WalletURL, cfg.WalletToken)
		if err != nil {
			return err
		}
	} else {
		wapi, err = NewWalletService(storePath)
		if err != nil {
			return err
		}
	}

	GlobalWalletService = NewLimitedWallet(wapi, cfg)
	return nil
}
