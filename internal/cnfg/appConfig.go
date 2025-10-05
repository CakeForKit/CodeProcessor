package cnfg

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Port int `mapstructure:"port"`
}

type RabbitMQConfig struct {
	RabbitMQURL string `mapstructure:"RABBITMQ_URL"`
}

var (
	ErrConfigRead   = errors.New("ReadInConfig")
	ErrUnmarshalKey = errors.New("ErrUnmarshalKey")
)

func LoadAppConfig(path string, fname string, ftype string) (config *AppConfig, err error) {
	config = &AppConfig{}
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(fname)
	v.SetConfigType(ftype)
	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if err = v.UnmarshalKey("app", config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalKey, err)
	}
	return
}

func LoadRabbitMQConfig(path string, fname string, ftype string) (config *RabbitMQConfig, err error) {
	config = &RabbitMQConfig{}
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName(fname)
	v.SetConfigType(ftype)
	if err = v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrConfigRead, err)
	}
	if err = v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnmarshalKey, err)
	}
	return
}
