package ipfs

// Config specifies the configuration for the IPFS protocol.
type Config struct {
	APIURL     string // URL of an IPFS API endpoint (for Ls and Stat calls).
	GatewayURL string // URL of an IPFS Gateway (to request content).
}

// DefaultConfig returns the default configuration for a Sniffer.
func DefaultConfig() *Config {
	return &Config{
		APIURL:     "http://localhost:5001",
		GatewayURL: "http://localhost:8080",
	}
}
