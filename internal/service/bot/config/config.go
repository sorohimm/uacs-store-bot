package config

import "context"

type LoggerConf struct {
	Level   string `short:"l" long:"level" env:"LEVEL" description:"logging level" default:"DEBUG"`
	EncType string `long:"enctype" env:"ENCTYPE" description:"log as json or not (console|json)" default:"json" `
}

type Config struct {
	Log   *LoggerConf `group:"logger option" namespace:"log" env-namespace:"LOG"`
	Token string      `long:"token" env:"TOKEN" description:"telegram bot token"`
	Debug bool        `long:"debug" env:"DEBUG" description:"debug mode (default is false)"`
}

type confKey struct{} // or exported to use outside the package

func WithContext(ctx context.Context, c *Config) context.Context {
	return context.WithValue(ctx, confKey{}, c)
}

func FromContext(ctx context.Context) *Config {
	if cc, ok := ctx.Value(confKey{}).(*Config); ok {
		return cc
	}
	return NewDefaultConfig()
}

func NewDefaultConfig() *Config {
	return &Config{}
}
