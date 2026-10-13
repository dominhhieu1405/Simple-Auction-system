package config

import (
	"log"
	"os"
)

type Config struct {
	JWTSecret []byte
	DBPath    string
	Addr      string
}

func Load() Config {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Println("WARNING: Thêm JWT_SECRET vào .env đi")
		secret = "abcdefghijklmnopqrstuvwxyz123456"
	}
	return Config{
		JWTSecret: []byte(secret),
		DBPath:    getenv("DB_PATH", "data/auction.db"),
		Addr:      ":" + getenv("SERVER_PORT", "8080"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
