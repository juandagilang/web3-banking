package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/juandagilang/web3bank-backend/internal/usecase/account"
	"github.com/juandagilang/web3bank-backend/pkg/response"
)

type AccountHandler struct {
	accountUC *account.UseCase
}

func NewAccountHandler(accountUC *account.UseCase) *AccountHandler {
	return &AccountHandler{accountUC: accountUC}
}

func (h *AccountHandler) GetBalance(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		response.BadRequest(c, "Address is required")
		return
	}

	info, err := h.accountUC.GetBalance(c.Request.Context(), address)
	if err != nil {
		response.InternalError(c, "Failed to get balance")
		return
	}

	response.Success(c, info)
}

func (h *AccountHandler) GetContractInfo(c *gin.Context) {
	info, err := h.accountUC.GetContractInfo(c.Request.Context(), "")
	if err != nil {
		response.InternalError(c, "Failed to get contract info")
		return
	}

	response.Success(c, info)
}
