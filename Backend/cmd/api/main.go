package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/juandagilang/web3bank-backend/internal/config"
	"github.com/juandagilang/web3bank-backend/internal/delivery/handler"
	"github.com/juandagilang/web3bank-backend/internal/delivery/middleware"
	"github.com/juandagilang/web3bank-backend/internal/infrastructure/blockchain"
	"github.com/juandagilang/web3bank-backend/internal/infrastructure/database"
	"github.com/juandagilang/web3bank-backend/internal/repository/transaction"
	"github.com/juandagilang/web3bank-backend/internal/repository/user"
	"github.com/juandagilang/web3bank-backend/internal/usecase/account"
	"github.com/juandagilang/web3bank-backend/internal/usecase/auth"
	"github.com/juandagilang/web3bank-backend/internal/usecase/transaction"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	cfg, err := config.Load(".env")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	db, err := database.NewPostgres(cfg.GetDSN())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	log.Info().Msg("Connected to database")

	userRepo := user.NewPostgresRepository(db)
	txRepo := transaction.NewPostgresRepository(db)

	blockchainClient, err := blockchain.NewClient(cfg.RPCURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to blockchain")
	}
	defer blockchainClient.Close()

	log.Info().Msg("Connected to blockchain")

	eventListener := blockchain.NewEventListener(
		blockchainClient,
		cfg.BankAddress,
		cfg.TokenAddress,
		txRepo,
		cfg.StartBlock,
		cfg.PollIntervalSec,
	)
	go eventListener.Start()

	authUC := auth.NewUseCase(userRepo, cfg.JWTSecret, cfg.JWTExpiryHours)
	accountUC := account.NewUseCase(blockchainClient, cfg.BankAddress)
	txUC := transaction.NewUseCase(txRepo)

	authHandler := handler.NewAuthHandler(authUC)
	accountHandler := handler.NewAccountHandler(accountUC)
	txHandler := handler.NewTransactionHandler(txUC)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	router.Use(middleware.CORS())

	router.GET("/health", handler.HealthCheck)
	router.GET("/ready", handler.ReadyCheck(db, blockchainClient))

	api := router.Group("/api/v1")
	{
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/nonce", authHandler.GetNonce)
			authGroup.POST("/login", authHandler.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWTSecret))
		{
			protected.GET("/account/:address", accountHandler.GetBalance)
			protected.GET("/transactions", txHandler.GetTransactions)
			protected.GET("/contract/info", accountHandler.GetContractInfo)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Info().Str("port", cfg.Port).Msg("Starting server")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	eventListener.Stop()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
