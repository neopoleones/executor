package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log/slog"
	"os"
	"sync"
)

var cfg *Configuration
var once sync.Once

const (
	EnvConfPathKey  = "CFG_PATH"
	DefaultConfPath = "./etc/conf.yml"
)

func envConfPath() string {
	return os.Getenv(EnvConfPathKey)
}

type Configuration struct {
	Service  ServiceConfiguration  `yaml:"service"`
	Database DatabaseConfiguration `yaml:"database"`
	Executor ExecutorConfiguration `yaml:"executor"`
}

type ServiceConfiguration struct {
	Addr string `yaml:"addr"`
}

type DatabaseConfiguration struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"DB"`
}

type ExecutorConfiguration struct {
	InterpreterPath string `yaml:"interpreter_path"`
}

func parseConfiguration(path string) error {
	fd, err := os.Open(path)
	if err != nil {
		return err
	}

	rawCfg := make([]byte, 0, 1024)
	if _, err := fd.Read(rawCfg); err != nil {
		return err
	}

	return yaml.Unmarshal(rawCfg, cfg)
}

func GetConfiguration(possibleLocations ...string) *Configuration {
	once.Do(func() {
		cfg = new(Configuration)

		// I sure possibleLocations isn't that large, so those slow operations
		// won't affect the service

		possibleLocations = append([]string{EnvConfPathKey}, possibleLocations...) // first priority
		possibleLocations = append(possibleLocations, DefaultConfPath)             // last priority

		for _, candidatePath := range possibleLocations {
			if _, err := os.Stat(candidatePath); err == nil {
				// Try to parse configuration with this path
				if err := parseConfiguration(candidatePath); err == nil {
					slog.Info(
						"loaded the configuration",
						slog.String("path", candidatePath),
					)
					return
				}
			}
		}

		// If we got here, there is no good configuration
		panic(fmt.Errorf("failed to load configuration: candidates=%v", possibleLocations))
	})

	return cfg
}
