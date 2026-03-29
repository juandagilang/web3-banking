package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/juandagilang/web3bank-backend/internal/usecase/auth"
	"github.com/juandagilang/web3bank-backend/pkg/response"
)

type AuthHandler struct {
	authUC *auth.UseCase
}

func NewAuthHandler(authUC *auth.UseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

type NonceRequest struct {
	Address string `json:"address" binding:"required"`
}

type NonceResponse struct {
	Nonce   string `json:"nonce"`
	Message string `json:"message"`
}

func (h *AuthHandler) GetNonce(c *gin.Context) {
	var req NonceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: address is required")
		return
	}

	nonce, err := h.authUC.GetNonce(c.Request.Context(), req.Address)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidAddress) {
			response.BadRequest(c, "Invalid wallet address")
			return
		}
		response.InternalError(c, "Failed to generate nonce")
		return
	}

	response.Success(c, NonceResponse{
		Nonce:   nonce,
		Message: "Sign this message: " + nonce,
	})
}

type LoginRequest struct {
	Address   string `json:"address" binding:"required"`
	Signature string `json:"signature" binding:"required"`
}

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: address and signature are required")
		return
	}

	token, expiresIn, err := h.authUC.Login(c.Request.Context(), req.Address, req.Signature)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidAddress) {
			response.Unauthorized(c, "Invalid wallet address")
			return
		}
		if errors.Is(err, auth.ErrInvalidSignature) {
			response.Unauthorized(c, "Invalid signature")
			return
		}
		response.InternalError(c, "Login failed")
		return
	}

	response.Success(c, LoginResponse{
		Token:     token,
		ExpiresIn: expiresIn,
	})
}
