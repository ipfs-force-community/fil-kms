package sign_service

import (
	"context"
	"fil-kms/app/config"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/builtin/v10/market"
	"github.com/filecoin-project/lotus/api"
	"github.com/stretchr/testify/assert"
)

func TestLimitWallet(t *testing.T) {
	client, err := address.NewFromString("f3tginto55sulyzzjoy3byof7u7t5vicwa372z4hn4zdhjkdwdkhalsn6mytk53paoc63uf575rfdvqetyvpka")
	assert.NoError(t, err)
	miner := address.TestAddress2
	cfg := &config.Config{
		WalletURL:   "/ip4/192.168.200.48/tcp/5678/http",
		WalletToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJyZWFkIiwid3JpdGUiLCJzaWduIiwiYWRtaW4iXX0.ityzQ7Xkb7C_K-JzZODVWoz_dQSoDceMSovGCHRAzQI",
		Filters: []*config.Filter{
			{
				Client: config.Address(client),
				Miner:  config.Address(miner),
				Limit:  "10Gib",
			},
		},
	}
	cfg.Filters[0].SetLimit(10 * 1024 * 1024 * 1024)

	wapi, err := NewRemoteWallet(context.Background(), cfg.WalletURL, cfg.WalletToken)
	assert.NoError(t, err)

	lwapi := NewLimitedWallet(wapi, cfg)
	go lwapi.Start(context.Background())
	defer lwapi.Stop()

	prop := &market.DealProposal{
		PieceSize: 1,
		Client:    client,
		Provider:  miner,
	}
	for i := 0; i < 10; i++ {
		sig, err := lwapi.WalletSign(context.Background(), client, []byte("hello"), api.MsgMeta{Type: api.MTUnknown}, prop)
		assert.NoError(t, err)
		assert.NotNil(t, sig)
	}
}
