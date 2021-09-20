package chain

// Operation -
type Operation struct {
	Hash         string
	ChainType    ChainType
	Status       OperationStatus
	HashedSecret Hex
}

// OperationStatus -
type OperationStatus int

const (
	Pending OperationStatus = iota + 1
	Applied
	Failed
)

// String -
func (status OperationStatus) String() string {
	switch status {
	case Pending:
		return "pending"
	case Applied:
		return "applied"
	case Failed:
		return "failed"
	}
	return "unknown"
}
