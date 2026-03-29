package entity

import "time"

type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
)

type Transaction struct {
	ID             uint            `json:"id" gorm:"primaryKey"`
	TxHash         string          `json:"tx_hash" gorm:"uniqueIndex;not null;size:66"`
	EventType      TransactionType `json:"event_type" gorm:"not null;size:20"`
	FromAddress    string          `json:"from_address" gorm:"not null;size:42;index"`
	ToAddress      string          `json:"to_address" gorm:"size:42"`
	Amount         string          `json:"amount" gorm:"not null;size:78"`
	BlockNumber    int64           `json:"block_number" gorm:"not null;index"`
	BlockTimestamp int64           `json:"block_timestamp" gorm:"not null;index"`
	CreatedAt      time.Time       `json:"created_at"`
}

func (Transaction) TableName() string {
	return "transactions"
}

type TransactionResponse struct {
	ID          uint            `json:"id"`
	Type        TransactionType `json:"type"`
	From        string          `json:"from"`
	To          string          `json:"to,omitempty"`
	Amount      string          `json:"amount"`
	BlockNumber int64           `json:"block_number"`
	Timestamp   int64           `json:"timestamp"`
	TxHash      string          `json:"tx_hash"`
}

func (t *Transaction) ToResponse() TransactionResponse {
	return TransactionResponse{
		ID:          t.ID,
		Type:        t.EventType,
		From:        t.FromAddress,
		To:          t.ToAddress,
		Amount:      t.Amount,
		BlockNumber: t.BlockNumber,
		Timestamp:   t.BlockTimestamp,
		TxHash:      t.TxHash,
	}
}
