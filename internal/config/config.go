package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
)

var Config *config

type config struct {
	DbDSN string
}

func absPath(filename string) string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(b), "../..", filename)
}

func Load(filename string) error {
	filename = absPath(filename)

	if err := godotenv.Load(filename); err != nil {
		return fmt.Errorf("load env: %v", err)
	}

	Config = &config{
		DbDSN: os.Getenv("DB_DSN"),
	}

	return nil
}
