package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type (
	Config struct {
		App
		HTTP
		HLF
	}

	App struct {
		AppName    string `env:"APP_NAME" env-default:"saac-v2 client"`
		AppVersion string `env:"APP_VERSION" env-default:"1.0"`
		LogLevel   string `env:"LOG_LEVEL" env-default:"info"`
	}

	HTTP struct {
		BindAddress string `env:"BIND_ADDRESS" env-default:"0.0.0.0"`
		BindPort    uint   `env:"BIND_PORT" env-default:"8082"`
	}

	HLF struct {
		WalletPath              string `env:"WALLET_PATH" env-default:"wallet"`
		ChannelName             string `env:"CHANNEL_NAME" env-default:"testchannel"`
		ChaincodeName           string `env:"CHAINCODE_NAME" env-default:"saacv2"`
		OrganizationsFolderPath string `env:"ORGANIZATION_FOLDER_PATH" env-default:"../../organizations"`
	}
)

func NewWebApiConfig() *Config {
	cfg := &Config{}

	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	return cfg
}
