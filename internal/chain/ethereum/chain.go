package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"sync"
	"time"

	"github.com/aopoltorzhicky/watch_tower/internal/chain"
	"github.com/aopoltorzhicky/watch_tower/internal/config"
	"github.com/ethereum/go-ethereum"
	abi "github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/pkg/errors"
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
	initChan   chan chain.InitEvent
	redeemChan chan chain.RedeemEvent
	refundChan chan chain.RefundEvent
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
		initChan:   make(chan chain.InitEvent, 1024),
		redeemChan: make(chan chain.RedeemEvent, 1024),
		refundChan: make(chan chain.RefundEvent, 1024),
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

	log.Info().Str("address", cfg.UserAddress).Str("blockchain", chain.ChainTypeEthereum.String()).Msg("using address")

	return nil
}

// Run -
func (e *Ethereum) Run() error {
	if err := LoadAbi(); err != nil {
		return err
	}

	latest, err := e.client.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	e.latest = int64(latest)

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

	e.wg.Add(1)
	go e.listen()

	return nil
}

// Close -
func (e *Ethereum) Close() error {
	e.stop <- struct{}{}
	e.wg.Wait()

	e.subLogs.Unsubscribe()
	e.subHead.Unsubscribe()
	e.wss.Close()
	e.client.Close()

	close(e.logs)
	close(e.head)
	close(e.initChan)
	close(e.redeemChan)
	close(e.refundChan)
	close(e.operations)
	close(e.stop)
	return nil
}

// InitEvents -
func (e *Ethereum) InitEvents() <-chan chain.InitEvent {
	return e.initChan
}

// RedeemEvents -
func (e *Ethereum) RedeemEvents() <-chan chain.RedeemEvent {
	return e.redeemChan
}

// RefundEvents -
func (e *Ethereum) RefundEvents() <-chan chain.RefundEvent {
	return e.refundChan
}

// Operations -
func (e *Ethereum) Operations() <-chan chain.Operation {
	return e.operations
}

// Redeem -
func (e *Ethereum) Redeem(hashedSecret, secret chain.Hex, contract string) error {
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
	if err := e.restoreEth(); err != nil {
		return err
	}
	if err := e.restoreErc20(); err != nil {
		return err
	}
	return nil
}

func (e *Ethereum) restoreEth() error {
	events := make(map[string]chain.InitEvent)

	iterInit, err := e.eth.FilterInitiated(nil, nil, nil)
	if err != nil {
		return err
	}
	for iterInit.Next() {
		hashedSecret := chain.NewHexFromBytes32(iterInit.Event.HashedSecret)

		if e.minPayoff.Cmp(iterInit.Event.Payoff) > 0 {
			continue
		}

		events[hashedSecret.String()] = chain.InitEvent{
			Event: chain.Event{
				HashedSecret: hashedSecret,
				Contract:     e.cfg.EthAddress,
				Chain:        chain.ChainTypeEthereum,
			},
			Participant: iterInit.Event.Participant.Hex(),
			Initiator:   iterInit.Event.Initiator.Hex(),
			Amount:      decimal.NewFromBigInt(iterInit.Event.Value, 0),
			PayOff:      decimal.NewFromBigInt(iterInit.Event.Payoff, 0),
			RefundTime:  time.Unix(iterInit.Event.RefundTimestamp.Int64(), 0),
		}
	}
	if err := iterInit.Close(); err != nil {
		return err
	}

	iterRedeemed, err := e.eth.FilterRedeemed(nil, nil)
	if err != nil {
		return err
	}
	for iterRedeemed.Next() {
		hashedSecret := hex.EncodeToString(iterRedeemed.Event.HashedSecret[:])
		delete(events, hashedSecret)
	}
	if err := iterRedeemed.Close(); err != nil {
		return err
	}

	iterRefunded, err := e.eth.FilterRefunded(nil, nil)
	if err != nil {
		return err
	}
	for iterRefunded.Next() {
		hashedSecret := hex.EncodeToString(iterRefunded.Event.HashedSecret[:])
		delete(events, hashedSecret)
	}
	if err := iterRefunded.Close(); err != nil {
		return err
	}

	log.Info().Int("count", len(events)).Msg("initiated swaps found in eth")

	for hashedSecret, event := range events {
		e.initChan <- event
		delete(events, hashedSecret)
	}

	return nil
}

func (e *Ethereum) restoreErc20() error {
	events := make(map[string]chain.InitEvent)

	iterInit, err := e.erc20.FilterInitiated(nil, nil, nil, nil)
	if err != nil {
		return err
	}
	for iterInit.Next() {
		hashedSecret := chain.NewHexFromBytes32(iterInit.Event.HashedSecret)

		if e.minPayoff.Cmp(iterInit.Event.Payoff) > 0 {
			continue
		}

		events[hashedSecret.String()] = chain.InitEvent{
			Event: chain.Event{
				HashedSecret: hashedSecret,
				Contract:     e.cfg.Erc20Address,
				Chain:        chain.ChainTypeEthereum,
			},
			Participant: iterInit.Event.Participant.Hex(),
			Initiator:   iterInit.Event.Initiator.Hex(),
			Amount:      decimal.NewFromBigInt(iterInit.Event.Value, 0),
			PayOff:      decimal.NewFromBigInt(iterInit.Event.Payoff, 0),
			RefundTime:  time.Unix(iterInit.Event.RefundTimestamp.Int64(), 0),
		}
	}
	if err := iterInit.Close(); err != nil {
		return err
	}

	iterRedeemed, err := e.erc20.FilterRedeemed(nil, nil)
	if err != nil {
		return err
	}
	for iterRedeemed.Next() {
		hashedSecret := hex.EncodeToString(iterRedeemed.Event.HashedSecret[:])
		delete(events, hashedSecret)
	}
	if err := iterRedeemed.Close(); err != nil {
		return err
	}

	iterRefunded, err := e.erc20.FilterRefunded(nil, nil)
	if err != nil {
		return err
	}
	for iterRefunded.Next() {
		hashedSecret := hex.EncodeToString(iterRefunded.Event.HashedSecret[:])
		delete(events, hashedSecret)
	}
	if err := iterRefunded.Close(); err != nil {
		return err
	}

	log.Info().Int("count", len(events)).Msg("initiated swaps found in erc20")

	for hashedSecret, event := range events {
		e.initChan <- event
		delete(events, hashedSecret)
	}

	return nil
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
				log.Err(err).Msg("")
			}
		case head := <-e.head:
			if err := e.parseHead(head); err != nil {
				log.Err(err).Msg("")
			}
		case err := <-e.subLogs.Err():
			log.Err(err).Msg("ethereum subscription error")
		case err := <-e.subHead.Err():
			log.Err(err).Msg("ethereum subscription error")
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
			log.Warn().Str("hashed_secret", hashedSecret.String()).Msg("skip because of small pay off")
			return nil
		}

		e.initChan <- chain.InitEvent{
			Event: chain.Event{
				HashedSecret: hashedSecret,
				Contract:     e.cfg.EthAddress,
				Chain:        chain.ChainTypeEthereum,
			},
			Participant: common.BytesToAddress(l.Topics[2].Bytes()).Hex(),
			Initiator:   args.Initiator.Hex(),
			Amount:      decimal.NewFromBigInt(args.Value, 0),
			PayOff:      decimal.NewFromBigInt(args.PayOff, 0),
			RefundTime:  time.Unix(args.RefundTimestamp.Int64(), 0),
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
			log.Warn().Str("hashed_secret", hashedSecret.String()).Msg("skip because of small pay off")
			return nil
		}

		e.initChan <- chain.InitEvent{
			Event: chain.Event{
				HashedSecret: chain.Hex(l.Topics[1].Hex()[2:]),
				Contract:     e.cfg.Erc20Address,
				Chain:        chain.ChainTypeEthereum,
			},
			Participant: l.Topics[3].Hex(),
			Initiator:   args.Initiator.Hex(),
			Amount:      decimal.NewFromBigInt(args.Value, 0),
			PayOff:      decimal.NewFromBigInt(args.PayOff, 0),
			RefundTime:  time.Unix(args.RefundTimestamp.Int64(), 0),
		}
	}
	return nil
}

func (e *Ethereum) handleRefunded(l types.Log) error {
	e.refundChan <- chain.RefundEvent{
		Event: chain.Event{
			HashedSecret: chain.Hex(l.Topics[1].Hex()[2:]),
			Chain:        chain.ChainTypeEthereum,
			Contract:     l.Address.Hex(),
		},
	}
	return nil
}

func (e *Ethereum) handleRedeemed(abi *abi.ABI, l types.Log, event *abi.Event) error {
	var args redeemedArgs
	if err := abi.UnpackIntoInterface(&args, event.Name, l.Data); err != nil {
		return err
	}

	e.redeemChan <- chain.RedeemEvent{
		Event: chain.Event{
			HashedSecret: chain.Hex(l.Topics[1].Hex()[2:]),
			Chain:        chain.ChainTypeEthereum,
			Contract:     l.Address.Hex(),
		},
		Secret: chain.NewHexFromBytes32(args.Secret),
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
