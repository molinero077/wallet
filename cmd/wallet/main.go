package main

import (
	"context"
	"log"
	"wallet/internal/model"
	"wallet/internal/model/postgres"
)

func main() {
	var (
		storage model.WalletProvider
		err     error
	)
	ctx := context.Background()
	storage, err = postgres.New(&ctx, &postgres.ConnectionParameters{
		User:     "postgres",
		Password: "postgres",
		Host:     "192.168.88.103",
		Port:     "5432",
	})

	if err != nil {
		log.Fatal(err)
	}

	walletId := "17ccd802-4a56-48ca-abbd-b76986c734d9"

	balance, err := storage.GetBalance(walletId)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Balance for wallet %s is %f", walletId, balance)

}
