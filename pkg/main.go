package main

import (
	"context"
	"fmt"
	"keyless-auth/api"
	"keyless-auth/repository"
	"keyless-auth/storage"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

var (
	redisHost     string
	redisPort     string
	redisPassword string
	redisUsername string
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func init() {
	redisHost = getEnvOrDefault("REDIS_HOST", "localhost")
	redisPort = getEnvOrDefault("REDIS_PORT", "6379")
	redisPassword = getEnvOrDefault("REDIS_PASSWORD", "")
	redisUsername = getEnvOrDefault("REDIS_USERNAME", "")
}

func main() {
	db, err := storage.NewRedisClient(context.Background(), &redis.Options{
		Username: redisUsername,
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
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

	serverAddr := fmt.Sprintf(":%s", getEnvOrDefault("APP_PORT", "8080"))
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
