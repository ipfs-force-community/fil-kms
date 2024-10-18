package config

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/docker/go-units"
	"github.com/filecoin-project/go-address"
	logging "github.com/ipfs/go-log/v2"
)

const (
	ConfigFile = "config.toml"
)

var log = logging.Logger("config")

var GlobalConfig *Config

type Config struct {
	WalletURL   string
	WalletToken string
	Filters     []*Filter

	bathPath string
}

type Filter struct {
	Client Address
	Miner  Address
	Limit  string
	limit  int64
	Used   int64
	// Start  Time
	// End    Time
}

func (f *Filter) GetLimit() int64 {
	return f.limit
}

func (f *Filter) SetLimit(i int64) {
	f.limit = i
}

func InitConfig(basePath string) (*Config, error) {
	has, err := hasConfig(filepath.Join(basePath, ConfigFile))
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, os.WriteFile(filepath.Join(basePath, ConfigFile), []byte{}, 0644)
		}
		return nil, err
	}
	if !has {
		return nil, fmt.Errorf("config file not found")
	}

	cfg, err := loadConfig(filepath.Join(basePath, ConfigFile))
	if err != nil {
		return nil, err
	}
	cfg.bathPath = basePath

	for _, filter := range cfg.Filters {
		if filter.Client.Empty() || filter.Miner.Empty() {
			log.Warnf("invalid filter: %v", filter)
			continue
		}

		if filter.limit, err = units.RAMInBytes(filter.Limit); err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func hasConfig(cfgPath string) (bool, error) {
	s, err := os.Stat(cfgPath)
	if err != nil {
		return false, err
	}

	return !s.IsDir(), nil
}

func loadConfig(cfgPath string) (*Config, error) {
	cfg := &Config{}
	f, err := os.Open(cfgPath)
	if err != nil {
		return nil, err
	}

	_, err = toml.NewDecoder(f).Decode(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) SaveConfig() error {
	buf := new(bytes.Buffer)
	err := toml.NewEncoder(buf).Encode(c)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(c.bathPath, ConfigFile), buf.Bytes(), 0644)
}

func (c *Config) GetFilter(clientAddr, minerAddr address.Address) *Filter {
	for _, filter := range c.Filters {
		if filter.Client.Address() == clientAddr && filter.Miner.Address() == minerAddr {
			return filter
		}
	}
	return nil
}

func (c *Config) SaveFilter(clientAddr, minerAddr address.Address, f *Filter) error {
	for idx, filter := range c.Filters {
		if filter.Client.Address() == clientAddr && filter.Miner.Address() == minerAddr {
			c.Filters[idx] = f
			return nil
		}
	}
	return fmt.Errorf("filter not found")
}
