package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/juandagilang/web3bank-backend/internal/infrastructure/blockchain"
	"github.com/juandagilang/web3bank-backend/pkg/response"
)

func HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status":  "healthy",
		"service": "web3bank-backend",
	})
}

func ReadyCheck(db *pgxpool.Pool, bc *blockchain.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		checks := make(map[string]string)
		healthy := true

		if err := db.Ping(c.Request.Context()); err != nil {
			checks["database"] = "unhealthy: " + err.Error()
			healthy = false
		} else {
			checks["database"] = "healthy"
		}

		if bc == nil {
			checks["blockchain"] = "unhealthy: client not initialized"
			healthy = false
		} else {
			checks["blockchain"] = "healthy"
		}

		if healthy {
			response.Success(c, gin.H{
				"status": "ready",
				"checks": checks,
			})
		} else {
			c.JSON(http.StatusServiceUnavailable, response.Response{
				Success: false,
				Data: gin.H{
					"status": "not ready",
					"checks": checks,
				},
			})
		}
	}
}
