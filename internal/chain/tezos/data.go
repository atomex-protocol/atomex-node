package tezos

// actions
const (
	BigMapActionAllocate  = "allocate"
	BigMapActionAddKey    = "add_key"
	BigMapActionRemoveKey = "remove_key"
	BigMapActionUpdateKey = "update_key"
)

// NewAtomexValue -
type NewAtomexValue struct {
	Amount      string `json:"amount"`
	Payoff      string `json:"payoff"`
	RefundTime  string `json:"refund_time"`
	Initiator   string `json:"initiator"`
	Participant string `json:"participant"`
}

// AtomexValue -
type AtomexValue struct {
	Settings struct {
		Amount     string `json:"amount"`
		Payoff     string `json:"payoff"`
		RefundTime string `json:"refund_time"`
	} `json:"settings"`
	Recipients struct {
		Initiator   string `json:"initiator"`
		Participant string `json:"participant"`
	} `json:"recipients"`
}

// AtomexTokenValue -
type AtomexTokenValue struct {
	Amount       string `json:"totalAmount"`
	Payoff       string `json:"payoffAmount"`
	RefundTime   string `json:"refundTime"`
	Initiator    string `json:"initiator"`
	Participant  string `json:"participant"`
	TokenAddress string `json:"tokenAddress"`
}

// Transaction -
type Transaction struct {
	Type   string `json:"type"`
	Hash   string `json:"hash"`
	Status string `json:"status"`
}

// OperationParamsByContracts -
type OperationParamsByContracts map[string]OperationParams

// OperationParams -
type OperationParams struct {
	GasLimit     ContractParams `yaml:"gas_limit"`
	StorageLimit ContractParams `yaml:"storage_limit"`
}

// ContractParams -
type ContractParams struct {
	Initiate string `yaml:"initiate"`
	Refund   string `yaml:"refund"`
	Redeem   string `yaml:"redeem"`
}
