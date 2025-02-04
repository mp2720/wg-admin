package config

import "github.com/caarlos0/env/v11"

type Config struct {
	SQLiteDBPath        string `env:"SQLITE_DB_PATH"`
	AuthTokenSigningKey string `env:"AUTH_TOKEN_SIGNING_KEY"`
	AuthTokenIssuer     string `env:"AUTH_TOKEN_ISSUER" envDefault:"wgadmin"`
}

func Load() (*Config, error) {
	cfg := Config{}
	if err := env.ParseWithOptions(&cfg, env.Options{
		RequiredIfNoDef: true,
		Prefix:          "WGA_",
	}); err != nil {
		return nil, err
	}

	return &cfg, nil
}
