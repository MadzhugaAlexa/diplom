package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Config struct {
	RequestPeriod int      `json:"request_period"`
	RSS           []string `json:"rss"`
	Duration      time.Duration
}

// LoadConfig читает конфигурацию из файла config.json и загружает из нее
// список каналов RSS ыи период перезагрузки данных из каналаы
func LoadConfig(path string) Config {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Не смогли открыть config.json")
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)

	cfg := Config{}
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatal("Не смогли загрузить config")
	}

	cfg.Duration = time.Minute
	return cfg
}
