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
	DBKindLocal    = "local"
	DBKindPostgres = "postgres"
)

const (
	EnvConfPathKey      = "CFG_PATH"
	EnvDatabasePassword = "DB_PASS"
	DefaultConfPath     = "./etc/conf.yml"
)

func envConfPath() string {
	return os.Getenv(EnvConfPathKey)
}

type Configuration struct {
	Version string `yaml:"version"`

	Service  ServiceConfiguration  `yaml:"service"`
	Database DatabaseConfiguration `yaml:"database"`
	Executor ExecutorConfiguration `yaml:"executor"`
}

func (c *Configuration) loadSensitiveFromEnv() {
	if v := os.Getenv(EnvDatabasePassword); len(v) > 0 {
		c.Database.Password = v
	}
}

type ServiceConfiguration struct {
	Addr string `yaml:"addr"`
}

type DatabaseConfiguration struct {
	Kind string `yaml:"kind"`

	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
}

type ExecutorConfiguration struct {
	InterpreterPath string `yaml:"interpreter_path"`
}

func parseConfiguration(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}

	rawCfg, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(rawCfg, cfg)
}

func GetConfiguration(possibleLocations ...string) *Configuration {
	once.Do(func() {
		cfg = new(Configuration)

		// I sure possibleLocations isn't that large, so those slow operations
		// won't affect the service

		possibleLocations = append([]string{envConfPath()}, possibleLocations...) // first priority
		possibleLocations = append(possibleLocations, DefaultConfPath)            // last priority

		for _, candidatePath := range possibleLocations {
			if _, err := os.Stat(candidatePath); err == nil {
				// Try to parse configuration with this path
				if err := parseConfiguration(candidatePath); err == nil {
					cfg.loadSensitiveFromEnv()

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
