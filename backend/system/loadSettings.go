package system

import (
	"encoding/json"
	"os"
)

type Config struct {
	DBPath        string `json:"db_path"`
	JWTSecret     string `env:"SECRET_JWT"`
	JWTTTL        string `env:"JWT_TTL"`
	RefreshSecret string `env:"SECRET_REFRESH"`
	RefreshTTL    string `env:"REFRESH_TTL"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
