package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Enviroment string

const (
	Dev  Enviroment = "dev"
	Qc   Enviroment = "qc"
	Prod Enviroment = "prod"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Logger   LoggerConfig
	Database DatabaseConfig
	Data     DataConfig
	Ai       AiConfig
}

func LoadConfig(env Enviroment, configPath string) (*Config, error) {
	if env == "" {
		env = Dev // Default environment
	}

	//Load base config first
	baseConfig, err := loadBaseConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load base config: %w", err)
	}
	// load environment specific config
	envConfig, err := loadEnvConfig(env, configPath)

	if err != nil {
		return nil, fmt.Errorf("failed to load %s config: %w", env, err)
	}

	// Merge configs
	config := mergeConfigs(baseConfig, envConfig)

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	return config, nil
}

func loadBaseConfig(configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName("base")
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

func loadEnvConfig(env Enviroment, configPath string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigName(string(env))
	v.SetConfigType("yaml")
	v.AddConfigPath(configPath)

	// Enable environment variable substitution
	v.SetEnvPrefix("MMO_AFFILIATE_APP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

func mergeConfigs(base, env *viper.Viper) *Config {
	var config Config

	// Merge base config
	if err := base.Unmarshal(&config); err != nil {
		panic("failed to unmarshal base config: " + err.Error())
	}

	// Override with environment specific settings
	if err := env.Unmarshal(&config); err != nil {
		panic("failed to unmarshal env config: " + err.Error())
	}

	return &config
}

func validateConfig(config *Config) error {
	if config.Server.HTTP.Port <= 0 {
		return fmt.Errorf("server port must be greater than 0")
	}
	return nil
}
