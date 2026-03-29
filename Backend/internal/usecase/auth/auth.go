package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/juandagilang/web3bank-backend/internal/domain/entity"
	"github.com/juandagilang/web3bank-backend/internal/repository/user"
)

var (
	ErrInvalidSignature = errors.New("invalid signature")
	ErrInvalidAddress   = errors.New("invalid wallet address")
)

type UseCase struct {
	userRepo     *user.PostgresRepository
	jwtSecret    string
	jwtExpiryHrs int
}

func NewUseCase(userRepo *user.PostgresRepository, jwtSecret string, jwtExpiryHrs int) *UseCase {
	return &UseCase{
		userRepo:     userRepo,
		jwtSecret:    jwtSecret,
		jwtExpiryHrs: jwtExpiryHrs,
	}
}

func (uc *UseCase) GetNonce(ctx context.Context, walletAddress string) (string, error) {
	if !isValidAddress(walletAddress) {
		return "", ErrInvalidAddress
	}

	nonce := generateNonce()

	userEntity := &entity.User{
		WalletAddress: walletAddress,
		Nonce:         nonce,
	}

	if err := uc.userRepo.Create(ctx, userEntity); err != nil {
		log.Error().Err(err).Str("wallet", walletAddress).Msg("Failed to create user")
		return "", err
	}

	if err := uc.userRepo.UpdateNonce(ctx, walletAddress, nonce); err != nil {
		log.Error().Err(err).Str("wallet", walletAddress).Msg("Failed to update nonce")
		return "", err
	}

	return nonce, nil
}

func (uc *UseCase) Login(ctx context.Context, walletAddress, signature string) (string, int64, error) {
	if !isValidAddress(walletAddress) {
		return "", 0, ErrInvalidAddress
	}

	userEntity, err := uc.userRepo.GetByWalletAddress(ctx, walletAddress)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return "", 0, ErrInvalidAddress
		}
		log.Error().Err(err).Str("wallet", walletAddress).Msg("Failed to get user")
		return "", 0, err
	}

	message := fmt.Sprintf("Sign this message: %s", userEntity.Nonce)

	if !verifySignature(message, signature, walletAddress) {
		return "", 0, ErrInvalidSignature
	}

	newNonce := generateNonce()
	if err := uc.userRepo.UpdateNonce(ctx, walletAddress, newNonce); err != nil {
		log.Error().Err(err).Str("wallet", walletAddress).Msg("Failed to update nonce after login")
	}

	token, expiry, err := uc.generateToken(walletAddress)
	if err != nil {
		return "", 0, err
	}

	return token, expiry, nil
}

func (uc *UseCase) generateToken(walletAddress string) (string, int64, error) {
	expiry := time.Now().Add(time.Duration(uc.jwtExpiryHrs) * time.Hour)

	claims := jwt.MapClaims{
		"address": walletAddress,
		"exp":     expiry.Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(uc.jwtSecret))
	if err != nil {
		return "", 0, err
	}

	return signedToken, int64(uc.jwtExpiryHrs * 3600), nil
}

func generateNonce() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func isValidAddress(address string) bool {
	return common.IsHexAddress(address)
}

func verifySignature(message, signature, address string) bool {
	signatureBytes := common.FromHex(signature)
	if len(signatureBytes) != 65 {
		return false
	}

	signatureBytes[64] = 27

	messageHash := crypto.Keccak256Hash([]byte(message))

	sigPublicKey, err := crypto.SigToPub(messageHash.Bytes(), signatureBytes)
	if err != nil {
		return false
	}

	recoveredAddress := crypto.PubkeyToAddress(*sigPublicKey)

	return recoveredAddress.Hex() == common.HexToAddress(address).Hex()
}
