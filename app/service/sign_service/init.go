package sign_service

func Init(storePath string) error {
	walletService, err := NewWalletService(storePath)
	if err != nil {
		return err
	}

	GlobalWalletService = walletService
	return nil
}
