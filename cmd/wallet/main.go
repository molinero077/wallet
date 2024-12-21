package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"wallet/internal/api/v1/app"
	"wallet/internal/model"
	"wallet/internal/model/postgres"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetLevel(log.DebugLevel)
	filePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	configFile := fmt.Sprintf("%s/%s", path.Dir(filePath), "config.env")
	log.Info("config file will be used: ", configFile)
	err = godotenv.Load(configFile)
	if err != nil {
		log.Fatal("no .env file")
	}

	if net.ParseIP(os.Getenv("API_HOST")) == nil {
		log.Fatal("API_HOST must be a IP address")
	}

	if _, err := strconv.Atoi(os.Getenv("API_PORT")); err != nil {
		log.Fatal("API_PORT must be a number")
	}

	if l, ok := os.LookupEnv("LOG_LEVEL"); ok && strings.ToUpper(l) == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	app := app.New()

	// init storage
	var storage model.WalletProvider
	storage, err = postgres.New(app.GetContext(), &postgres.ConnectionParameters{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Database: os.Getenv("DB_NAME"),
	})
	if err != nil {
		log.Fatal(err)
	}

	app.AssignStorage(&storage)
	app.Run(fmt.Sprintf("%s:%s", os.Getenv("API_HOST"), os.Getenv("API_PORT")))

}
