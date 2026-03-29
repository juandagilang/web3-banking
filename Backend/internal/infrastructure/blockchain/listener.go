package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/rs/zerolog/log"

	"github.com/juandagilang/web3bank-backend/internal/domain/entity"
	"github.com/juandagilang/web3bank-backend/internal/repository/transaction"
)

type Event struct {
	Type        string
	From        common.Address
	To          common.Address
	Amount      *big.Int
	BlockNumber uint64
	Timestamp   uint64
	TxHash      string
}

type EventListener struct {
	client       *Client
	bankAddress  string
	txRepo       transaction.Repository
	startBlock   int64
	pollInterval int
	stopChan     chan struct{}
	wg           sync.WaitGroup
	mu           sync.Mutex
	lastBlock    int64
}

func NewEventListener(
	client *Client,
	bankAddress string,
	txRepo transaction.Repository,
	startBlock int64,
	pollInterval int,
) *EventListener {
	return &EventListener{
		client:       client,
		bankAddress:  bankAddress,
		txRepo:       txRepo,
		startBlock:   startBlock,
		pollInterval: pollInterval,
		stopChan:     make(chan struct{}),
		lastBlock:    startBlock,
	}
}

func (el *EventListener) Start() {
	el.wg.Add(1)
	go el.poll()
}

func (el *EventListener) Stop() {
	close(el.stopChan)
	el.wg.Wait()
	log.Info().Msg("Event listener stopped")
}

func (el *EventListener) poll() {
	defer el.wg.Done()

	ticker := time.NewTicker(time.Duration(el.pollInterval) * time.Second)
	defer ticker.Stop()

	log.Info().Int64("start_block", el.lastBlock).Msg("Event listener started")

	for {
		select {
		case <-el.stopChan:
			return
		case <-ticker.C:
			el.processEvents()
		}
	}
}

func (el *EventListener) processEvents() {
	el.mu.Lock()
	defer el.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	currentBlock, err := el.client.GetBlockNumber(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get current block")
		return
	}

	if currentBlock <= uint64(el.lastBlock) {
		return
	}

	address := common.HexToAddress(el.bankAddress)
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(uint64(el.lastBlock + 1)),
		ToBlock:   new(big.Int).SetUint64(currentBlock),
		Addresses: []common.Address{address},
	}

	logs, err := el.client.Client.FilterLogs(ctx, query)
	if err != nil {
		log.Error().Err(err).Msg("Failed to filter logs")
		return
	}

	bankABI, _ := abi.JSON([]byte(BankABI))

	for _, vLog := range logs {
		event, err := el.parseLog(vLog, bankABI)
		if err != nil {
			log.Warn().Err(err).Str("tx", vLog.TxHash.Hex()).Msg("Failed to parse log")
			continue
		}
		if event == nil {
			continue
		}

		event.BlockNumber = vLog.BlockNumber
		event.TxHash = vLog.TxHash.Hex()

		tx := &entity.Transaction{
			TxHash:         event.TxHash,
			EventType:      entity.TransactionType(event.Type),
			FromAddress:    event.From.Hex(),
			ToAddress:      event.To.Hex(),
			Amount:         event.Amount.String(),
			BlockNumber:    int64(event.BlockNumber),
			BlockTimestamp: int64(event.Timestamp),
		}

		if err := el.txRepo.Create(ctx, tx); err != nil {
			log.Error().Err(err).Str("tx", event.TxHash).Msg("Failed to save transaction")
			continue
		}

		log.Info().
			Str("type", event.Type).
			Str("tx", event.TxHash).
			Str("from", event.From.Hex()).
			Str("amount", event.Amount.String()).
			Msg("Transaction saved")
	}

	el.lastBlock = int64(currentBlock)
}

func (el *EventListener) parseLog(vLog ethereum.Log, bankABI abi.ABI) (*Event, error) {
	for _, event := range bankABI.Events {
		if event.ID == vLog.Topics[0] {
			switch event.Name {
			case "Deposit":
				return parseDepositEvent(vLog)
			case "Withdrawal":
				return parseWithdrawalEvent(vLog)
			case "Transfer":
				return parseTransferEvent(vLog)
			}
		}
	}
	return nil, nil
}

func parseDepositEvent(vLog ethereum.Log) (*Event, error) {
	if len(vLog.Topics) < 2 || len(vLog.Data) < 64 {
		return nil, fmt.Errorf("invalid deposit event data")
	}

	user := common.HexToAddress(vLog.Topics[1].Hex())
	amount := math.ParseBig256OrNil(common.Bytes2Hex(vLog.Data[:32]))
	timestamp := math.ParseBig256OrNil(common.Bytes2Hex(vLog.Data[32:64]))

	return &Event{
		Type:        "deposit",
		From:        user,
		Amount:      amount,
		Timestamp:   timestamp.Uint64(),
		BlockNumber: vLog.BlockNumber,
		TxHash:      vLog.TxHash.Hex(),
	}, nil
}

func parseWithdrawalEvent(vLog ethereum.Log) (*Event, error) {
	if len(vLog.Topics) < 2 || len(vLog.Data) < 64 {
		return nil, fmt.Errorf("invalid withdrawal event data")
	}

	user := common.HexToAddress(vLog.Topics[1].Hex())
	amount := math.ParseBig256OrNil(common.Bytes2Hex(vLog.Data[:32]))
	timestamp := math.ParseBig256OrNil(common.Bytes2Hex(vLog.Data[32:64]))

	return &Event{
		Type:        "withdrawal",
		From:        user,
		Amount:      amount,
		Timestamp:   timestamp.Uint64(),
		BlockNumber: vLog.BlockNumber,
		TxHash:      vLog.TxHash.Hex(),
	}, nil
}

func parseTransferEvent(vLog ethereum.Log) (*Event, error) {
	if len(vLog.Topics) < 3 || len(vLog.Data) < 64 {
		return nil, fmt.Errorf("invalid transfer event data")
	}

	from := common.HexToAddress(vLog.Topics[1].Hex())
	to := common.HexToAddress(vLog.Topics[2].Hex())
	amount := math.ParseBig256OrNil(common.Bytes2Hex(vLog.Data[:32]))
	timestamp := math.ParseBig256OrNil(common.Bytes2Hex(vLog.Data[32:64]))

	return &Event{
		Type:        "transfer",
		From:        from,
		To:          to,
		Amount:      amount,
		Timestamp:   timestamp.Uint64(),
		BlockNumber: vLog.BlockNumber,
		TxHash:      vLog.TxHash.Hex(),
	}, nil
}
