package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/juandagilang/web3bank-backend/internal/usecase/transaction"
	"github.com/juandagilang/web3bank-backend/pkg/response"
)

type TransactionHandler struct {
	txUC *transaction.UseCase
}

func NewTransactionHandler(txUC *transaction.UseCase) *TransactionHandler {
	return &TransactionHandler{txUC: txUC}
}

func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		address = c.GetString("wallet_address")
	}

	if address == "" {
		response.BadRequest(c, "Address is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.txUC.GetTransactions(c.Request.Context(), address, page, limit)
	if err != nil {
		response.InternalError(c, "Failed to get transactions")
		return
	}

	totalPages := (result.Total + limit - 1) / limit

	response.SuccessPaginated(c, result.Transactions, &response.Pagination{
		Page:       page,
		Limit:      limit,
		Total:      result.Total,
		TotalPages: totalPages,
	})
}
