package tezos

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	atomextez "github.com/atomex-protocol/watch_tower/internal/chain/tezos/atomex_tez"
	atomexteztoken "github.com/atomex-protocol/watch_tower/internal/chain/tezos/atomex_tez_token"
	"github.com/atomex-protocol/watch_tower/internal/logger"
	"github.com/dipdup-net/go-lib/node"
	"github.com/dipdup-net/go-lib/tools/forge"
	"github.com/dipdup-net/go-lib/tools/tezgen"
	"github.com/dipdup-net/go-lib/tzkt/api"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/shopspring/decimal"

	"github.com/goat-systems/go-tezos/v4/keys"
)

// Tezos -
type Tezos struct {
	cfg       Config
	rpc       *node.NodeRPC
	api       *api.API
	key       *keys.Key
	counter   int64
	ttl       string
	minPayoff decimal.Decimal

	tezContract   *atomextez.Atomextez
	tokenContract map[string]*atomexteztoken.Atomexteztoken

	log zerolog.Logger

	events     chan chain.Event
	operations chan chain.Operation

	wg sync.WaitGroup
}

// Config -
type Config struct {
	Node            string
	TzKT            string
	MinPayOff       string
	Contract        string
	Tokens          []string
	LogLevel        zerolog.Level
	TTL             int64
	OperaitonParams OperationParamsByContracts
}

// New -
func New(cfg Config) (*Tezos, error) {
	if cfg.LogLevel == 0 {
		cfg.LogLevel = zerolog.InfoLevel
	}
	if len(cfg.OperaitonParams) == 0 {
		return nil, errors.New("empty operations params for tezos. you have to create tezos.yml")
	}

	secret, err := chain.LoadSecret("TEZOS_PRIVATE")
	if err != nil {
		return nil, err
	}

	key, err := keys.FromBase58(secret, keys.Ed25519)
	if err != nil {
		return nil, err
	}

	if cfg.TTL < 1 {
		cfg.TTL = 5
	}

	tokens := make(map[string]*atomexteztoken.Atomexteztoken)
	for i := range cfg.Tokens {
		tokenContract := atomexteztoken.New(cfg.TzKT)
		tokenContract.ChangeAddress(cfg.Tokens[i])
		tokens[cfg.Tokens[i]] = tokenContract
	}

	tez := &Tezos{
		cfg:           cfg,
		rpc:           node.NewNodeRPC(cfg.Node),
		key:           key,
		api:           api.New(cfg.TzKT),
		tezContract:   atomextez.New(cfg.TzKT),
		tokenContract: tokens,
		ttl:           fmt.Sprintf("%d", cfg.TTL),
		log:           logger.New(logger.WithLogLevel(cfg.LogLevel), logger.WithModuleName("tezos")),
		events:        make(chan chain.Event, 1024*16),
		operations:    make(chan chain.Operation, 1024),
	}

	tez.tezContract.ChangeAddress(cfg.Contract)

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
		Address:   t.key.PubKey.GetAddress(),
		PublicKey: t.key.PubKey.GetBytes(),
		Private:   t.key.GetBytes(),
	}
}

// Init -
func (t *Tezos) Init(ctx context.Context) error {
	t.log.Info().Msg("initializing...")

	counterCtx, counterCancel := context.WithTimeout(ctx, 10*time.Second)
	defer counterCancel()
	counterValue, err := t.rpc.Counter(t.key.PubKey.GetAddress(), "head", node.WithContext(counterCtx))
	if err != nil {
		return errors.Wrap(err, "counter")
	}
	counter, err := strconv.ParseInt(counterValue, 10, 64)
	if err != nil {
		return errors.Wrap(err, "invalid counter")
	}
	atomic.StoreInt64(&t.counter, counter)

	return nil
}

// Run -
func (t *Tezos) Run(ctx context.Context) error {
	t.log.Info().Msg("running...")

	t.wg.Add(1)
	go t.listenTezosContract(ctx)

	for _, contract := range t.tokenContract {
		t.wg.Add(1)
		go t.listenTezosTokenContract(ctx, contract)

		if err := contract.Subscribe(ctx); err != nil {
			return err
		}
	}

	return t.tezContract.Subscribe(ctx)
}

// Close -
func (t *Tezos) Close() error {
	t.log.Info().Msg("closing...")
	t.wg.Wait()

	close(t.events)
	close(t.operations)
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

func (t *Tezos) listenTezosContract(ctx context.Context) {
	defer t.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case update := <-t.tezContract.BigMapUpdates():
			if err := t.parseTezosContractUpdate(ctx, update); err != nil {
				t.log.Err(err).Interface("update", update).Msg("parseTezosContractUpdate")
				continue
			}

		case initiate := <-t.tezContract.InitiateEvents():
			t.operations <- chain.Operation{
				Hash:      initiate.Hash,
				ChainType: chain.ChainTypeTezos,
				Status:    toOperationStatus(initiate.Status),
			}
		case add := <-t.tezContract.AddEvents():
			t.operations <- chain.Operation{
				Hash:      add.Hash,
				ChainType: chain.ChainTypeTezos,
				Status:    toOperationStatus(add.Status),
			}
		case redeem := <-t.tezContract.RedeemEvents():
			t.operations <- chain.Operation{
				Hash:      redeem.Hash,
				ChainType: chain.ChainTypeTezos,
				Status:    toOperationStatus(redeem.Status),
			}
		case refund := <-t.tezContract.RefundEvents():
			t.operations <- chain.Operation{
				Hash:      refund.Hash,
				ChainType: chain.ChainTypeTezos,
				Status:    toOperationStatus(refund.Status),
			}
		}
	}
}

func (t *Tezos) listenTezosTokenContract(ctx context.Context, token *atomexteztoken.Atomexteztoken) {
	defer t.wg.Done()

	for {
		select {
		case <-ctx.Done():
			return

		case update := <-token.BigMap0Updates():
			if err := t.parseTokenContractUpdate(ctx, update); err != nil {
				t.log.Err(err).Interface("update", update).Msg("parseTokenContractUpdate")
				continue
			}
		case initiate := <-token.InitiateEvents():
			t.operations <- chain.Operation{
				Hash:      initiate.Hash,
				ChainType: chain.ChainTypeTezos,
				Status:    toOperationStatus(initiate.Status),
			}
		case redeem := <-token.RedeemEvents():
			t.operations <- chain.Operation{
				Hash:      redeem.Hash,
				ChainType: chain.ChainTypeTezos,
				Status:    toOperationStatus(redeem.Status),
			}
		case refund := <-token.RefundEvents():
			t.operations <- chain.Operation{
				Hash:      refund.Hash,
				ChainType: chain.ChainTypeTezos,
				Status:    toOperationStatus(refund.Status),
			}
		}
	}
}

// Initiate -
func (t *Tezos) Initiate(ctx context.Context, args chain.InitiateArgs) error {
	t.log.Info().Str("hashed_secret", args.HashedSecret.String()).Msg("initiate")

	operationParams, ok := t.cfg.OperaitonParams[args.Contract]
	if !ok {
		return errors.Errorf("can't find operation parameters for %s", args.Contract)
	}

	tx := node.Transaction{
		Source:       t.key.PubKey.GetAddress(),
		StorageLimit: operationParams.StorageLimit.Initiate,
		GasLimit:     operationParams.GasLimit.Initiate,
		Fee:          "1000",
		Destination:  args.Contract,
		Amount:       "0",
		Parameters: &node.Parameters{
			Entrypoint: "initiate",
		},
	}
	hashed, err := args.HashedSecret.Bytes()
	if err != nil {
		return err
	}

	var value []byte
	switch args.Contract {
	case t.cfg.Contract:
		tx.Amount = args.Amount.String()

		value, err = t.tezContract.BuildInitiateParameters(ctx, atomextez.Initiate{
			Participant: tezgen.Address(args.Participant),
			Settings: atomextez.Settings{
				HashedSecret: tezgen.Bytes(hashed),
				RefundTime:   tezgen.NewTimestamp(args.RefundTime),
				Payoff:       tezgen.NewInt(args.PayOff.BigInt().Int64()),
			},
		})

	default:
		contract, ok := t.tokenContract[args.Contract]
		if !ok {
			return errors.Errorf("unknown contract: %s", args.Contract)
		}
		value, err = contract.BuildInitiateParameters(ctx, atomexteztoken.Initiate{
			TokenAddress: tezgen.Address(args.TokenAddress),
			Participant:  tezgen.Address(args.Participant),
			HashedSecret: tezgen.Bytes(hashed),
			RefundTime:   tezgen.NewTimestamp(args.RefundTime),
			PayoffAmount: tezgen.NewInt(args.PayOff.BigInt().Int64()),
			TotalAmount:  tezgen.NewInt(args.Amount.BigInt().Int64()),
		})
	}
	if err != nil {
		return err
	}
	if value == nil {
		return nil
	}
	params := json.RawMessage(value)
	tx.Parameters.Value = &params

	opHash, err := t.sendTransaction(ctx, tx)
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
func (t *Tezos) Redeem(ctx context.Context, hashedSecret, secret chain.Hex, contract string) error {
	t.log.Info().Str("hashed_secret", hashedSecret.String()).Str("contract", contract).Msg("redeem")

	value, err := json.Marshal(map[string]interface{}{
		"bytes": secret,
	})
	if err != nil {
		return err
	}

	operationParams, ok := t.cfg.OperaitonParams[contract]
	if !ok {
		return errors.Errorf("can't find operation parameters for %s", contract)
	}

	params := json.RawMessage(value)
	opHash, err := t.sendTransaction(ctx, node.Transaction{
		Source:       t.key.PubKey.GetAddress(),
		Amount:       "0",
		StorageLimit: operationParams.StorageLimit.Redeem,
		GasLimit:     operationParams.GasLimit.Redeem,
		Fee:          "3000",
		Destination:  contract,
		Parameters: &node.Parameters{
			Entrypoint: "redeem",
			Value:      &params,
		},
	})
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
func (t *Tezos) Refund(ctx context.Context, hashedSecret chain.Hex, contract string) error {
	t.log.Info().Str("hashed_secret", hashedSecret.String()).Str("contract", contract).Msg("refund")

	value, err := json.Marshal(map[string]interface{}{
		"bytes": hashedSecret,
	})
	if err != nil {
		return err
	}

	operationParams, ok := t.cfg.OperaitonParams[contract]
	if !ok {
		return errors.Errorf("can't find operation parameters for %s", contract)
	}

	params := json.RawMessage(value)
	opHash, err := t.sendTransaction(ctx, node.Transaction{
		Source:       t.key.PubKey.GetAddress(),
		Amount:       "0",
		StorageLimit: operationParams.StorageLimit.Refund,
		GasLimit:     operationParams.GasLimit.Refund,
		Fee:          "10000",
		Destination:  contract,
		Parameters: &node.Parameters{
			Entrypoint: "refund",
			Value:      &params,
		},
	})
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
func (t *Tezos) Restore(ctx context.Context) error {
	t.log.Info().Msg("restoring...")

	getBigmapsCtx, getBigmapsCancel := context.WithTimeout(ctx, 10*time.Second)
	defer getBigmapsCancel()

	bigMaps, err := t.api.GetBigmaps(getBigmapsCtx, map[string]string{
		"contract.in": strings.Join(append(t.cfg.Tokens, t.cfg.Contract), ","),
	})
	if err != nil {
		return err
	}

	for _, bm := range bigMaps {
		if err := t.restoreFromBigMap(ctx, bm); err != nil {
			return err
		}
	}

	t.events <- chain.RestoredEvent{Chain: chain.ChainTypeTezos}
	return nil
}

func (t *Tezos) restoreFromBigMap(ctx context.Context, bm api.BigMap) error {
	limit := 100
	var end bool
	var offset int
	for !end {
		getBigmapKeysCtx, getBigmapKeysCancel := context.WithTimeout(ctx, 10*time.Second)
		defer getBigmapKeysCancel()

		keys, err := t.api.GetBigmapKeys(getBigmapKeysCtx, uint64(bm.Ptr), map[string]string{
			"limit":  fmt.Sprintf("%d", limit),
			"offset": fmt.Sprintf("%d", offset),
		})
		if err != nil {
			return err
		}

		for i := range keys {
			if err := t.restoreFinilizationSwap(ctx, bm, keys[i]); err != nil {
				return err
			}
		}

		end = len(keys) != limit
		offset += len(keys)
	}
	return nil
}

func (t *Tezos) restoreFinilizationSwap(ctx context.Context, bm api.BigMap, key api.BigMapKey) error {
	updates, err := t.api.GetBigmapKeyUpdates(ctx, uint64(bm.Ptr), key.Key, nil)
	if err != nil {
		return err
	}
	decodedKey, err := hex.DecodeString(key.Key)
	if err != nil {
		return err
	}

	for i := range updates {
		if t.cfg.Contract == bm.Contract.Address {
			var bmUpdate atomextez.BigMapUpdate
			if err := json.Unmarshal(updates[i].Value, &bmUpdate.BigMap.Value); err != nil {
				return err
			}
			bmUpdate.BigMap.Key = atomextez.KeyBigMap(decodedKey)
			bmUpdate.BigMap.Ptr = &bm.Ptr
			bmUpdate.Action = updates[i].Action
			bmUpdate.Contract = bm.Contract.Address
			bmUpdate.Level = updates[i].Level
			if err := t.parseTezosContractUpdate(ctx, bmUpdate); err != nil {
				return err
			}
		} else {
			var bmUpdate atomexteztoken.BigMap0Update
			if err := json.Unmarshal(updates[i].Value, &bmUpdate.BigMap0.Value); err != nil {
				return err
			}
			bmUpdate.BigMap0.Key = atomexteztoken.Key0(decodedKey)
			bmUpdate.BigMap0.Ptr = &bm.Ptr
			bmUpdate.Action = updates[i].Action
			bmUpdate.Contract = bm.Contract.Address
			bmUpdate.Level = updates[i].Level
			if err := t.parseTokenContractUpdate(ctx, bmUpdate); err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *Tezos) parseTezosContractUpdate(ctx context.Context, update atomextez.BigMapUpdate) error {
	hashedSecret := chain.NewHexFromBytes(update.BigMap.Key)

	switch update.Action {
	case BigMapActionAddKey:
		event := chain.InitEvent{
			HashedSecretHex: hashedSecret,
			Chain:           chain.ChainTypeTezos,
			ContractAddress: update.Contract,
			BlockNumber:     update.Level,
			Initiator:       string(update.BigMap.Value.Recipients.Initiator),
			Participant:     string(update.BigMap.Value.Recipients.Participant),
			RefundTime:      update.BigMap.Value.Settings.RefundTime.Value(),
			Amount:          decimal.NewFromBigInt(update.BigMap.Value.Settings.Amount.Int, 0),
		}

		if err := event.SetPayOff(update.BigMap.Value.Settings.Payoff.Int, t.minPayoff); err != nil {
			if errors.Is(err, chain.ErrMinPayoff) {
				t.log.Warn().Str("hashed_secret", event.HashedSecretHex.String()).Msg("skip because of small pay off")
				return nil
			}
			return err
		}

		t.events <- event
	case BigMapActionUpdateKey:
	case BigMapActionRemoveKey:
		ops, err := t.api.GetTransactions(ctx, map[string]string{
			"level":  fmt.Sprintf("%d", update.Level),
			"target": update.Contract,
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
			case atomextez.EntrypointRedeem:
				var secret atomexteztoken.Redeem
				if err := json.Unmarshal(ops[i].Parameters.Value, &secret); err != nil {
					return err
				}

				t.events <- chain.RedeemEvent{
					HashedSecretHex: hashedSecret,
					Chain:           chain.ChainTypeTezos,
					ContractAddress: update.Contract,
					BlockNumber:     update.Level,
					Secret:          chain.NewHexFromBytes(secret),
				}
				return nil
			case atomextez.EntrypointRefund:
				t.events <- chain.RefundEvent{
					HashedSecretHex: hashedSecret,
					Chain:           chain.ChainTypeTezos,
					ContractAddress: update.Contract,
					BlockNumber:     update.Level,
				}
				return nil
			}
		}
	}

	return nil
}

func (t *Tezos) parseTokenContractUpdate(ctx context.Context, update atomexteztoken.BigMap0Update) error {
	hashedSecret := chain.NewHexFromBytes(update.BigMap0.Key)

	switch update.Action {
	case BigMapActionAddKey:
		event := chain.InitEvent{
			HashedSecretHex: hashedSecret,
			Chain:           chain.ChainTypeTezos,
			ContractAddress: update.Contract,
			BlockNumber:     update.Level,
			Initiator:       string(update.BigMap0.Value.Initiator),
			Participant:     string(update.BigMap0.Value.Participant),
			RefundTime:      update.BigMap0.Value.RefundTime.Value(),
			Amount:          decimal.NewFromBigInt(update.BigMap0.Value.TotalAmount.Int, 0),
		}

		if err := event.SetPayOff(update.BigMap0.Value.PayoffAmount.Int, t.minPayoff); err != nil {
			if errors.Is(err, chain.ErrMinPayoff) {
				t.log.Warn().Str("hashed_secret", event.HashedSecretHex.String()).Msg("skip because of small pay off")
				return nil
			}
			return err
		}

		t.events <- event
	case BigMapActionUpdateKey:
	case BigMapActionRemoveKey:
		ops, err := t.api.GetTransactions(ctx, map[string]string{
			"level":  fmt.Sprintf("%d", update.Level),
			"target": update.Contract,
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
			case atomexteztoken.EntrypointRedeem:
				var secret atomexteztoken.Redeem
				if err := json.Unmarshal(ops[i].Parameters.Value, &secret); err != nil {
					return err
				}

				t.events <- chain.RedeemEvent{
					HashedSecretHex: hashedSecret,
					Chain:           chain.ChainTypeTezos,
					ContractAddress: update.Contract,
					BlockNumber:     update.Level,
					Secret:          chain.NewHexFromBytes(secret),
				}
				return nil
			case atomexteztoken.EntrypointRefund:
				t.events <- chain.RefundEvent{
					HashedSecretHex: hashedSecret,
					Chain:           chain.ChainTypeTezos,
					ContractAddress: update.Contract,
					BlockNumber:     update.Level,
				}
				return nil
			}
		}
	}

	return nil
}

func (t *Tezos) sendTransaction(ctx context.Context, transaction node.Transaction) (string, error) {
	headerCtx, headerCancel := context.WithTimeout(ctx, 10*time.Second)
	defer headerCancel()
	header, err := t.rpc.Header(fmt.Sprintf("head~%s", t.ttl), node.WithContext(headerCtx))
	if err != nil {
		return "", err
	}

	atomic.AddInt64(&t.counter, 1)
	transaction.Counter = fmt.Sprintf("%d", t.counter)

	encoded, err := forge.OPG(header.Hash, node.Operation{
		Body: transaction,
		Kind: node.KindTransaction,
	})
	if err != nil {
		return "", err
	}

	msg := hex.EncodeToString(encoded)
	signature, err := t.key.SignHex(msg)
	if err != nil {
		return "", err
	}

	injectCtx, injectCancel := context.WithTimeout(ctx, 10*time.Second)
	defer injectCancel()
	return t.rpc.InjectOperaiton(node.InjectOperationRequest{
		Operation: signature.AppendToHex(msg),
		ChainID:   header.ChainID,
	}, node.WithContext(injectCtx))
}
