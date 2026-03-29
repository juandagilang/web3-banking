package transaction

import (
	"context"

	"github.com/juandagilang/web3bank-backend/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, tx *entity.Transaction) error
	GetByTxHash(ctx context.Context, txHash string) (*entity.Transaction, error)
	GetByAddress(ctx context.Context, address string, page, limit int) ([]entity.Transaction, int, error)
	GetLastProcessedBlock(ctx context.Context) (int64, error)
	SaveLastProcessedBlock(ctx context.Context, blockNumber int64) error
}
