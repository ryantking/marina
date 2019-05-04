package config

import "time"

// Config holds all the configuration variables for Marina
type Config struct {
	RootURL     string `mapstructure:"root_url"`
	Environment string `mapstructure:"environment"`
	Version     string `mapstructure:"version"`

	Server struct {
		Address      string        `mapstructure:"address"`
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	} `mapstructure:"server"`

	DB struct {
		DSN           string        `mapstructure:"dsn"`
		Type          string        `mapstructure:"type"`
		Timeout       time.Duration `mapstructure:"timeout"`
		MaxConns      int           `mapstructure:"max_conns"`
		MigrationsDir string        `mapstructure:"migrations_dir"`
	} `mapstructure:"db"`

	Swagger struct {
		URL      string `mapstructure:"url"`
		DocsPath string `mapstructure:"docs_path"`
	} `mapstructure:"swagger"`

	S3 struct {
		Endpoint        string `mapstructure:"endpoint"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	} `mapstructure:"s3"`
}
