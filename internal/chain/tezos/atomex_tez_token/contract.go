// DO NOT EDIT!!!
package atomexteztoken

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	jsoniter "github.com/json-iterator/go"

	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/dipdup-net/go-lib/tzkt/events"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// entrypoint names
const (   
	EntrypointInitiate = "initiate"   
	EntrypointRedeem = "redeem"   
	EntrypointRefund = "refund"
)


// InitiateTx - `initiate` transaction entity
type InitiateTx struct {
	*events.Transaction
	Initiate Initiate
}

// RedeemTx - `redeem` transaction entity
type RedeemTx struct {
	*events.Transaction
	Redeem Redeem
}

// RefundTx - `refund` transaction entity
type RefundTx struct {
	*events.Transaction
	Refund Refund
}


// BigMap0Update - `BigMap0` update entity
type BigMap0Update struct {
	BigMap0 BigMap0
	Level uint64
	Contract string
	Action string
}

// Atomexteztoken - struct which implementing contract interaction
type Atomexteztoken struct {
	tzktAPI *api.API
	tzktEvents *events.TzKT
	address string

	initiate chan InitiateTx
	redeem chan RedeemTx
	refund chan RefundTx
	bigmap0 chan BigMap0Update

	wg sync.WaitGroup
}

// New - constructor of contract entity
func New(baseURL string) *Atomexteztoken  {
	return &Atomexteztoken {
		tzktAPI: api.New(baseURL),
		tzktEvents: events.NewTzKT(fmt.Sprintf("%s/v1/events", baseURL)),
		address: "KT1KHdx4EC7WNVeeTGMr8utCSRczgu2EPCzX",
		initiate: make(chan InitiateTx, 1024),
		redeem: make(chan RedeemTx, 1024),
		refund: make(chan RefundTx, 1024),
		bigmap0: make(chan BigMap0Update, 1024),
	}
}

// ChangeAddress - replaces using contract address. Default: value from generating arguments.
func (contract *Atomexteztoken) ChangeAddress(address string) {
	contract.address = address
}

// Subscribe - subscribe on all contract's transaction
func (contract *Atomexteztoken) Subscribe(ctx context.Context) error {
	if err := contract.tzktEvents.Connect(); err != nil {
		return err
	}

	contract.wg.Add(1)
	go contract.listen(ctx)

	if err := contract.tzktEvents.SubscribeToBigMaps(nil, contract.address, ""); err != nil {
		return err
	}

	return contract.tzktEvents.SubscribeToOperations(contract.address, api.KindTransaction)
}

// Close - close all contract's connections
func (contract *Atomexteztoken) Close() error {
	contract.wg.Wait()

	if err := contract.tzktEvents.Close(); err != nil {
		return err
	}
	close(contract.initiate)
	close(contract.redeem)
	close(contract.refund)
	return nil
}


// InitiateEvents - listen `initiate` events channel
func (contract *Atomexteztoken) InitiateEvents() <-chan InitiateTx {
	return contract.initiate
}

// RedeemEvents - listen `redeem` events channel
func (contract *Atomexteztoken) RedeemEvents() <-chan RedeemTx {
	return contract.redeem
}

// RefundEvents - listen `refund` events channel
func (contract *Atomexteztoken) RefundEvents() <-chan RefundTx {
	return contract.refund
}


// BigMap0Updates - listen `BigMap0` updates channel
func (contract *Atomexteztoken) BigMap0Updates() <-chan BigMap0Update {
	return contract.bigmap0
}


func (contract *Atomexteztoken) listen(ctx context.Context) {
	defer contract.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <- contract.tzktEvents.Listen():
			switch msg.Type {
			case events.MessageTypeData:

				switch msg.Channel {
				case events.ChannelBigMap:
					items := msg.Body.([]events.BigMapUpdate)
					for i := range items {
						switch items[i].Path {
						case "0":
							var key Key0
							if err := json.Unmarshal(items[i].Content.Key, &key); err != nil {
								log.Println(err)
								continue
							}

							var value Value0
							if err := json.Unmarshal(items[i].Content.Value, &value); err != nil {
								log.Println(err)
								continue
							}
							contract.bigmap0 <- BigMap0Update{
								BigMap0: BigMap0{
									Key: key,
									Value: value,
								},
								Level: items[i].Level,
								Action: items[i].Action,
								Contract: contract.address,								
							}
						}
					}

				case events.ChannelOperations:
					items := msg.Body.([]interface{})
					for _, item := range items {
						tx, ok := item.(*events.Transaction)
						if !ok {
							continue
						}
						if tx.Parameter == nil {
							continue
						}

						switch tx.Parameter.Entrypoint {
						case "initiate":
							var data Initiate
							if err := json.Unmarshal(tx.Parameter.Value, &data); err != nil {
								log.Println(err)
								continue
							}
							contract.initiate <- InitiateTx{
								tx, data,
							}
						case "redeem":
							var data Redeem
							if err := json.Unmarshal(tx.Parameter.Value, &data); err != nil {
								log.Println(err)
								continue
							}
							contract.redeem <- RedeemTx{
								tx, data,
							}
						case "refund":
							var data Refund
							if err := json.Unmarshal(tx.Parameter.Value, &data); err != nil {
								log.Println(err)
								continue
							}
							contract.refund <- RefundTx{
								tx, data,
							}
						}
					}
				}

			case events.MessageTypeReorg:
			case events.MessageTypeState:
			case events.MessageTypeSubscribed:
			}
		}
	}
}

   
// GetInitiate - get `initiate` transactions
func (contract *Atomexteztoken) GetInitiate(ctx context.Context, page Page) ([]Initiate, error) {
	operations, err := getTransactions(ctx, contract.tzktAPI, "initiate", contract.address, page)
	if err != nil {
		return nil, err
	}
	values := make([]Initiate, 0)
	for i := range operations {
		if operations[i].Parameter == nil {
			continue
		}
		var value Initiate
		if err := json.Unmarshal(operations[i].Parameter.Value, &value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}
   
// GetRedeem - get `redeem` transactions
func (contract *Atomexteztoken) GetRedeem(ctx context.Context, page Page) ([]Redeem, error) {
	operations, err := getTransactions(ctx, contract.tzktAPI, "redeem", contract.address, page)
	if err != nil {
		return nil, err
	}
	values := make([]Redeem, 0)
	for i := range operations {
		if operations[i].Parameter == nil {
			continue
		}
		var value Redeem
		if err := json.Unmarshal(operations[i].Parameter.Value, &value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}
   
// GetRefund - get `refund` transactions
func (contract *Atomexteztoken) GetRefund(ctx context.Context, page Page) ([]Refund, error) {
	operations, err := getTransactions(ctx, contract.tzktAPI, "refund", contract.address, page)
	if err != nil {
		return nil, err
	}
	values := make([]Refund, 0)
	for i := range operations {
		if operations[i].Parameter == nil {
			continue
		}
		var value Refund
		if err := json.Unmarshal(operations[i].Parameter.Value, &value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

// Page -
type Page struct {
	Limit  uint64
	Offset uint64
}

func getLimits(p Page) Page {
	var newPage Page
	if p.Limit == 0 || p.Limit > 10000 {
		newPage.Limit = 100
	} else {
		newPage.Limit = p.Limit
	}

	if p.Offset == 0 || p.Offset > 10000 {
		newPage.Offset = 100
	} else {
		newPage.Offset = p.Offset
	}

	return newPage
}

func getTransactions(ctx context.Context, tzktAPI *api.API, entrypoint, contract string, page Page) ([]api.Transaction, error) {
	limits := getLimits(page)
	return tzktAPI.GetTransactions(ctx, map[string]string{
		"entrypoint": entrypoint,
		"target":     contract,
		"limit":      strconv.FormatUint(limits.Limit, 10),
		"offset":     strconv.FormatUint(limits.Offset, 10),
	})
}

// GetStorage - get `KT1KHdx4EC7WNVeeTGMr8utCSRczgu2EPCzX` current storage
func (contract *Atomexteztoken) GetStorage(ctx context.Context) (Storage, error) {
	var storage Storage
	err := contract.tzktAPI.GetContractStorage(ctx, contract.address, &storage)
	return storage, err
}
   
// BuildInitiateParameters - build `initiate` parameters
func (contract *Atomexteztoken) BuildInitiateParameters(ctx context.Context, params Initiate) ([]byte, error) {
	return contract.tzktAPI.BuildContractParameters(ctx, contract.address, "initiate", params)
}
   
// BuildRedeemParameters - build `redeem` parameters
func (contract *Atomexteztoken) BuildRedeemParameters(ctx context.Context, params Redeem) ([]byte, error) {
	return contract.tzktAPI.BuildContractParameters(ctx, contract.address, "redeem", params)
}
   
// BuildRefundParameters - build `refund` parameters
func (contract *Atomexteztoken) BuildRefundParameters(ctx context.Context, params Refund) ([]byte, error) {
	return contract.tzktAPI.BuildContractParameters(ctx, contract.address, "refund", params)
}

