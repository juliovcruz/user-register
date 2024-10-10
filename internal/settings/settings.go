package settings

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment string

const (
	Local      Environment = "local"
	Staging    Environment = "staging"
	Production Environment = "production"
)

type Settings struct {
	ZipCodeSettings ZipCode
	TokenSecret     string
	Database        Database
}

type ZipCode struct {
	ViaCEPBaseURL string
}

type Database struct {
	FilePath string
	Driver   string
	Secrets  Secrets
}

type Secrets struct {
	Previous string
	Current  string
}

var settingsByEnvironment = map[Environment]Settings{
	Local: {
		ZipCodeSettings: ZipCode{
			ViaCEPBaseURL: "https://viacep.com.br/ws/",
		},
		Database: Database{
			FilePath: "./database.db",
			Driver:   "sqlite3",
		},
	},
	Staging: {
		ZipCodeSettings: ZipCode{
			ViaCEPBaseURL: "https://viacep.com.br/ws/",
		},
		Database: Database{
			FilePath: "./database.db",
			Driver:   "sqlite3",
		},
	},
	Production: {
		ZipCodeSettings: ZipCode{
			ViaCEPBaseURL: "https://viacep.com.br/ws/",
		},
		Database: Database{
			FilePath: "./database.db",
			Driver:   "sqlite3",
		},
	},
}

func LoadSettings(environment Environment) (Settings, error) {
	settings := settingsByEnvironment[environment]

	if err := godotenv.Load("./.env"); err != nil {
		log.Fatalf("Error loading .env file")
	}

	value := os.Getenv("TOKEN_SECRET")
	if value == "" {
		return Settings{}, fmt.Errorf("TOKEN_SECRET is required")
	}
	settings.TokenSecret = value

	value = os.Getenv("DB_CURRENT_SECRET")
	if value == "" {
		return Settings{}, fmt.Errorf("DB_CURRENT_SECRET is required")
	}
	settings.Database.Secrets.Current = value

	value = os.Getenv("DB_PREVIOUS_SECRET")
	if value == "" {
		return Settings{}, fmt.Errorf("DB_PREVIOUS_SECRET is required")
	}
	settings.Database.Secrets.Previous = value

	return settings, nil
}
