package config

import xpostgres "github.com/TienHien1201/go-mmo-affiliate/pkg/postgres"

type PostgresConfig struct {
	Host            string
	Port            int
	Database        string
	Username        string
	Password        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime int
	Timezone        string
	SSLMode         string
}

type RedisConfig struct {
	Host         string
	Port         int
	Password     string
	DB           int
	PoolSize     int
	PoolTimeout  int
	MinIdleConns int
}

type ElasticSearchConfig struct {
	Addresses     []string
	Username      string
	Password      string
	APIKey        string
	EnableLogging bool
	Header        map[string][]string
}

type DatabaseConfig struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Elastic  ElasticSearchConfig
}

func (c *Config) InitPostgresDB() (*xpostgres.Client, error) {
	return xpostgres.NewClient(&xpostgres.Config{
		Host:            c.Database.Postgres.Host,
		Port:            c.Database.Postgres.Port,
		Username:        c.Database.Postgres.Username,
		Password:        c.Database.Postgres.Password,
		Database:        c.Database.Postgres.Database,
		Timezone:        c.Database.Postgres.Timezone,
		SSLMode:         c.Database.Postgres.SSLMode,
		MaxIdleConns:    c.Database.Postgres.MaxIdleConns,
		MaxOpenConns:    c.Database.Postgres.MaxOpenConns,
		ConnMaxLifetime: c.Database.Postgres.ConnMaxLifetime,
	})
}
