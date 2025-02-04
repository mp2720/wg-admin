package config

import "github.com/caarlos0/env/v11"

type Config struct {
	SQLiteDBPath string `env:"SQLITE_DB_PATH"`
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
