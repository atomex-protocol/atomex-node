// DO NOT EDIT!!!
package atomexteztoken

import (
	"github.com/dipdup-net/go-lib/tools/tezgen"
	"strconv"
)

// Initiate
type Initiate struct {
	HashedSecret tezgen.Bytes     `json:"hashedSecret" validate:"string"`
	Participant  tezgen.Address   `json:"participant" validate:"string"`
	PayoffAmount tezgen.Int       `json:"payoffAmount,string" validate:"string"`
	RefundTime   tezgen.Timestamp `json:"refundTime" validate:"string"`
	TokenAddress tezgen.Address   `json:"tokenAddress" validate:"string"`
	TotalAmount  tezgen.Int       `json:"totalAmount,string" validate:"string"`
}

// Redeem
type Redeem tezgen.Bytes

// Refund
type Refund tezgen.Bytes

// Key0
type Key0 tezgen.Bytes

// Value0
type Value0 struct {
	Initiator    tezgen.Address   `json:"initiator" validate:"string"`
	Participant  tezgen.Address   `json:"participant" validate:"string"`
	PayoffAmount tezgen.Int       `json:"payoffAmount,string" validate:"string"`
	RefundTime   tezgen.Timestamp `json:"refundTime" validate:"string"`
	TokenAddress tezgen.Address   `json:"tokenAddress" validate:"string"`
	TotalAmount  tezgen.Int       `json:"totalAmount,string" validate:"string"`
}

// BigMap0
type BigMap0 struct {
	Key   Key0
	Value Value0
	Ptr   *int64
}

// UnmarshalJSON
func (b *BigMap0) UnmarshalJSON(data []byte) error {
	ptr, err := strconv.ParseInt(string(data), 10, 64)
	if err == nil {
		b.Ptr = &ptr
		return nil
	}
	parts := []interface{}{b.Key, b.Value}
	return json.Unmarshal(data, &parts)
}

// Storage
type Storage struct {
	BigMap0 BigMap0     `json:"0" validate:"string"`
	Unit1   tezgen.Unit `json:"1" validate:"string"`
}
