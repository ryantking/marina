package config

import (
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	rootURL     = "root_url"
	environment = "environment"
	version     = "version"

	serverAddr         = "server.address"
	serverWriteTimeout = "server.write_timeout"
	serverReadTimeout  = "server.read_timeout"

	dbDSN           = "db.dsn"
	dbType          = "db.type"
	dbTimeout       = "db.timeout"
	dbMaxConns      = "db.max_conns"
	dbMigrationsDir = "db.migrations_dir"

	swaggerURL      = "swagger.url"
	swaggerDocsPath = "swagger.docs_path"

	s3Endpoint        = "s3.endpoint"
	s3Bucket          = "s3.bucket"
	s3Region          = "s3.region"
	s3AccessKeyID     = "s3.access_key_id"
	s3SecretAccessKey = "s3.secret_access_key"
	s3UseSSL          = "s3.use_ssl"
)

var (
	cfg  *Config
	cfgL sync.Mutex

	// Version holds the git SHA or tag
	Version string
)

// Get returns the configuration, initializing it if it has not been
func Get() *Config {
	cfgL.Lock()
	defer cfgL.Unlock()
	if cfg != nil {
		return cfg
	}

	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// General configuration
	v.SetDefault(rootURL, "http://localhost:8080")
	v.SetDefault(environment, "DEVELOPMENT")
	v.SetDefault(version, Version)

	// Server configuration
	v.SetDefault(serverAddr, "0.0.0.0:8080")
	v.SetDefault(serverWriteTimeout, time.Minute*5)
	v.SetDefault(serverReadTimeout, time.Minute*5)

	// Database configuration
	v.SetDefault(dbDSN, "marina:marina@tcp(localhost:3306)/marinatest")
	v.SetDefault(dbType, "mysql")
	v.SetDefault(dbTimeout, time.Minute)
	v.SetDefault(dbMaxConns, 200)
	v.SetDefault(dbMigrationsDir, "migrations")

	// Swagger configuration
	v.SetDefault(swaggerURL, "https://swagger.cogolo.net")
	v.SetDefault(swaggerDocsPath, "/api/docs")

	// S3 configuration
	v.SetDefault(s3Endpoint, "localhost:9000")
	v.SetDefault(s3Bucket, "marina")
	v.SetDefault(s3Region, "us-east-1")
	v.SetDefault(s3AccessKeyID, "minio")
	v.SetDefault(s3SecretAccessKey, "minio123")
	v.SetDefault(s3UseSSL, false)
	v.AutomaticEnv()

	c := Config{}
	err := v.Unmarshal(&c)
	if err != nil {
		log.WithError(err).Fatal("error unmarshalling configuration")
	}
	cfg = &c

	return cfg
}

// Destroy deletes the config so on next call it will be reinitialized
func Destroy() {
	cfgL.Lock()
	defer cfgL.Unlock()
	cfg = nil
}

// Set sets a specific config value
func Set(key string, val interface{}) {
	if cfg == nil {
		Get()
	}

	v := viper.New()
	v.Set(key, val)
	cfgL.Lock()
	defer cfgL.Unlock()

	err := v.Unmarshal(cfg)
	if err != nil {
		log.WithError(err).Fatalln("error unmarshalling config")
	}
}
