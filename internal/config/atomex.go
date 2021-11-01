package config

// Atomex -
type Atomex struct {
	Symbols  map[string]string `yaml:"symbols"`
	Settings AtomexSetiings    `yaml:"settings"`
	RestAPI  string            `yaml:"rest_api"`
	WsAPI    string            `yaml:"wss"`
}

// AtomexSetiings -
type AtomexSetiings struct {
	LockTime        int64   `yaml:"lock_time"`
	RewardForRedeem float64 `yaml:"reward_for_redeem"`
}
