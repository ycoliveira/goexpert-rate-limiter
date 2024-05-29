package configs

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/viper"
)

type conf struct {
	WebServerPort          string         `mapstructure:"WEB_SERVER_PORT"`
	RateLimiterMaxRequests int            `mapstructure:"RATE_LIMITER_MAX_REQUESTS"`
	BlockTimeSeconds       int            `mapstructure:"BLOCK_TIME_SECONDS"`
	TokenLimits            map[string]int `mapstructure:"-"`
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	// Manually parse the TOKEN_LIMITS
	tokenLimits := viper.GetString("TOKEN_LIMITS")
	var tokenLimitsMap map[string]int
	err = json.Unmarshal([]byte(tokenLimits), &tokenLimitsMap)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TOKEN_LIMITS: %v", err)
	}

	cfg.TokenLimits = tokenLimitsMap
	return cfg, nil
}
