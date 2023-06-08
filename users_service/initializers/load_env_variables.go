package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const pathEnvVariablesDefault = "./conf/.env"
const conf_file_env_name = "CONF_FILE_PATH"

func LoadEnvVariables() {
	pathEnvVariables := getPathEnvVariables()

	err := godotenv.Load(pathEnvVariables)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getPathEnvVariables() string {
	pathEnvVariables := os.Getenv(conf_file_env_name)
	if pathEnvVariables == "" {
		pathEnvVariables = pathEnvVariablesDefault
	}
	return pathEnvVariables
}
