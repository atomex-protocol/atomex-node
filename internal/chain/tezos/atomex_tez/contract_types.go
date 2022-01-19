// DO NOT EDIT!!!
package atomextez

import (
	"github.com/dipdup-net/go-lib/tools/tezgen"
	"strconv"
)

// Initiate
type Initiate struct {
	Settings    Settings       `json:"settings" validate:"string"`
	Participant tezgen.Address `json:"participant" validate:"string"`
}

// Settings
type Settings struct {
	RefundTime   tezgen.Timestamp `json:"refund_time" validate:"string"`
	Payoff       tezgen.Int       `json:"payoff" validate:"string"`
	HashedSecret tezgen.Bytes     `json:"hashed_secret" validate:"string"`
}

// Add
type Add tezgen.Bytes

// Redeem
type Redeem tezgen.Bytes

// Refund
type Refund tezgen.Bytes

// KeyBigMap
type KeyBigMap tezgen.Bytes

// ValueBigMap
type ValueBigMap struct {
	Recipients Recipients `json:"recipients" validate:"string"`
	Settings   Settings1  `json:"settings" validate:"string"`
}

// Recipients
type Recipients struct {
	Participant tezgen.Address `json:"participant" validate:"string"`
	Initiator   tezgen.Address `json:"initiator" validate:"string"`
}

// Settings1
type Settings1 struct {
	Amount     tezgen.Int       `json:"amount" validate:"string"`
	RefundTime tezgen.Timestamp `json:"refund_time" validate:"string"`
	Payoff     tezgen.Int       `json:"payoff" validate:"string"`
}

// BigMap
type BigMap struct {
	Key   KeyBigMap
	Value ValueBigMap
	Ptr   *int64
}

// UnmarshalJSON
func (b *BigMap) UnmarshalJSON(data []byte) error {
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
	BigMap BigMap      `json:"big_map" validate:"string"`
	Unit   tezgen.Unit `json:"unit" validate:"string"`
}
