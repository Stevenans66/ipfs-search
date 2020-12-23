package config

import (
	"github.com/ipfs-search/ipfs-search/protocol/ipfs"
)

// IPFS specifies the configuration for the IPFS protocol.
type IPFS struct {
	APIURL     string `yaml:"api_url" env:"IPFS_API_URL"`
	GatewayURL string `yaml:"gateway_url"`
}

func IPFSDefaults() IPFS {
	return IPFS(*ipfs.DefaultConfig())
}

func (c *Config) IPFSConfig() *ipfs.Config {
	cfg := ipfs.Config(c.IPFS)
	return &cfg
}
