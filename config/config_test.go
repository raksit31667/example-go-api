package config

import (
	"testing"
)

type StubEnvGetter map[string]string

func (s StubEnvGetter) Getenv(key string) string {
	return s[key]
}

func TestGetStringEnv(t *testing.T) {
	envGetter := StubEnvGetter{
		"STRING_VAR": "hello",
	}
	configProvider := ConfigProvider{Getter: envGetter}
	t.Run("get value given key exists", func(t *testing.T) {
		got := configProvider.GetStringEnv("STRING_VAR", "")
		want := "hello"

		if got != want {
			t.Errorf("expected %q but got %q", got, want)
		}
	})

	t.Run("get default value given key does not exist", func(t *testing.T) {
		got := configProvider.GetStringEnv("NOT_EXIST", "world")
		want := "world"

		if got != want {
			t.Errorf("expected %q but got %q", got, want)
		}
	})
}

func TestGetIntEnv(t *testing.T) {
	t.Run("get value given key exists", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"INT_VAR": "42",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetIntEnv("INT_VAR", 0)
		want := 42

		if got != want {
			t.Errorf("expected %d but got %d", got, want)
		}
	})

	t.Run("get default value given key does not exist", func(t *testing.T) {
		envGetter := StubEnvGetter{}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetIntEnv("NOT_EXIST", 10)
		want := 10

		if got != want {
			t.Errorf("expected %d but got %d", got, want)
		}
	})

	t.Run("get default value given value cannot be converted to int", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"INT_VAR": "42.5",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetIntEnv("INT_VAR", 10)
		want := 10

		if got != want {
			t.Errorf("expected %d but got %d", got, want)
		}
	})
}

func TestGetBoolEnv(t *testing.T) {
	t.Run("get value given key exists", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"BOOL_VAR": "true",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetBoolEnv("BOOL_VAR", false)
		want := true

		if got != want {
			t.Errorf("expected %v but got %v", got, want)
		}
	})

	t.Run("get default value given key does not exist", func(t *testing.T) {
		envGetter := StubEnvGetter{}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetBoolEnv("NOT_EXIST", false)
		want := false

		if got != want {
			t.Errorf("expected %v but got %v", got, want)
		}
	})

	t.Run("get default value given value cannot be converted to bool", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"BOOL_VAR": "42.5",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		got := configProvider.GetBoolEnv("BOOL_VAR", true)
		want := true

		if got != want {
			t.Errorf("expected %v but got %v", got, want)
		}
	})
}

func TestGetConfig(t *testing.T) {
	t.Run("get server given keys exist", func(t *testing.T) {
		envGetter := StubEnvGetter{
			"HOSTNAME": "127.0.0.1",
			"PORT":     "5000",
		}
		configProvider := ConfigProvider{Getter: envGetter}
		config := configProvider.GetConfig()

		got := config
		want := Config{
			Server{
				Hostname: "127.0.0.1",
				Port:     5000,
			},
		}

		if got != want {
			t.Errorf("expected %v but got %v", got, want)
		}
	})

	t.Run("get server given keys do not exist", func(t *testing.T) {
		envGetter := StubEnvGetter{}
		configProvider := ConfigProvider{Getter: envGetter}
		config := configProvider.GetConfig()

		got := config
		want := Config{
			Server{
				Hostname: "localhost",
				Port:     1323,
			},
		}

		if got != want {
			t.Errorf("expected %v but got %v", got, want)
		}
	})
}
