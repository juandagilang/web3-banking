package account

import (
	"context"

	"github.com/juandagilang/web3bank-backend/internal/infrastructure/blockchain"
)

type UseCase struct {
	blockchainClient *blockchain.Client
	bankAddress      string
}

func NewUseCase(blockchainClient *blockchain.Client, bankAddress string) *UseCase {
	return &UseCase{
		blockchainClient: blockchainClient,
		bankAddress:      bankAddress,
	}
}

type AccountInfo struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Symbol  string `json:"symbol"`
}

func (uc *UseCase) GetBalance(ctx context.Context, walletAddress string) (*AccountInfo, error) {
	balance, err := uc.blockchainClient.GetBalance(ctx, uc.bankAddress, walletAddress)
	if err != nil {
		return nil, err
	}

	return &AccountInfo{
		Address: walletAddress,
		Balance: balance.String(),
		Symbol:  "W3B",
	}, nil
}

type ContractInfo struct {
	TokenAddress string `json:"token_address"`
	BankAddress  string `json:"bank_address"`
	Name         string `json:"name"`
	Symbol       string `json:"symbol"`
	Decimals     string `json:"decimals"`
	TotalSupply  string `json:"total_supply"`
}

func (uc *UseCase) GetContractInfo(ctx context.Context, tokenAddress string) (*ContractInfo, error) {
	name, symbol, decimals, totalSupply, err := uc.blockchainClient.GetTokenInfo(ctx, tokenAddress)
	if err != nil {
		return nil, err
	}

	return &ContractInfo{
		TokenAddress: tokenAddress,
		BankAddress:  uc.bankAddress,
		Name:         name,
		Symbol:       symbol,
		Decimals:     decimals,
		TotalSupply:  totalSupply.String(),
	}, nil
}
