package ethereum

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/atomex-protocol/watch_tower/internal/chain"
	"github.com/atomex-protocol/watch_tower/internal/config"
	"github.com/ethereum/go-ethereum"
	abi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

// Ethereum -
type Ethereum struct {
	cfg        config.Ethereum
	client     *ethclient.Client
	wss        *ethclient.Client
	subLogs    ethereum.Subscription
	subHead    ethereum.Subscription
	privateKey *ecdsa.PrivateKey

	eth       *AtomexEth
	erc20     *AtomexErc20
	chainID   *big.Int
	minPayoff *big.Int

	latest int64

	logs       chan types.Log
	head       chan *types.Header
	events     chan chain.Event
	operations chan chain.Operation
	stop       chan struct{}
	wg         sync.WaitGroup
}

// New -
func New(cfg config.Ethereum) (*Ethereum, error) {
	client, err := ethclient.Dial(cfg.Node)
	if err != nil {
		return nil, err
	}
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return nil, err
	}

	wss, err := ethclient.Dial(cfg.Wss)
	if err != nil {
		return nil, err
	}

	atomexEth, err := NewAtomexEth(common.HexToAddress(cfg.EthAddress), client)
	if err != nil {
		return nil, err
	}
	atomexErc20, err := NewAtomexErc20(common.HexToAddress(cfg.Erc20Address), client)
	if err != nil {
		return nil, err
	}

	minPayoff := big.NewInt(0)
	if _, ok := minPayoff.SetString(cfg.MinPayOff, 10); !ok {
		return nil, errors.Errorf("invalid minimal payoff value: %s", cfg.MinPayOff)
	}

	eth := Ethereum{
		cfg:        cfg,
		minPayoff:  minPayoff,
		client:     client,
		chainID:    chainID,
		eth:        atomexEth,
		erc20:      atomexErc20,
		wss:        wss,
		logs:       make(chan types.Log, 1024),
		head:       make(chan *types.Header, 16),
		events:     make(chan chain.Event, 1024),
		operations: make(chan chain.Operation, 1024),
		stop:       make(chan struct{}, 1),
	}

	if err := initKeystore(cfg, &eth); err != nil {
		return nil, err
	}

	return &eth, nil
}

func initKeystore(cfg config.Ethereum, e *Ethereum) error {
	secret, err := chain.LoadSecret("ETHEREUM_PRIVATE")
	if err != nil {
		return err
	}
	privateKey, err := crypto.HexToECDSA(secret)
	if err != nil {
		return err
	}
	e.privateKey = privateKey

	log.Info().Str("blockchain", chain.ChainTypeEthereum.String()).Str("address", cfg.UserAddress).Msg("using address")

	return nil
}

func (e *Ethereum) info() *zerolog.Event {
	return log.Info().Str("blockchain", chain.ChainTypeEthereum.String())
}

func (e *Ethereum) warn() *zerolog.Event {
	return log.Warn().Str("blockchain", chain.ChainTypeEthereum.String())
}

func (e *Ethereum) err() *zerolog.Event {
	return log.Error().Str("blockchain", chain.ChainTypeEthereum.String())
}

// Init -
func (e *Ethereum) Init() error {
	e.info().Msg("initializing...")
	return LoadAbi()
}

// Run -
func (e *Ethereum) Run() error {
	e.info().Msg("running...")
	latest, err := e.client.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	e.latest = int64(latest)

	if err := e.subscribe(); err != nil {
		return errors.Wrap(err, "subscribe")
	}

	e.wg.Add(1)
	go e.listen()

	return nil
}

func (e *Ethereum) subscribe() error {
	query := ethereum.FilterQuery{
		Addresses: []common.Address{
			common.HexToAddress(e.cfg.EthAddress),
			common.HexToAddress(e.cfg.Erc20Address),
		},
		FromBlock: big.NewInt(e.latest),
	}

	subLogs, err := e.wss.SubscribeFilterLogs(context.Background(), query, e.logs)
	if err != nil {
		return err
	}
	e.subLogs = subLogs

	subHead, err := e.wss.SubscribeNewHead(context.Background(), e.head)
	if err != nil {
		return err
	}
	e.subHead = subHead
	return nil
}

func (e *Ethereum) reconnect() error {
	wss, err := ethclient.Dial(e.cfg.Wss)
	if err != nil {
		return errors.Wrap(err, "reconnect Dial")
	}
	e.wss = wss
	if err := e.subscribe(); err != nil {
		return errors.Wrap(err, "reconnect subscribe")
	}
	return nil
}

// Close -
func (e *Ethereum) Close() error {
	e.info().Msg("closing...")
	e.stop <- struct{}{}
	e.wg.Wait()

	e.subLogs.Unsubscribe()
	e.subHead.Unsubscribe()
	e.wss.Close()
	e.client.Close()

	close(e.logs)
	close(e.head)
	close(e.events)
	close(e.operations)
	close(e.stop)
	return nil
}

// InitEvents -
func (e *Ethereum) Events() <-chan chain.Event {
	return e.events
}

// Operations -
func (e *Ethereum) Operations() <-chan chain.Operation {
	return e.operations
}

// Redeem -
func (e *Ethereum) Redeem(hashedSecret, secret chain.Hex, contract string) error {
	e.info().Str("hashed_secret", hashedSecret.String()).Str("contract", contract).Msg("redeem")

	opts, err := e.buildTxOpts()
	if err != nil {
		return err
	}

	hashedSecretBytes, err := hashedSecret.Bytes32()
	if err != nil {
		return err
	}
	secretBytes, err := secret.Bytes32()
	if err != nil {
		return err
	}

	var tx *types.Transaction

	switch contract {
	case e.cfg.EthAddress:
		tx, err = e.eth.Redeem(opts, hashedSecretBytes, secretBytes)
		if err != nil {
			return err
		}
	case e.cfg.Erc20Address:
		tx, err = e.erc20.Redeem(opts, hashedSecretBytes, secretBytes)
		if err != nil {
			return err
		}
	}

	e.operations <- chain.Operation{
		Status:       chain.Pending,
		Hash:         tx.Hash().Hex(),
		ChainType:    chain.ChainTypeEthereum,
		HashedSecret: hashedSecret,
	}
	return nil
}

// Refund -
func (e *Ethereum) Refund(hashedSecret chain.Hex, contract string) error {
	e.info().Str("hashed_secret", hashedSecret.String()).Str("contract", contract).Msg("refund")

	opts, err := e.buildTxOpts()
	if err != nil {
		return err
	}

	hashedSecretBytes, err := hashedSecret.Bytes32()
	if err != nil {
		return err
	}
	var tx *types.Transaction

	switch contract {
	case e.cfg.EthAddress:
		tx, err = e.eth.Refund(opts, hashedSecretBytes)
		if err != nil {
			return err
		}
	case e.cfg.Erc20Address:
		tx, err = e.erc20.Refund(opts, hashedSecretBytes)
		if err != nil {
			return err
		}
	}

	e.operations <- chain.Operation{
		Status:       chain.Pending,
		Hash:         tx.Hash().Hex(),
		ChainType:    chain.ChainTypeEthereum,
		HashedSecret: hashedSecret,
	}
	return nil
}

func (e *Ethereum) buildTxOpts() (*bind.TransactOpts, error) {
	auth, err := bind.NewKeyedTransactorWithChainID(e.privateKey, e.chainID)
	if err != nil {
		return nil, err
	}

	gasPrice, err := e.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}

	nonce, err := e.client.PendingNonceAt(context.Background(), common.HexToAddress(e.cfg.UserAddress))
	if err != nil {
		return nil, err
	}

	auth.Nonce = big.NewInt(int64(nonce))

	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)
	auth.GasPrice = gasPrice
	return auth, nil
}

// Restore -
func (e *Ethereum) Restore() error {
	e.info().Msg("restoring...")
	ethEvents, err := e.restoreEth()
	if err != nil {
		return err
	}
	erc20Events, err := e.restoreErc20()
	if err != nil {
		return err
	}
	events := append(ethEvents, erc20Events...)
	sort.Sort(chain.ByLevel(events))

	for i := range events {
		e.events <- events[i]
	}
	e.events <- chain.RestoredEvent{Chain: chain.ChainTypeEthereum}
	return nil
}

func (e *Ethereum) restoreEth() ([]chain.Event, error) {
	events := make([]chain.Event, 0)
	iterInit, err := e.eth.FilterInitiated(nil, nil, nil)
	if err != nil {
		return nil, err
	}
	for iterInit.Next() {
		events = append(events, chain.InitEvent{
			HashedSecretHex: chain.NewHexFromBytes32(iterInit.Event.HashedSecret),
			ContractAddress: e.cfg.EthAddress,
			Chain:           chain.ChainTypeEthereum,
			BlockNumber:     iterInit.Event.Raw.BlockNumber,
			Participant:     iterInit.Event.Participant.Hex(),
			Initiator:       iterInit.Event.Initiator.Hex(),
			Amount:          decimal.NewFromBigInt(iterInit.Event.Value, 0),
			PayOff:          decimal.NewFromBigInt(iterInit.Event.Payoff, 0),
			RefundTime:      time.Unix(iterInit.Event.RefundTimestamp.Int64(), 0),
		})
	}
	if err := iterInit.Close(); err != nil {
		return nil, err
	}

	iterRedeemed, err := e.eth.FilterRedeemed(nil, nil)
	if err != nil {
		return nil, err
	}
	for iterRedeemed.Next() {
		events = append(events, chain.RedeemEvent{
			HashedSecretHex: chain.NewHexFromBytes32(iterRedeemed.Event.HashedSecret),
			ContractAddress: e.cfg.EthAddress,
			Chain:           chain.ChainTypeEthereum,
			BlockNumber:     iterRedeemed.Event.Raw.BlockNumber,
			Secret:          chain.NewHexFromBytes32(iterRedeemed.Event.Secret),
		})
	}
	if err := iterRedeemed.Close(); err != nil {
		return nil, err
	}

	iterRefunded, err := e.eth.FilterRefunded(nil, nil)
	if err != nil {
		return nil, err
	}
	for iterRefunded.Next() {
		events = append(events, chain.RefundEvent{
			HashedSecretHex: chain.NewHexFromBytes32(iterRefunded.Event.HashedSecret),
			ContractAddress: e.cfg.EthAddress,
			Chain:           chain.ChainTypeEthereum,
			BlockNumber:     iterRefunded.Event.Raw.BlockNumber,
		})
	}
	if err := iterRefunded.Close(); err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Ethereum) restoreErc20() ([]chain.Event, error) {
	events := make([]chain.Event, 0)

	iterInit, err := e.erc20.FilterInitiated(nil, nil, nil, nil)
	if err != nil {
		return nil, err
	}
	for iterInit.Next() {
		events = append(events, chain.InitEvent{
			HashedSecretHex: chain.NewHexFromBytes32(iterInit.Event.HashedSecret),
			ContractAddress: e.cfg.Erc20Address,
			Chain:           chain.ChainTypeEthereum,
			BlockNumber:     iterInit.Event.Raw.BlockNumber,
			Participant:     iterInit.Event.Participant.Hex(),
			Initiator:       iterInit.Event.Initiator.Hex(),
			Amount:          decimal.NewFromBigInt(iterInit.Event.Value, 0),
			PayOff:          decimal.NewFromBigInt(iterInit.Event.Payoff, 0),
			RefundTime:      time.Unix(iterInit.Event.RefundTimestamp.Int64(), 0),
		})
	}
	if err := iterInit.Close(); err != nil {
		return nil, err
	}

	iterRedeemed, err := e.erc20.FilterRedeemed(nil, nil)
	if err != nil {
		return nil, err
	}
	for iterRedeemed.Next() {
		events = append(events, chain.RedeemEvent{
			HashedSecretHex: chain.NewHexFromBytes32(iterRedeemed.Event.HashedSecret),
			ContractAddress: e.cfg.Erc20Address,
			Chain:           chain.ChainTypeEthereum,
			BlockNumber:     iterRedeemed.Event.Raw.BlockNumber,
			Secret:          chain.NewHexFromBytes32(iterRedeemed.Event.Secret),
		})
	}
	if err := iterRedeemed.Close(); err != nil {
		return nil, err
	}

	iterRefunded, err := e.erc20.FilterRefunded(nil, nil)
	if err != nil {
		return nil, err
	}
	for iterRefunded.Next() {
		events = append(events, chain.RefundEvent{
			HashedSecretHex: chain.NewHexFromBytes32(iterRefunded.Event.HashedSecret),
			ContractAddress: e.cfg.Erc20Address,
			Chain:           chain.ChainTypeEthereum,
			BlockNumber:     iterRefunded.Event.Raw.BlockNumber,
		})
	}
	if err := iterRefunded.Close(); err != nil {
		return nil, err
	}

	return events, nil
}

func (e *Ethereum) parseLog(l types.Log) error {
	switch l.Address.Hex() {
	case e.cfg.EthAddress:
		return e.parseLogForContract(abiAtomexEth, ContractTypeEth, l)
	case e.cfg.Erc20Address:
		return e.parseLogForContract(abiAtomexErc20, ContractTypeErc20, l)
	}

	return nil
}

func (e *Ethereum) parseLogForContract(abi *abi.ABI, typ string, l types.Log) error {
	if len(l.Topics) == 0 {
		return nil
	}

	event, err := abi.EventByID(l.Topics[0])
	if err != nil {
		return err
	}

	switch event.Name {
	case EventActivated:
	case EventAdded:
	case EventInitiated:
		return e.handleInitiated(abi, l, event, typ)
	case EventRedeemed:
		return e.handleRedeemed(abi, l, event)
	case EventRefunded:
		return e.handleRefunded(l)
	default:
		return errors.Errorf("unknown event: %s", event.Name)
	}
	return nil
}

func (e *Ethereum) listen() {
	defer e.wg.Done()

	for {
		select {
		case <-e.stop:
			return
		case l := <-e.logs:
			if err := e.parseLog(l); err != nil {
				e.err().Err(err).Msg("")
			}
		case head := <-e.head:
			if err := e.parseHead(head); err != nil {
				e.err().Err(err).Msg("")
			}
		case err := <-e.subLogs.Err():
			e.err().Err(err).Msg("ethereum subscription error")
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				if err := e.reconnect(); err != nil {
					e.err().Err(err).Msg("")
				}
			}
		case err := <-e.subHead.Err():
			e.err().Err(err).Msg("ethereum subscription error")
			if websocket.IsCloseError(err, websocket.CloseAbnormalClosure) {
				if err := e.reconnect(); err != nil {
					e.err().Err(err).Msg("")
				}
			}
		}
	}
}

func (e *Ethereum) handleInitiated(abi *abi.ABI, l types.Log, event *abi.Event, typ string) error {
	switch typ {
	case ContractTypeEth:
		if len(l.Topics) != 3 {
			return nil
		}
		var args initiatedEthArgs
		if err := abi.UnpackIntoInterface(&args, event.Name, l.Data); err != nil {
			return err
		}
		hashedSecret := chain.Hex(l.Topics[1].Hex()[2:])
		if e.minPayoff.Cmp(args.PayOff) > 0 {
			e.warn().Str("hashed_secret", hashedSecret.String()).Msg("skip because of small pay off")
			return nil
		}

		e.events <- chain.InitEvent{
			HashedSecretHex: hashedSecret,
			ContractAddress: e.cfg.EthAddress,
			Chain:           chain.ChainTypeEthereum,
			BlockNumber:     l.BlockNumber,
			Participant:     common.BytesToAddress(l.Topics[2].Bytes()).Hex(),
			Initiator:       args.Initiator.Hex(),
			Amount:          decimal.NewFromBigInt(args.Value, 0),
			PayOff:          decimal.NewFromBigInt(args.PayOff, 0),
			RefundTime:      time.Unix(args.RefundTimestamp.Int64(), 0),
		}
	case ContractTypeErc20:
		if len(l.Topics) != 4 {
			return nil
		}
		var args initiatedErc20Args
		if err := abi.UnpackIntoInterface(&args, event.Name, l.Data); err != nil {
			return err
		}

		hashedSecret := chain.Hex(l.Topics[1].Hex()[2:])
		if e.minPayoff.Cmp(args.PayOff) > 0 {
			e.warn().Str("hashed_secret", hashedSecret.String()).Msg("skip because of small pay off")
			return nil
		}

		e.events <- chain.InitEvent{
			HashedSecretHex: chain.Hex(l.Topics[1].Hex()[2:]),
			ContractAddress: e.cfg.Erc20Address,
			Chain:           chain.ChainTypeEthereum,
			BlockNumber:     l.BlockNumber,
			Participant:     l.Topics[3].Hex(),
			Initiator:       args.Initiator.Hex(),
			Amount:          decimal.NewFromBigInt(args.Value, 0),
			PayOff:          decimal.NewFromBigInt(args.PayOff, 0),
			RefundTime:      time.Unix(args.RefundTimestamp.Int64(), 0),
		}
	}
	return nil
}

func (e *Ethereum) handleRefunded(l types.Log) error {
	e.events <- chain.RefundEvent{
		HashedSecretHex: chain.Hex(l.Topics[1].Hex()[2:]),
		Chain:           chain.ChainTypeEthereum,
		ContractAddress: l.Address.Hex(),
		BlockNumber:     l.BlockNumber,
	}
	return nil
}

func (e *Ethereum) handleRedeemed(abi *abi.ABI, l types.Log, event *abi.Event) error {
	var args redeemedArgs
	if err := abi.UnpackIntoInterface(&args, event.Name, l.Data); err != nil {
		return err
	}

	e.events <- chain.RedeemEvent{
		HashedSecretHex: chain.Hex(l.Topics[1].Hex()[2:]),
		Chain:           chain.ChainTypeEthereum,
		ContractAddress: l.Address.Hex(),
		BlockNumber:     l.BlockNumber,
		Secret:          chain.NewHexFromBytes32(args.Secret),
	}
	return nil
}

func (e *Ethereum) parseHead(head *types.Header) error {
	block, err := e.client.BlockByHash(context.Background(), head.Hash())
	if err != nil {
		return err
	}

	txs := block.Transactions()
	if len(txs) == 0 {
		return nil
	}

	for i := range txs {
		to := txs[i].To()
		if to == nil {
			continue
		}
		address := to.Hex()
		if address != e.cfg.EthAddress && address != e.cfg.Erc20Address {
			continue
		}
		receipt, err := e.client.TransactionReceipt(context.Background(), txs[i].Hash())
		if err != nil {
			return err
		}

		e.operations <- chain.Operation{
			ChainType: chain.ChainTypeEthereum,
			Hash:      txs[i].Hash().Hex(),
			Status:    toOperationStatus(receipt.Status),
		}
	}
	return nil
}

func toOperationStatus(status uint64) chain.OperationStatus {
	switch status {
	case 0:
		return chain.Failed
	case 1:
		return chain.Applied
	default:
		return chain.Pending
	}
}
