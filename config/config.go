package config

import (
	"os"
	"strconv"
)

type EnvGetter interface {
	Getenv(key string) string
}

type OsEnvGetter struct{}

func (o *OsEnvGetter) Getenv(key string) string {
	return os.Getenv(key)
}

type ConfigProvider struct {
	Getter EnvGetter
}

type Config struct {
	Server Server
}

type Server struct {
	Hostname string
	Port     int
}

func (c *ConfigProvider) GetStringEnv(key string, defaultValue string) string {
	value := c.Getter.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func (c *ConfigProvider) GetIntEnv(key string, defaultValue int) int {
	value := c.Getter.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intValue
}

func (c *ConfigProvider) GetBoolEnv(key string, defaultValue bool) bool {
	value := c.Getter.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}

func (c *ConfigProvider) GetConfig() Config {
	return Config{
		Server: Server{
			Hostname: c.GetStringEnv("HOSTNAME", "localhost"),
			Port:     c.GetIntEnv("PORT", 1323),
		},
	}
}
