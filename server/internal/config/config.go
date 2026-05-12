package config

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string
	StoragePath string
	GRPC        GRPCConfig
}

type GRPCConfig struct {
	Port    string
	Timeout time.Duration
}

func Load() (*Config, error) {
	path := fetchConfigPath()

	if path == "" {
		return nil, fmt.Errorf("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", path)
	}

	var rawCfg struct {
		Env         string `yaml:"env" env-default:"local"`
		StoragePath string `yaml:"storage_path" env-required:"true"`
		GRPC        struct {
			Port            string `yaml:"port" env-default:"50051"`
			Timeout_seconds int    `yaml:"timeout_seconds" env-default:"5"`
		} `yaml:"grpc"`
	}

	if err := cleanenv.ReadConfig(path, &rawCfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	cfg := Config{
		Env:         rawCfg.Env,
		StoragePath: rawCfg.StoragePath,
		GRPC: GRPCConfig{
			Port:    rawCfg.GRPC.Port,
			Timeout: time.Duration(rawCfg.GRPC.Timeout_seconds) * time.Second,
		},
	}

	return &cfg, nil
}

func fetchConfigPath() string {
	var res string
	f := flag.Lookup("config")
	if f != nil {
		res = f.Value.String()
	} else {
		flag.StringVar(&res, "config", "", "path to config file")
		flag.Parse()
	}

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
