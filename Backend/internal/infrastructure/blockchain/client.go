package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

type Client struct {
	Client   *ethclient.Client
	bankABI  abi.ABI
	tokenABI abi.ABI
}

func NewClient(rpcURL string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to blockchain: %w", err)
	}

	bankABI, err := abi.JSON([]byte(BankABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse bank ABI: %w", err)
	}

	tokenABI, err := abi.JSON([]byte(TokenABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse token ABI: %w", err)
	}

	log.Info().Str("rpc", rpcURL).Msg("Connected to blockchain")

	return &Client{
		Client:   client,
		bankABI:  bankABI,
		tokenABI: tokenABI,
	}, nil
}

func (c *Client) Close() {
	c.Client.Close()
}

func (c *Client) GetBlockNumber(ctx context.Context) (uint64, error) {
	return c.Client.BlockNumber(ctx)
}

func (c *Client) GetBalance(ctx context.Context, bankAddress, userAddress string) (*big.Int, error) {
	bank := common.HexToAddress(bankAddress)
	user := common.HexToAddress(userAddress)

	data, err := c.bankABI.Pack("balanceOf", user)
	if err != nil {
		return nil, err
	}

	var result hexBigInt
	err = c.Client.CallContract(ctx, ethereum.CallMsg{
		To:   &bank,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call balanceOf: %w", err)
	}

	return result.ToInt(), nil
}

func (c *Client) GetTokenBalance(ctx context.Context, tokenAddress, userAddress string) (*big.Int, error) {
	token := common.HexToAddress(tokenAddress)
	user := common.HexToAddress(userAddress)

	data, err := c.tokenABI.Pack("balanceOf", user)
	if err != nil {
		return nil, err
	}

	var result hexBigInt
	err = c.Client.CallContract(ctx, ethereum.CallMsg{
		To:   &token,
		Data: data,
	}, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call balanceOf: %w", err)
	}

	return result.ToInt(), nil
}

func (c *Client) GetTokenInfo(ctx context.Context, tokenAddress string) (name, symbol, decimals string, totalSupply *big.Int, err error) {
	token := common.HexToAddress(tokenAddress)

	nameData, err := c.tokenABI.Pack("name")
	if err != nil {
		return "", "", "", nil, err
	}

	var nameResult string
	err = c.Client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: nameData}, nil)
	if err != nil {
		nameResult = "Web3Bank Token"
	} else {
		nameResult = decodeString(nameResult)
	}

	symbolData, err := c.tokenABI.Pack("symbol")
	if err != nil {
		return "", "", "", nil, err
	}

	var symbolResult string
	err = c.Client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: symbolData}, nil)
	if err != nil {
		symbolResult = "W3B"
	} else {
		symbolResult = decodeString(symbolResult)
	}

	decimalsData, err := c.tokenABI.Pack("decimals")
	if err != nil {
		return "", "", "", nil, err
	}

	var decimalsResult uint8 = 18
	err = c.Client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: decimalsData}, nil)
	if err != nil {
		decimalsResult = 18
	}

	totalSupplyData, err := c.tokenABI.Pack("totalSupply")
	if err != nil {
		return "", "", "", nil, err
	}

	var totalResult hexBigInt
	err = c.Client.CallContract(ctx, ethereum.CallMsg{To: &token, Data: totalSupplyData}, nil)
	if err != nil {
		return "", "", "", nil, err
	}

	return nameResult, symbolResult, fmt.Sprintf("%d", decimalsResult), totalResult.ToInt(), nil
}

func (c *Client) FilterEvents(ctx context.Context, contractAddress string, startBlock, endBlock uint64) ([]Event, error) {
	address := common.HexToAddress(contractAddress)

	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(startBlock),
		ToBlock:   new(big.Int).SetUint64(endBlock),
		Addresses: []common.Address{address},
	}

	logs, err := c.Client.FilterLogs(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to filter logs: %w", err)
	}

	var events []Event
	for _, vLog := range logs {
		event, err := c.parseLog(vLog)
		if err != nil {
			log.Warn().Err(err).Str("tx", vLog.TxHash.Hex()).Msg("Failed to parse log")
			continue
		}
		if event != nil {
			events = append(events, *event)
		}
	}

	return events, nil
}

func (c *Client) parseLog(vLog ethereum.Log) (*Event, error) {
	for _, event := range c.bankABI.Events {
		if event.ID == vLog.Topics[0] {
			switch event.Name {
			case "Deposit":
				return parseDeposit(vLog)
			case "Withdrawal":
				return parseWithdrawal(vLog)
			case "Transfer":
				return parseTransfer(vLog)
			}
		}
	}
	return nil, nil
}

type hexBigInt struct {
	big.Int
}

func (h *hexBigInt) UnmarshalJSON(data []byte) error {
	s := string(data)
	if s == "0x" || s == `"0x"` {
		h.SetInt64(0)
		return nil
	}
	h.SetString(s, 0)
	return nil
}

func decodeString(data []byte) string {
	if len(data) > 64 {
		offset := 32
		length := new(big.Int).SetBytes(data[offset : offset+32]).Uint64()
		return string(data[offset+32 : offset+32+int(length)])
	}
	return ""
}
