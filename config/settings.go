package config

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
	ERROR LogLevel = "ERROR"
	PANIC LogLevel = "PANIC"
	FATAL LogLevel = "FATAL"
)

type _Setting struct {
	CORSAllowOrigins []string `json:"CORS_ALLOW_ORIGINS"`
	CORSAllowHeaders []string `json:"CORS_ALLOW_HEADERS"`
	CORSAllowMethods []string `json:"CORS_ALLOW_METHODS"`

	DBPoolSize      int    `json:"DB_POOL_SIZE" default:"10"`
	ConnMaxIdleTime int    `json:"CONN_MAX_IDLE_TIME" default:"30" description:"30 seconds"`
	ConnMaxLifetime int    `json:"CONN_MAX_LIFETIME" default:"30" description:"30 seconds"`
	Debug           bool   `json:"DEBUG"`
	DBHost          string `json:"DB_HOST"`
	DBPort          int    `json:"DB_PORT"`
	DBUser          string `json:"DB_USER"`
	DBPass          string `json:"DB_PASS"`
	DBName          string `json:"DB_NAME"`

	DefaultPageSize   int `json:"DEFAULT_PAGE_SIZE"`
	DefaultPageNumber int `json:"DEFAULT_PAGE_NUMBER"`

	LogLevel LogLevel `json:"LOG_LEVEL"`
}

var Settings *_Setting

func init() {
	// Initialize Settings
	Settings = &_Setting{}
	isNotDebug := os.Getenv("NOT_DEBUG")
	fmt.Printf("isNotDebug: '%s'\n", isNotDebug)

	if isNotDebug != "true" {
		jsonFilePath := os.Getenv("ENV_PATH")
		if jsonFilePath == "" {
			log.Fatal("ENV_PATH is not set")
		}
		jsonFileData, err := os.ReadFile(jsonFilePath)
		if err != nil {
			log.Fatalf("Error reading file: %v", err)
		}
		err = json.Unmarshal(jsonFileData, Settings)
		if err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
		}
	} else {
		secretManagerBase64 := os.Getenv("SECRET_MANAGER_BASE64")
		if secretManagerBase64 == "" {
			log.Fatal("SECRET_MANAGER_BASE64 is not set")
		}

		decodedData, err := base64.StdEncoding.DecodeString(secretManagerBase64)
		if err != nil {
			log.Fatalf("Error decoding SECRET_MANAGER_BASE64: %v", err)
		}

		err = json.Unmarshal(decodedData, Settings)
		if err != nil {
			log.Fatalf("Error unmarshalling SECRET_MANAGER_BASE64: %v", err)
		}
	}

	step := new(string)
	*step = "Starting Settings"
	fmt.Printf("Initializing settings: %s\n", *step)

	defer func(s *string) {
		fmt.Printf("Settings initialization step: %s\n", *s)
	}(step)
	if Settings.LogLevel == "" {
		Settings.LogLevel = INFO
	}

	fmt.Println("================================================")
	fmt.Println("   Finished Settings")
	fmt.Println("================================================")
}
