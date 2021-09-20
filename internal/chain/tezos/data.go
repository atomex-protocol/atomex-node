package tezos

// BigMapUpdate -
type BigMapUpdate struct {
	ID       int64 `mapstructure:"id"`
	Level    int64 `mapstructure:"level"`
	Bigmap   int64 `mapstructure:"bigmap"`
	Contract struct {
		Alias   string `mapstructure:"alias"`
		Address string `mapstructure:"address"`
	} `mapstructure:"contract"`
	Path    string       `mapstructure:"path"`
	Action  BigMapAction `mapstructure:"action"`
	Content struct {
		Hash  string      `mapstructure:"hash"`
		Key   string      `mapstructure:"key"`
		Value interface{} `mapstructure:"value"`
	} `mapstructure:"content"`
}

// BigMapAction -
type BigMapAction string

// actions
const (
	BigMapActionAllocate  BigMapAction = "allocate"
	BigMapActionAddKey    BigMapAction = "add_key"
	BigMapActionRemoveKey BigMapAction = "remove_key"
	BigMapActionUpdateKey BigMapAction = "update_key"
)

// AtomexValue -
type AtomexValue struct {
	Settings struct {
		Amount     string `json:"amount" mapstructure:"amount"`
		Payoff     string `json:"payoff" mapstructure:"payoff"`
		RefundTime string `json:"refund_time" mapstructure:"refund_time"`
	} `json:"settings" mapstructure:"settings"`
	Recipients struct {
		Initiator   string `json:"initiator" mapstructure:"initiator"`
		Participant string `json:"participant" mapstructure:"participant"`
	} `json:"recipients" mapstructure:"recipients"`
}

// AtomexTokenValue -
type AtomexTokenValue struct {
	Amount       string `json:"totalAmount" mapstructure:"totalAmount"`
	Payoff       string `json:"payoffAmount" mapstructure:"payoffAmount"`
	RefundTime   string `json:"refundTime" mapstructure:"refundTime"`
	Initiator    string `json:"initiator" mapstructure:"initiator"`
	Participant  string `json:"participant" mapstructure:"participant"`
	TokenAddress string `json:"tokenAddress" mapstructure:"tokenAddress"`
}

// Transaction -
type Transaction struct {
	Type   string `mapstructure:"type" json:"type"`
	Hash   string `mapstructure:"hash" json:"hash"`
	Status string `mapstructure:"status" json:"status"`
}
