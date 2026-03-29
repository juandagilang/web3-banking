package user

import (
	"context"

	"github.com/juandagilang/web3bank-backend/internal/domain/entity"
)

type Repository interface {
	Create(ctx context.Context, user *entity.User) error
	GetByWalletAddress(ctx context.Context, walletAddress string) (*entity.User, error)
	UpdateNonce(ctx context.Context, walletAddress, nonce string) error
}
