package config

// Atomex -
type Atomex struct {
	ToSymbols   map[string]string `yaml:"to_symbols"`
	FromSymbols map[string]string `yaml:"from_symbols"`
	Settings    AtomexSetiings    `yaml:"settings"`
	RestAPI     string            `yaml:"rest_api"`
	WsAPI       string            `yaml:"wss"`
}

// AtomexSetiings -
type AtomexSetiings struct {
	LockTime        int64   `yaml:"lock_time"`
	RewardForRedeem float64 `yaml:"reward_for_redeem"`
}
