package tezos

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/logger"
	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/dipdup-net/go-lib/tzkt/events"
	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/goat-systems/go-tezos/v4/forge"
	"github.com/goat-systems/go-tezos/v4/keys"
	"github.com/goat-systems/go-tezos/v4/rpc"
)

// Tezos -
type Tezos struct {
	cfg        Config
	rpc        *rpc.Client
	api        *api.API
	eventsTzKT *events.TzKT
	key        *keys.Key
	bigMaps    []api.BigMap
	counter    int64
	chainID    string
	minPayoff  decimal.Decimal

	log zerolog.Logger

	events     chan chain.Event
	operations chan chain.Operation

	stop chan struct{}
	wg   sync.WaitGroup
}

// Config -
type Config struct {
	Node      string
	TzKT      string
	MinPayOff string
	Contract  string
	Tokens    []string
	LogLevel  zerolog.Level
}

// New -
func New(cfg Config) (*Tezos, error) {
	if cfg.LogLevel == 0 {
		cfg.LogLevel = zerolog.InfoLevel
	}

	client, err := rpc.New(cfg.Node)
	if err != nil {
		return nil, err
	}

	secret, err := chain.LoadSecret("TEZOS_PRIVATE")
	if err != nil {
		return nil, err
	}

	key, err := keys.FromBase58(secret, keys.Ed25519)
	if err != nil {
		return nil, err
	}

	tez := &Tezos{
		cfg:        cfg,
		rpc:        client,
		key:        key,
		api:        api.New(cfg.TzKT),
		eventsTzKT: events.NewTzKT(fmt.Sprintf("%s/v1/events", cfg.TzKT)),
		log:        logger.New(logger.WithLogLevel(cfg.LogLevel), logger.WithModuleName("tezos")),
		events:     make(chan chain.Event, 1024*16),
		operations: make(chan chain.Operation, 1024),
		stop:       make(chan struct{}, 1),
	}

	tez.log.Info().Str("address", key.PubKey.GetAddress()).Msg("using address")

	minPayoff, err := decimal.NewFromString(cfg.MinPayOff)
	if err != nil {
		return nil, err
	}
	tez.minPayoff = minPayoff

	return tez, nil
}

// Wallet -
func (t *Tezos) Wallet() chain.Wallet {
	return chain.Wallet{
		Address: t.key.PubKey.GetAddress(),
		Private: t.key.GetBytes(),
	}
}

// Init -
func (t *Tezos) Init() error {
	t.log.Info().Msg("initializing...")
	blockHead := new(rpc.BlockIDHead)
	_, block, err := t.rpc.Block(blockHead)
	if err != nil {
		return err
	}
	t.chainID = block.ChainID

	_, counter, err := t.rpc.ContractCounter(rpc.ContractCounterInput{
		BlockID:    blockHead,
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
	return nil
}

// Run -
func (t *Tezos) Run() error {
	t.log.Info().Msg("running...")
	if err := t.eventsTzKT.Connect(); err != nil {
		return err
	}

	t.wg.Add(1)
	go t.listen()

	for _, bm := range t.bigMaps {
		if err := t.eventsTzKT.SubscribeToBigMaps(&bm.Ptr, bm.Contract.Address, ""); err != nil {
			return err
		}
	}

	if err := t.eventsTzKT.SubscribeToOperations(t.cfg.Contract, events.KindTransaction); err != nil {
		return err
	}

	for i := range t.cfg.Tokens {
		if err := t.eventsTzKT.SubscribeToOperations(t.cfg.Tokens[i], events.KindTransaction); err != nil {
			return err
		}
	}

	return nil
}

// Close -
func (t *Tezos) Close() error {
	t.log.Info().Msg("closing...")
	t.stop <- struct{}{}
	t.wg.Wait()

	close(t.events)
	close(t.operations)
	close(t.stop)
	return nil
}

// Events -
func (t *Tezos) Events() <-chan chain.Event {
	return t.events
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
		case update := <-t.eventsTzKT.Listen():
			switch update.Channel {
			case events.ChannelBigMap:
				if err := t.handleBigMapChannel(update); err != nil {
					t.log.Error().Err(err).Msg("handleBigMapChannel")
					continue
				}
			case events.ChannelOperations:
				if err := t.handleOperationsChannel(update); err != nil {
					t.log.Error().Err(err).Msg("handleOperationsChannel")
					continue
				}
			}
		}
	}
}

// Initiate -
func (t *Tezos) Initiate(args chain.InitiateArgs) error {
	t.log.Info().Str("hashed_secret", args.HashedSecret.String()).Msg("initiate")

	var value []byte
	var err error
	var amount string
	switch args.Contract {
	case t.cfg.Contract:
		amount = args.Amount.String()
		value, err = json.Marshal(map[string]interface{}{
			"prim": "Pair",
			"args": []map[string]interface{}{
				{
					"string": args.Participant,
				}, {
					"prim": "Pair",
					"args": []map[string]interface{}{
						{
							"prim": "Pair",
							"args": []map[string]interface{}{
								{
									"bytes": args.HashedSecret.String(),
								}, {
									"int": args.RefundTime.Unix(),
								},
							},
						}, {
							"int": args.PayOff.String(),
						},
					},
				},
			},
		})
	default:
		amount = "0"
		for i := range t.cfg.Tokens {
			if t.cfg.Tokens[i] == args.Contract {
				value, err = json.Marshal(map[string]interface{}{
					"prim": "Pair",
					"args": []map[string]interface{}{
						{
							"prim": "Pair",
							"args": []map[string]interface{}{
								{
									"prim": "Pair",
									"args": []map[string]interface{}{
										{
											"bytes": args.HashedSecret.String(),
										},
										{
											"string": args.Participant,
										},
									},
								},
								{
									"prim": "Pair",
									"args": []map[string]interface{}{
										{
											"int": args.PayOff.String(),
										},
										{
											"int": args.RefundTime.Unix(),
										},
									},
								},
							},
						}, {
							"prim": "Pair",
							"args": []map[string]interface{}{
								{
									"string": args.TokenAddress,
								},
								{
									"int": args.Amount.String(),
								},
							},
						},
					},
				})
				break
			}
		}
	}

	if err != nil {
		return err
	}

	if value == nil {
		return nil
	}

	opHash, err := t.sendTransaction(args.Contract, amount, "50000", "22000", "initiate", json.RawMessage(value))
	if err != nil {
		return err
	}
	t.operations <- chain.Operation{
		Status:       chain.Pending,
		Hash:         opHash,
		ChainType:    chain.ChainTypeTezos,
		HashedSecret: args.HashedSecret,
	}

	return nil
}

// Redeem -
func (t *Tezos) Redeem(hashedSecret, secret chain.Hex, contract string) error {
	t.log.Info().Str("hashed_secret", hashedSecret.String()).Str("contract", contract).Msg("redeem")

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
	t.log.Info().Str("hashed_secret", hashedSecret.String()).Str("contract", contract).Msg("refund")

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
	t.log.Info().Msg("restoring...")
	for _, bm := range t.bigMaps {
		if err := t.restoreFromBigMap(bm); err != nil {
			return err
		}
	}

	t.events <- chain.RestoredEvent{Chain: chain.ChainTypeTezos}
	return nil
}

func (t *Tezos) restoreFromBigMap(bm api.BigMap) error {
	limit := 100
	var end bool
	var offset int
	for !end {
		keys, err := t.api.GetBigmapKeys(uint64(bm.Ptr), map[string]string{
			"limit":  fmt.Sprintf("%d", limit),
			"offset": fmt.Sprintf("%d", offset),
		})
		if err != nil {
			return err
		}

		for i := range keys {
			if err := t.handleBigMapKey(keys[i], bm.Contract.Address); err != nil {
				return err
			}

			if err := t.restoreFinilizationSwap(bm, keys[i]); err != nil {
				return err
			}
		}

		end = len(keys) != limit
		offset += len(keys)
	}
	return nil
}

func (t *Tezos) restoreFinilizationSwap(bm api.BigMap, key api.BigMapKey) error {
	updates, err := t.api.GetBigmapKeyUpdates(uint64(bm.Ptr), key.Key, map[string]string{
		"sort.desc": "id",
	})
	if err != nil {
		return err
	}
	switch len(updates) {
	case 2:
		var value map[string]interface{}
		if err := json.Unmarshal(updates[0].Value, &value); err != nil {
			return err
		}
		return t.handleBigMapUpdate(BigMapUpdate{
			ID:     int64(updates[0].ID),
			Level:  int64(updates[0].Level),
			Bigmap: bm.Ptr,
			Contract: struct {
				Alias   string `mapstructure:"alias"`
				Address string `mapstructure:"address"`
			}{
				Address: bm.Contract.Address,
				Alias:   bm.Contract.Alias,
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
		})
	default:
		return nil
	}
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
	var value NewAtomexValue
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
			HashedSecretHex: chain.Hex(bigMapUpdate.Content.Key),
			Chain:           chain.ChainTypeTezos,
			ContractAddress: bigMapUpdate.Contract.Address,
			BlockNumber:     uint64(bigMapUpdate.Level),
			Initiator:       value.Initiator,
			Participant:     value.Participant,
			RefundTime:      refundTime,
		}

		if err := event.SetPayOff(value.Payoff, t.minPayoff); err != nil {
			if errors.Is(err, chain.ErrMinPayoff) {
				t.log.Warn().Str("hashed_secret", event.HashedSecretHex.String()).Msg("skip because of small pay off")
				return nil
			}
			return err
		}

		if err := event.SetAmountFromString(value.Amount); err != nil {
			return err
		}

		t.events <- event
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

				t.events <- chain.RedeemEvent{
					HashedSecretHex: chain.Hex(bigMapUpdate.Content.Key),
					Chain:           chain.ChainTypeTezos,
					ContractAddress: bigMapUpdate.Contract.Address,
					BlockNumber:     uint64(bigMapUpdate.Level),
					Secret:          secret,
				}
				return nil
			case "refund":
				t.events <- chain.RefundEvent{
					HashedSecretHex: chain.Hex(bigMapUpdate.Content.Key),
					Chain:           chain.ChainTypeTezos,
					ContractAddress: bigMapUpdate.Contract.Address,
					BlockNumber:     uint64(bigMapUpdate.Level),
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
			HashedSecretHex: chain.Hex(bigMapUpdate.Content.Key),
			Chain:           chain.ChainTypeTezos,
			ContractAddress: bigMapUpdate.Contract.Address,
			BlockNumber:     uint64(bigMapUpdate.Level),
			Initiator:       value.Initiator,
			Participant:     value.Participant,
			RefundTime:      refundTime,
		}

		if err := event.SetPayOff(value.Payoff, t.minPayoff); err != nil {
			if errors.Is(err, chain.ErrMinPayoff) {
				t.log.Warn().Str("hashed_secret", event.HashedSecretHex.String()).Msg("skip because of small pay off")
				return nil
			}
			return err
		}

		if err := event.SetAmountFromString(value.Amount); err != nil {
			return err
		}

		t.events <- event
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

				t.events <- chain.RedeemEvent{
					HashedSecretHex: chain.Hex(bigMapUpdate.Content.Key),
					Chain:           chain.ChainTypeTezos,
					ContractAddress: bigMapUpdate.Contract.Address,
					BlockNumber:     uint64(bigMapUpdate.Level),
					Secret:          secret,
				}
				return nil
			case "refund":
				t.events <- chain.RefundEvent{
					HashedSecretHex: chain.Hex(bigMapUpdate.Content.Key),
					Chain:           chain.ChainTypeTezos,
					ContractAddress: bigMapUpdate.Contract.Address,
					BlockNumber:     uint64(bigMapUpdate.Level),
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
	var value NewAtomexValue
	if err := json.Unmarshal(key.Value, &value); err != nil {
		return err
	}

	refundTime, err := time.Parse(time.RFC3339, value.RefundTime)
	if err != nil {
		return err
	}

	event := chain.InitEvent{
		HashedSecretHex: chain.Hex(key.Key),
		Chain:           chain.ChainTypeTezos,
		ContractAddress: contract,
		BlockNumber:     uint64(key.FirstLevel),
		Initiator:       value.Initiator,
		Participant:     value.Participant,
		RefundTime:      refundTime,
	}

	if err := event.SetPayOff(value.Payoff, t.minPayoff); err != nil {
		if errors.Is(err, chain.ErrMinPayoff) {
			t.log.Warn().Str("hashed_secret", event.HashedSecretHex.String()).Msg("skip because of small pay off")
			return nil
		}
		return err
	}

	if err := event.SetAmountFromString(value.Amount); err != nil {
		return err
	}

	t.events <- event
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
		HashedSecretHex: chain.Hex(key.Key),
		Chain:           chain.ChainTypeTezos,
		ContractAddress: contract,
		BlockNumber:     uint64(key.FirstLevel),
		Initiator:       value.Initiator,
		Participant:     value.Participant,
		RefundTime:      refundTime,
	}

	if err := event.SetPayOff(value.Payoff, t.minPayoff); err != nil {
		if errors.Is(err, chain.ErrMinPayoff) {
			t.log.Warn().Str("hashed_secret", event.HashedSecretHex.String()).Msg("skip because of small pay off")
			return nil
		}
		return err
	}

	if err := event.SetAmountFromString(value.Amount); err != nil {
		return err
	}

	t.events <- event
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
