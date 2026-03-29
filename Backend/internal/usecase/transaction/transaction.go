package transaction

import (
	"context"

	"github.com/juandagilang/web3bank-backend/internal/domain/entity"
	"github.com/juandagilang/web3bank-backend/internal/repository/transaction"
)

type UseCase struct {
	txRepo *transaction.PostgresRepository
}

func NewUseCase(txRepo *transaction.PostgresRepository) *UseCase {
	return &UseCase{txRepo: txRepo}
}

type TransactionList struct {
	Transactions []entity.TransactionResponse `json:"transactions"`
	Total        int                          `json:"total"`
}

func (uc *UseCase) GetTransactions(ctx context.Context, walletAddress string, page, limit int) (*TransactionList, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	transactions, total, err := uc.txRepo.GetByAddress(ctx, walletAddress, page, limit)
	if err != nil {
		return nil, err
	}

	responses := make([]entity.TransactionResponse, len(transactions))
	for i, tx := range transactions {
		responses[i] = tx.ToResponse()
	}

	return &TransactionList{
		Transactions: responses,
		Total:        total,
	}, nil
}
