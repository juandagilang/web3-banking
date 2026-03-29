package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port            string
	Env             string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	RPCURL          string
	TokenAddress    string
	BankAddress     string
	ChainID         int64
	JWTSecret       string
	JWTExpiryHours  int
	StartBlock      int64
	PollIntervalSec int
}

func Load(configPath string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(configPath)
	v.SetConfigType("env")

	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	cfg := &Config{
		Port:            v.GetString("PORT"),
		Env:             v.GetString("ENV"),
		DBHost:          v.GetString("DB_HOST"),
		DBPort:          v.GetString("DB_PORT"),
		DBUser:          v.GetString("DB_USER"),
		DBPassword:      v.GetString("DB_PASSWORD"),
		DBName:          v.GetString("DB_NAME"),
		DBSSLMode:       v.GetString("DB_SSLMODE"),
		RPCURL:          v.GetString("RPC_URL"),
		TokenAddress:    v.GetString("TOKEN_ADDRESS"),
		BankAddress:     v.GetString("BANK_ADDRESS"),
		ChainID:         v.GetInt64("CHAIN_ID"),
		JWTSecret:       v.GetString("JWT_SECRET"),
		JWTExpiryHours:  v.GetInt("JWT_EXPIRY_HOURS"),
		StartBlock:      v.GetInt64("START_BLOCK"),
		PollIntervalSec: v.GetInt("POLL_INTERVAL_SECONDS"),
	}

	return cfg, nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName, c.DBSSLMode,
	)
}
