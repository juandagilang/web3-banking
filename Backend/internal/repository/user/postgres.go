package user

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/juandagilang/web3bank-backend/internal/domain/entity"
)

var ErrUserNotFound = errors.New("user not found")

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, user *entity.User) error {
	query := `
		INSERT INTO users (wallet_address, nonce, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		ON CONFLICT (wallet_address) DO NOTHING
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.WalletAddress, user.Nonce).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error().Err(err).Str("wallet", user.WalletAddress).Msg("Failed to create user")
		return err
	}

	return nil
}

func (r *PostgresRepository) GetByWalletAddress(ctx context.Context, walletAddress string) (*entity.User, error) {
	query := `
		SELECT id, wallet_address, nonce, created_at, updated_at
		FROM users
		WHERE wallet_address = $1
	`

	user := &entity.User{}
	err := r.db.QueryRow(ctx, query, walletAddress).Scan(
		&user.ID,
		&user.WalletAddress,
		&user.Nonce,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrUserNotFound
	}

	if err != nil {
		log.Error().Err(err).Str("wallet", walletAddress).Msg("Failed to get user")
		return nil, err
	}

	return user, nil
}

func (r *PostgresRepository) UpdateNonce(ctx context.Context, walletAddress, nonce string) error {
	query := `
		UPDATE users
		SET nonce = $2, updated_at = NOW()
		WHERE wallet_address = $1
	`

	result, err := r.db.Exec(ctx, query, walletAddress, nonce)
	if err != nil {
		log.Error().Err(err).Str("wallet", walletAddress).Msg("Failed to update nonce")
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
