package config

import (
	"fmt"
	"testing"

	"github.com/filecoin-project/go-address"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	// now := time.Now()
	// oneDayBefore := now.Add(-time.Hour * 24)
	cfg := Config{
		Filters: []*Filter{
			{
				Client: Address(address.TestAddress),
				Miner:  Address(address.TestAddress),
				Limit:  "10Tib",
				Used:   0,
				// Start:  Time(now),
				// End:    Time(oneDayBefore),
			},
			{
				Client: Address(address.TestAddress2),
				Miner:  Address(address.TestAddress2),
				Limit:  "8Gib",
				Used:   0,
				// Start:  Time(now),
				// End:    Time(oneDayBefore),
			},
		},
		WalletURL:   "127.0.0.1:8000",
		WalletToken: "token",
	}

	assert.NoError(t, cfg.SaveConfig())
}

func TestLoadConfig(t *testing.T) {
	cfg, err := InitConfig("./")
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	fmt.Printf("%+v\n", cfg)
	for _, f := range cfg.Filters {
		fmt.Printf("%+v\n", f)
	}
}
