package config

import (
	"encoding/json"
	"errors"
	"os"
)

type ServerKey struct {
	KeyID string `json:"kid"`
	Value string `json:"value"`
}

type Config struct {
	Mode            Mode        `json:"Mode"`
	MySQLDSN        string      `json:"MySQLDSN"`
	MySQLServerName string      `json:"MySQLServerName"`
	DBMaxConns      float64     `json:"DBMaxConns"`
	JwtSecret       []ServerKey `json:"JwtSecret"`
	RedisCluster    []string    `json:"RedisCluster"`
}

var AppConfig Config

func LoadConfig() error {
	// Read config data from config file
	data, err := os.ReadFile("/credentials/config")
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &AppConfig)
	if err != nil {
		return err
	}

	// Validate config data
	switch AppConfig.Mode {
	case ModeLocal:
	case ModeStage:
	case ModeProduction:
	default:
		return errors.New("invalid mode: Mode must be [ local | staging | production ]")
	}

	// Validate JWT key
	if len(AppConfig.JwtSecret) == 0 {
		return errors.New("invalid jwt key")
	}
	if len(AppConfig.JwtSecret[0].KeyID) == 0 || len(AppConfig.JwtSecret[0].Value) == 0 {
		return errors.New("invalid jwt key")
	}

	return nil
}

type Mode string

const (
	ModeLocal      Mode = "local"
	ModeStage      Mode = "stage"
	ModeProduction Mode = "production"
)
