// Package config contains application configuration.
package config

import (
	"strings"

	"github.com/play/notifications/bot"

	"github.com/spf13/viper"
)

const (
	botTgbotapiKey = "tgBOT.TgbotapiKey"
)

type Config struct {
	v *viper.Viper
}

func NewConfig(configName, configPath string) (*Config, error) {
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return &Config{v: v}, nil
}

func (vcfg *Config) NewNotificationBotConfig() *bot.Config {
	return &bot.Config{
		TgbotapiKey: vcfg.v.GetString(botTgbotapiKey),
	}
}
