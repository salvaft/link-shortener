package cfg

import "os"

type Config struct {
	DbName string
	Port   string
	Host   string
}

func InitConfig() Config {
	return Config{
		DbName: getEnv("DB_NAME", "db.sqlite3"),
		Port:   getEnv("PORT", "8000"),
		Host:   getEnv("HOST", "localhost"),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
