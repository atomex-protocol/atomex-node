package tezos

// actions
const (
	BigMapActionAllocate  = "allocate"
	BigMapActionAddKey    = "add_key"
	BigMapActionRemoveKey = "remove_key"
	BigMapActionUpdateKey = "update_key"
)

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
