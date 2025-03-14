package config

import (
	"log"
	"os"
)

type Config struct {
	ClientID      string
	AppToken      string
	TelegramToken string
	ServerURL     string
	DBHost        string
	User          string
	Password      string
	DBName        string
	DBPort        string
	PayModeOn     bool
}

func NewConfig() *Config {
	host, exists := os.LookupEnv("DBHOST")
	if !exists {
		log.Fatal("env par HOST not found")
	}
	user, exists := os.LookupEnv("USER")
	if !exists {
		log.Fatal("env par HOST not found")
	}
	password, exists := os.LookupEnv("PASSWORD")
	if !exists {
		log.Fatal("env par HOST not found")
	}
	dbname, exists := os.LookupEnv("DBNAME")
	if !exists {
		log.Fatal("env par HOST not found")
	}
	port, exists := os.LookupEnv("DBPORT")
	if !exists {
		log.Fatal("env par HOST not found")
	}

	clientId, exists := os.LookupEnv("CLIENT_ID")
	if !exists {
		log.Fatal("CLIENT_ID env parametr not found!")
	}

	appToken, exists := os.LookupEnv("APP_TOKEN")
	if !exists {
		log.Fatal("APP_TOKEN env parametr not found!")
	}

	server_url, exists := os.LookupEnv("SERVER_URL")
	if !exists {
		log.Fatal("SERVER_URL env parametr not found!")
	}

	botToken, exists := os.LookupEnv("TELEGRAM_API_TOKEN")
	if !exists {
		log.Fatal("Telegram view token not found in .env file")
	}

	conf := &Config{
		ClientID:      clientId,
		AppToken:      appToken,
		TelegramToken: botToken,
		ServerURL:     server_url,
		DBHost:        host,
		User:          user,
		Password:      password,
		DBName:        dbname,
		DBPort:        port,
		PayModeOn:     true,
	}

	_, exists = os.LookupEnv("PAY_MODE")
	if !exists {
		conf.PayModeOn = false
		log.Print("PAY_MODE not found, setting to OFF")
	}

	return conf
}
