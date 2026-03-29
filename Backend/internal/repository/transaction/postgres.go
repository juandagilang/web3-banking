package transaction

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	"github.com/juandagilang/web3bank-backend/internal/domain/entity"
)

var ErrTxNotFound = errors.New("transaction not found")

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(db *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) Create(ctx context.Context, tx *entity.Transaction) error {
	query := `
		INSERT INTO transactions (tx_hash, event_type, from_address, to_address, amount, block_number, block_timestamp)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (tx_hash) DO NOTHING
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		tx.TxHash,
		tx.EventType,
		tx.FromAddress,
		tx.ToAddress,
		tx.Amount,
		tx.BlockNumber,
		tx.BlockTimestamp,
	).Scan(&tx.ID, &tx.CreatedAt)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Error().Err(err).Str("tx", tx.TxHash).Msg("Failed to create transaction")
		return err
	}

	return nil
}

func (r *PostgresRepository) GetByTxHash(ctx context.Context, txHash string) (*entity.Transaction, error) {
	query := `
		SELECT id, tx_hash, event_type, from_address, to_address, amount, block_number, block_timestamp, created_at
		FROM transactions
		WHERE tx_hash = $1
	`

	tx := &entity.Transaction{}
	err := r.db.QueryRow(ctx, query, txHash).Scan(
		&tx.ID,
		&tx.TxHash,
		&tx.EventType,
		&tx.FromAddress,
		&tx.ToAddress,
		&tx.Amount,
		&tx.BlockNumber,
		&tx.BlockTimestamp,
		&tx.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrTxNotFound
	}

	if err != nil {
		log.Error().Err(err).Str("tx", txHash).Msg("Failed to get transaction")
		return nil, err
	}

	return tx, nil
}

func (r *PostgresRepository) GetByAddress(ctx context.Context, address string, page, limit int) ([]entity.Transaction, int, error) {
	offset := (page - 1) * limit

	countQuery := `
		SELECT COUNT(*) FROM transactions
		WHERE from_address = $1 OR to_address = $1
	`

	var total int
	err := r.db.QueryRow(ctx, countQuery, address).Scan(&total)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("Failed to count transactions")
		return nil, 0, err
	}

	query := `
		SELECT id, tx_hash, event_type, from_address, to_address, amount, block_number, block_timestamp, created_at
		FROM transactions
		WHERE from_address = $1 OR to_address = $1
		ORDER BY block_timestamp DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, address, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("Failed to get transactions")
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var tx entity.Transaction
		err := rows.Scan(
			&tx.ID,
			&tx.TxHash,
			&tx.EventType,
			&tx.FromAddress,
			&tx.ToAddress,
			&tx.Amount,
			&tx.BlockNumber,
			&tx.BlockTimestamp,
			&tx.CreatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("Failed to scan transaction")
			continue
		}
		transactions = append(transactions, tx)
	}

	return transactions, total, nil
}

func (r *PostgresRepository) GetLastProcessedBlock(ctx context.Context) (int64, error) {
	query := `SELECT COALESCE(MAX(block_number), 0) FROM transactions`

	var blockNumber int64
	err := r.db.QueryRow(ctx, query).Scan(&blockNumber)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get last processed block")
		return 0, err
	}

	return blockNumber, nil
}

func (r *PostgresRepository) SaveLastProcessedBlock(ctx context.Context, blockNumber int64) error {
	return nil
}
