package config

// Atomex -
type Atomex struct {
	ToSymbols   map[string]string `yaml:"to_symbols" validate:"required"`
	FromSymbols map[string]string `yaml:"from_symbols" validate:"required"`
	Settings    AtomexSettings    `yaml:"settings"`
	RestAPI     string            `yaml:"rest_api" validate:"required,uri"`
	WsAPI       string            `yaml:"wss" validate:"required,uri"`
	UptimeAPI   string            `yaml:"uptimeUri" validate:"required,uri"`
}

// AtomexSettings -
type AtomexSettings struct {
	LockTime        int64   `yaml:"lock_time" validate:"required,numeric"`
	RewardForRedeem float64 `yaml:"reward_for_redeem" validate:"numeric"`
}
