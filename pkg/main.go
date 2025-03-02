package main

import (
	"context"
	"keyless-auth/api"
	"keyless-auth/repository"
	"keyless-auth/storage"
	"log"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func main() {
	db, err := storage.NewRedisClient(context.Background(), &redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	walletRepo := repository.NewWalletRepository(db)
	credentialsRepo := repository.NewCredentialsRepository(db)
	proofHandler := api.NewProofHandler(walletRepo)
	credentialsHandler := api.NewCredentialsHandler(credentialsRepo, walletRepo)

	router := mux.NewRouter()

	router.HandleFunc("/credentials/{credential}", credentialsHandler.GetWalletAddressByCredential).Methods("GET")
	router.HandleFunc("/credentials", credentialsHandler.GenerateCredential).Methods("POST")
	router.HandleFunc("/proof", proofHandler.GenerateProof).Methods("POST")
}
