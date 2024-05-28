package configs

import "github.com/spf13/viper"

type conf struct {
	WebServerPort          string `mapstructure:"WEB_SERVER_PORT"`
	RateLimiterMaxRequests string `mapstructure:"RATE_LIMITER_MAX_REQUESTS"`
	BlockTimeSeconds       int    `mapstructure:"BLOCK_TIME_SECONDS"`
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
	return cfg, nil
}
