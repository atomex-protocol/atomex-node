package tezos

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aopoltorzhicky/watch_tower/internal/chain"
	"github.com/aopoltorzhicky/watch_tower/internal/config"
	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/dipdup-net/go-lib/tzkt/events"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	"github.com/goat-systems/go-tezos/v4/forge"
	"github.com/goat-systems/go-tezos/v4/keys"
	"github.com/goat-systems/go-tezos/v4/rpc"
)

// Tezos -
type Tezos struct {
	cfg       config.Tezos
	rpc       *rpc.Client
	api       *api.API
	events    *events.TzKT
	key       *keys.Key
	bigMaps   []api.BigMap
	counter   int64
	chainID   string
	minPayoff decimal.Decimal

	initChan   chan chain.InitEvent
	redeemChan chan chain.RedeemEvent
	refundChan chan chain.RefundEvent
	operations chan chain.Operation

	stop chan struct{}
	wg   sync.WaitGroup
}

// New -
func New(cfg config.Tezos) (*Tezos, error) {
	client, err := rpc.New(cfg.Node)
	if err != nil {
		return nil, err
	}

	key, err := keys.FromBase58(os.Getenv("TEZOS_PRIVATE"), keys.Ed25519)
	if err != nil {
		return nil, err
	}

	log.Info().Str("address", key.PubKey.GetAddress()).Str("blockchain", "tezos").Msg("using address")

	_, block, err := client.Block(&rpc.BlockIDHead{})
	if err != nil {
		return nil, err
	}

	minPayoff, err := decimal.NewFromString(cfg.MinPayOff)
	if err != nil {
		return nil, err
	}

	return &Tezos{
		cfg:        cfg,
		minPayoff:  minPayoff,
		rpc:        client,
		key:        key,
		api:        api.New(cfg.TzKT),
		chainID:    block.ChainID,
		events:     events.NewTzKT(fmt.Sprintf("%s/v1/events", cfg.TzKT)),
		initChan:   make(chan chain.InitEvent, 1024),
		redeemChan: make(chan chain.RedeemEvent, 1024),
		refundChan: make(chan chain.RefundEvent, 1024),
		operations: make(chan chain.Operation, 1024),
		stop:       make(chan struct{}, 1),
	}, nil
}

// Run -
func (t *Tezos) Run() error {
	_, counter, err := t.rpc.ContractCounter(rpc.ContractCounterInput{
		BlockID:    &rpc.BlockIDHead{},
		ContractID: t.key.PubKey.GetAddress(),
	})
	if err != nil {
		return err
	}
	atomic.StoreInt64(&t.counter, int64(counter))

	bigMaps, err := t.api.GetBigmaps(map[string]string{
		"contract.in": strings.Join(append(t.cfg.Tokens, t.cfg.Contract), ","),
	})
	if err != nil {
		return err
	}
	t.bigMaps = bigMaps

	if err := t.events.Connect(); err != nil {
		return err
	}

	t.wg.Add(1)
	go t.listen()

	for _, bm := range bigMaps {
		if err := t.events.SubscribeToBigMaps(&bm.Ptr, bm.Contract.Address, ""); err != nil {
			return err
		}
	}

	if err := t.events.SubscribeToOperations(t.cfg.Contract, events.KindTransaction); err != nil {
		return err
	}

	for i := range t.cfg.Tokens {
		if err := t.events.SubscribeToOperations(t.cfg.Tokens[i], events.KindTransaction); err != nil {
			return err
		}
	}

	return nil
}

// Close -
func (t *Tezos) Close() error {
	t.stop <- struct{}{}
	t.wg.Wait()

	close(t.initChan)
	close(t.redeemChan)
	close(t.refundChan)
	close(t.operations)
	close(t.stop)
	return nil
}

// InitEvents -
func (t *Tezos) InitEvents() <-chan chain.InitEvent {
	return t.initChan
}

// RedeemEvents -
func (t *Tezos) RedeemEvents() <-chan chain.RedeemEvent {
	return t.redeemChan
}

// RefundEvents -
func (t *Tezos) RefundEvents() <-chan chain.RefundEvent {
	return t.refundChan
}

// Operations -
func (t *Tezos) Operations() <-chan chain.Operation {
	return t.operations
}

func (t *Tezos) listen() {
	defer t.wg.Done()

	for {
		select {
		case <-t.stop:
			return
		case update := <-t.events.Listen():
			switch update.Channel {
			case events.ChannelBigMap:
				if err := t.handleBigMapChannel(update); err != nil {
					log.Err(err).Msg("handleBigMapChannel")
					continue
				}
			case events.ChannelOperations:
				if err := t.handleOperationsChannel(update); err != nil {
					log.Err(err).Msg("handleOperationsChannel")
					continue
				}
			}
		}
	}
}

// Redeem -
func (t *Tezos) Redeem(hashedSecret, secret chain.Hex, contract string) error {
	log.Info().Msg("redeeming...")
	value, err := json.Marshal(map[string]interface{}{
		"bytes": secret,
	})
	if err != nil {
		return err
	}

	opHash, err := t.sendTransaction(contract, "0", "50000", "22000", "redeem", json.RawMessage(value))
	if err != nil {
		return err
	}
	t.operations <- chain.Operation{
		Status:       chain.Pending,
		Hash:         opHash,
		ChainType:    chain.ChainTypeTezos,
		HashedSecret: hashedSecret,
	}
	return nil
}

// Refund -
func (t *Tezos) Refund(hashedSecret chain.Hex, contract string) error {
	log.Info().Str("hashed_secret", hashedSecret.String()).Msg("refunding...")

	value, err := json.Marshal(map[string]interface{}{
		"bytes": hashedSecret,
	})
	if err != nil {
		return err
	}

	opHash, err := t.sendTransaction(contract, "0", "50000", "22000", "refund", json.RawMessage(value))
	if err != nil {
		return err
	}
	t.operations <- chain.Operation{
		Status:       chain.Pending,
		Hash:         opHash,
		ChainType:    chain.ChainTypeTezos,
		HashedSecret: hashedSecret,
	}
	return nil
}

// Restore -
func (t *Tezos) Restore() error {
	for _, bm := range t.bigMaps {
		keys, err := t.api.GetBigmapKeys(uint64(bm.Ptr), map[string]string{
			"active": "true",
		})
		if err != nil {
			return err
		}

		log.Info().Int("count", len(keys)).Str("contract", bm.Contract.Address).Msg("initiated swaps found in tezos")
		for i := range keys {
			if err := t.handleBigMapKey(keys[i], bm.Contract.Address); err != nil {
				return err
			}

			if err := t.restoreRedeem(bm, keys[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Tezos) restoreRedeem(bm api.BigMap, key api.BigMapKey) error {
	for _, bigMap := range t.bigMaps {
		if bigMap.Contract.Address == bm.Contract.Address {
			continue
		}

		updates, err := t.api.GetBigmapKeyUpdates(uint64(bigMap.Ptr), key.Key, map[string]string{
			"sort.desc": "id",
		})
		if err != nil {
			return err
		}
		switch len(updates) {
		case 0:
			continue
		case 1:
			return nil
		case 2:
			var value map[string]interface{}
			if err := json.Unmarshal(updates[0].Value, &value); err != nil {
				return err
			}
			if err := t.handleBigMapUpdate(BigMapUpdate{
				ID:     int64(updates[0].ID),
				Level:  int64(updates[0].Level),
				Bigmap: bigMap.Ptr,
				Contract: struct {
					Alias   string `mapstructure:"alias"`
					Address string `mapstructure:"address"`
				}{
					Address: bigMap.Contract.Address,
					Alias:   bigMap.Contract.Alias,
				},
				Action: BigMapAction(updates[0].Action),
				Content: struct {
					Hash  string      `mapstructure:"hash"`
					Key   string      `mapstructure:"key"`
					Value interface{} `mapstructure:"value"`
				}{
					Hash:  key.Hash,
					Key:   key.Key,
					Value: value,
				},
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Tezos) handleOperationsChannel(update events.Message) error {
	if update.Body == nil {
		return nil
	}
	switch update.Type {
	case events.MessageTypeData:
		return t.handleOperationsData(update)
	case events.MessageTypeReorg:
	case events.MessageTypeState:
	}
	return nil
}

func (t *Tezos) handleOperationsData(update events.Message) error {
	var txs []Transaction
	if err := mapstructure.Decode(update.Body, &txs); err != nil {
		return err
	}

	for i := range txs {
		if txs[i].Type != events.KindTransaction {
			continue
		}

		t.operations <- chain.Operation{
			Hash:      txs[i].Hash,
			ChainType: chain.ChainTypeTezos,
			Status:    toOperationStatus(txs[i].Status),
		}
	}
	return nil
}

func (t *Tezos) handleBigMapChannel(update events.Message) error {
	if update.Body == nil {
		return nil
	}
	switch update.Type {
	case events.MessageTypeData:
		return t.handleBigMapData(update)
	case events.MessageTypeReorg:
	case events.MessageTypeState:
	}
	return nil
}

func (t *Tezos) handleBigMapData(update events.Message) error {
	var bigMapUpdates []BigMapUpdate
	if err := mapstructure.Decode(update.Body, &bigMapUpdates); err != nil {
		return err
	}

	for i := range bigMapUpdates {
		if bigMapUpdates[i].Action == BigMapActionAllocate {
			return nil
		}

		if err := t.handleBigMapUpdate(bigMapUpdates[i]); err != nil {
			return err
		}
	}
	return nil
}

func (t *Tezos) handleBigMapUpdate(bigMapUpdate BigMapUpdate) error {
	switch bigMapUpdate.Contract.Address {
	case t.cfg.Contract:
		return t.parseContractValueUpdate(bigMapUpdate)
	default:
		return t.parseTokensValueUpdate(bigMapUpdate)
	}
}

func (t *Tezos) parseContractValueUpdate(bigMapUpdate BigMapUpdate) error {
	var value AtomexValue
	if err := mapstructure.Decode(bigMapUpdate.Content.Value, &value); err != nil {
		return err
	}

	switch bigMapUpdate.Action {
	case BigMapActionAddKey:
		refundTime, err := time.Parse(time.RFC3339, value.Settings.RefundTime)
		if err != nil {
			return err
		}

		event := chain.InitEvent{
			Event: chain.Event{
				HashedSecret: chain.Hex(bigMapUpdate.Content.Key),
				Chain:        chain.ChainTypeTezos,
				Contract:     bigMapUpdate.Contract.Address,
			},
			Initiator:   value.Recipients.Initiator,
			Participant: value.Recipients.Participant,
			RefundTime:  refundTime,
		}

		if err := event.SetPayOff(value.Settings.Payoff, t.minPayoff); err != nil {
			if errors.Is(err, chain.ErrMinPayoff) {
				log.Warn().Str("hashed_secret", event.HashedSecret.String()).Msg("skip because of small pay off")
				return nil
			}
			return err
		}

		if err := event.SetAmountFromString(value.Settings.Amount); err != nil {
			return err
		}

		t.initChan <- event
	case BigMapActionUpdateKey:
	case BigMapActionRemoveKey:
		ops, err := t.api.GetTransactions(map[string]string{
			"level":  fmt.Sprintf("%d", bigMapUpdate.Level),
			"target": bigMapUpdate.Contract.Address,
		})
		if err != nil {
			return err
		}
		if len(ops) == 0 {
			return nil
		}
		for i := range ops {
			if ops[i].Parameters == nil {
				continue
			}
			switch ops[i].Parameters.Entrypoint {
			case "redeem":
				var secret chain.Hex
				if len(ops[i].Parameters.Value) >= 2 {
					s := string(ops[i].Parameters.Value)
					secret = chain.Hex(strings.Trim(s, "\""))
				} else {
					secret = chain.Hex(ops[i].Parameters.Value)
				}

				t.redeemChan <- chain.RedeemEvent{
					Event: chain.Event{
						HashedSecret: chain.Hex(bigMapUpdate.Content.Key),
						Chain:        chain.ChainTypeTezos,
						Contract:     bigMapUpdate.Contract.Address,
					},
					Secret: secret,
				}
				return nil
			case "refund":
				t.refundChan <- chain.RefundEvent{
					Event: chain.Event{
						HashedSecret: chain.Hex(bigMapUpdate.Content.Key),
						Chain:        chain.ChainTypeTezos,
						Contract:     bigMapUpdate.Contract.Address,
					},
				}
				return nil
			}
		}
	}

	return nil
}

func (t *Tezos) parseTokensValueUpdate(bigMapUpdate BigMapUpdate) error {
	var value AtomexTokenValue
	if err := mapstructure.Decode(bigMapUpdate.Content.Value, &value); err != nil {
		return err
	}

	switch bigMapUpdate.Action {
	case BigMapActionAddKey:
		refundTime, err := time.Parse(time.RFC3339, value.RefundTime)
		if err != nil {
			return err
		}
		event := chain.InitEvent{
			Event: chain.Event{
				HashedSecret: chain.Hex(bigMapUpdate.Content.Key),
				Chain:        chain.ChainTypeTezos,
				Contract:     bigMapUpdate.Contract.Address,
			},
			Initiator:   value.Initiator,
			Participant: value.Participant,
			RefundTime:  refundTime,
		}

		if err := event.SetPayOff(value.Payoff, t.minPayoff); err != nil {
			if errors.Is(err, chain.ErrMinPayoff) {
				log.Warn().Str("hashed_secret", event.HashedSecret.String()).Msg("skip because of small pay off")
				return nil
			}
			return err
		}

		if err := event.SetAmountFromString(value.Amount); err != nil {
			return err
		}

		t.initChan <- event
	case BigMapActionUpdateKey:
	case BigMapActionRemoveKey:
		ops, err := t.api.GetTransactions(map[string]string{
			"level":  fmt.Sprintf("%d", bigMapUpdate.Level),
			"target": bigMapUpdate.Contract.Address,
		})
		if err != nil {
			return err
		}
		if len(ops) == 0 {
			return nil
		}
		for i := range ops {
			if ops[i].Parameters == nil {
				continue
			}
			switch ops[i].Parameters.Entrypoint {
			case "redeem":
				var secret chain.Hex
				if len(ops[i].Parameters.Value) >= 2 {
					s := string(ops[i].Parameters.Value)
					secret = chain.Hex(strings.Trim(s, "\""))
				} else {
					secret = chain.Hex(ops[i].Parameters.Value)
				}

				t.redeemChan <- chain.RedeemEvent{
					Event: chain.Event{
						HashedSecret: chain.Hex(bigMapUpdate.Content.Key),
						Chain:        chain.ChainTypeTezos,
						Contract:     bigMapUpdate.Contract.Address,
					},
					Secret: secret,
				}
				return nil
			case "refund":
				t.refundChan <- chain.RefundEvent{
					Event: chain.Event{
						HashedSecret: chain.Hex(bigMapUpdate.Content.Key),
						Chain:        chain.ChainTypeTezos,
						Contract:     bigMapUpdate.Contract.Address,
					},
				}
				return nil
			}
		}
	}

	return nil
}

func (t *Tezos) handleBigMapKey(key api.BigMapKey, contract string) error {
	switch contract {
	case t.cfg.Contract:
		return t.parseContractValueKeys(key, contract)
	default:
		return t.parseTokensValueKeys(key, contract)
	}
}

func (t *Tezos) parseContractValueKeys(key api.BigMapKey, contract string) error {
	var value AtomexValue
	if err := json.Unmarshal(key.Value, &value); err != nil {
		return err
	}

	refundTime, err := time.Parse(time.RFC3339, value.Settings.RefundTime)
	if err != nil {
		return err
	}

	event := chain.InitEvent{
		Event: chain.Event{
			HashedSecret: chain.Hex(key.Key),
			Chain:        chain.ChainTypeTezos,
			Contract:     contract,
		},
		Initiator:   value.Recipients.Initiator,
		Participant: value.Recipients.Participant,
		RefundTime:  refundTime,
	}

	if err := event.SetPayOff(value.Settings.Payoff, t.minPayoff); err != nil {
		if errors.Is(err, chain.ErrMinPayoff) {
			log.Warn().Str("hashed_secret", event.HashedSecret.String()).Msg("skip because of small pay off")
			return nil
		}
		return err
	}

	if err := event.SetAmountFromString(value.Settings.Amount); err != nil {
		return err
	}

	t.initChan <- event
	return nil
}

func (t *Tezos) parseTokensValueKeys(key api.BigMapKey, contract string) error {
	var value AtomexTokenValue
	if err := json.Unmarshal(key.Value, &value); err != nil {
		return err
	}

	refundTime, err := time.Parse(time.RFC3339, value.RefundTime)
	if err != nil {
		return err
	}

	event := chain.InitEvent{
		Event: chain.Event{
			HashedSecret: chain.Hex(key.Key),
			Chain:        chain.ChainTypeTezos,
			Contract:     contract,
		},
		Initiator:   value.Initiator,
		Participant: value.Participant,
		RefundTime:  refundTime,
	}

	if err := event.SetPayOff(value.Payoff, t.minPayoff); err != nil {
		if errors.Is(err, chain.ErrMinPayoff) {
			log.Warn().Str("hashed_secret", event.HashedSecret.String()).Msg("skip because of small pay off")
			return nil
		}
		return err
	}

	if err := event.SetAmountFromString(value.Amount); err != nil {
		return err
	}

	t.initChan <- event
	return nil
}

func (t *Tezos) sendTransaction(destination, amount, fee, gasLimit, entrypoint string, value json.RawMessage) (string, error) {
	atomic.AddInt64(&t.counter, 1)

	transaction := rpc.Transaction{
		Kind:         rpc.TRANSACTION,
		Source:       t.key.PubKey.GetAddress(),
		Fee:          fee,
		GasLimit:     gasLimit,
		StorageLimit: "257",
		Counter:      strconv.Itoa(int(t.counter)),
		Amount:       amount,
		Destination:  destination,
		Parameters: &rpc.Parameters{
			Entrypoint: entrypoint,
			Value:      &value,
		},
	}

	_, block, err := t.rpc.Block(&blockIDHead{"60"})
	if err != nil {
		return "", err
	}

	op, err := forge.Encode(block.Hash, transaction.ToContent())
	if err != nil {
		return "", err
	}

	signature, err := t.key.SignHex(op)
	if err != nil {
		return "", err
	}

	_, opHash, err := t.rpc.InjectionOperation(rpc.InjectionOperationInput{
		Operation: signature.AppendToHex(op),
		ChainID:   t.chainID,
	})

	return opHash, err
}
