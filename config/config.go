package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		println("No .env file found")
	} else {
		println("env loaded")
	}

}

//type Config struct {
//	RedirectUrl  string
//	ClientId     string
//	ClientSecret string
//	NumWorkers   int
//}
//
//func New() *Config {
//	return &Config{
//		RedirectUrl:  getEnv("REDIRECT_URL"),
//		ClientId:     getEnv("CLIENT_ID"),
//		ClientSecret: getEnv("CLIENT_SECRET"),
//		NumWorkers:   getEnvAsInt("NUMBER_OF_WORKERS", 5),
//	}
//}

func GetRedirectUrl() string {
	return getEnv("REDIRECT_URL")
}
func GetClientId() string {
	return getEnv("CLIENT_ID")
}
func GetClientSecret() string {
	return getEnv("CLIENT_SECRET")
}
func GetNumWorkers() int {
	return getEnvAsInt("NUMBER_OF_WORKERS", 5)
}
func GetNumThreads() int {
	return getEnvAsInt("NUMBER_OF_THREADS", 2)
}

func GetPostgresUrl() string {
	return getEnv("POSTGRES_URL")
}
func GetNewRelic() string {
	return getEnv("NEWRELIC")
}

func GetJWTSecret() string {
	return getEnv("JWT_SECRET")
}

func getEnv(key string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	println("Error with environment variables \n Make sure they are available.\n Key is: ", key)

	//need to change this to let main handle the problem and exit if necessary
	os.Exit(1)
	return ""
}

func getEnvAsInt(name string, defaultVal int) int {
	valStr := getEnv(name)
	if value, err := strconv.Atoi(valStr); err != nil {
		return value
	}
	return defaultVal
}
