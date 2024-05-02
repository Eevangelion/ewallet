package config

import (
	"log"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Address string
	Port    int
}

type DBConfig struct {
	MaxOpenConns int
	MinOpenConns int
	ConnTimeout  time.Duration
	Name         string
	Host         string
	Port         int
	User         string
	Password     string
}

type Config struct {
	Server *ServerConfig
	DB     *DBConfig
}

var Conf *Config = nil

func init() {
	projectName := regexp.MustCompile(`^(.*ewallet)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}
}

func GetConfig() *Config {
	if Conf == nil {
		Conf = &Config{
			Server: &ServerConfig{
				Address: getEnv("SERVER", "localhost"),
				Port:    getEnvAsInt("PORT", 8000),
			},
			DB: &DBConfig{
				MaxOpenConns: 4,
				MinOpenConns: 0,
				ConnTimeout:  time.Second * 5,
				Name:         getEnv("DB_NAME", "wallet"),
				Host:         getEnv("DB_HOST", "localhost"),
				Port:         getEnvAsInt("DB_PORT", 5432),
				User:         getEnv("DB_USER", "postgres"),
				Password:     getEnv("DB_PASSWORD", "root"),
			},
		}
	}
	return Conf
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if valueStr, exists := os.LookupEnv(key); exists {
		value, err := strconv.Atoi(valueStr)
		if err == nil {
			return value
		}
	}
	return defaultValue
}
