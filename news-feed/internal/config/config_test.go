package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	cfg := LoadConfig("../../config.json")

	if len(cfg.RSS) != 3 {
		t.Fatal("Не загрузились данные RSS из конфига")
	}

	if cfg.RequestPeriod != 5 {
		t.Fatal("Не загрузился Request Period")
	}
}
