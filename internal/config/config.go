package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	envPath := "../config/.env"

	absPath, err := filepath.Abs(envPath)
	if err != nil {
		log.Fatalf("не удалось определить абсолютный путь до %s: %v", envPath, err)
	}

	_, statErr := os.Stat(envPath)
	if os.IsNotExist(statErr) {
		log.Fatalf(".env файл не найден по пути: %s", absPath)
	} else if statErr != nil {
		log.Fatalf("Ошибка при попытке доступа к .env: %v", statErr)
	} else {
		log.Printf(".env файл найден по пути: %s\n", absPath)
	}

	loadErr := godotenv.Load(envPath)
	if loadErr != nil {
		log.Fatalf("Ошибка при загрузке .env файла (%s): %v", absPath, loadErr)
	}

	log.Println(".env успешно загружен.")
}

func GetDBConnString() string {
	return "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"
}

func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET не задан явно в .env")
	}
	return secret
}
