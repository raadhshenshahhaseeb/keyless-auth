package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"

	"keyless-auth/api"
	"keyless-auth/repository"
	"keyless-auth/repository/user"
	"keyless-auth/service/signer"
	"keyless-auth/services"
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
	db, err := services.NewRedisClient(context.Background(), &redis.Options{
		Username: redisUsername,
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	if err != nil {
		log.Fatal(err)
	}
	defer db.Client.Close()

	newSigner, _ := signer.New()
	sessionStore := api.NewSessionStore(10 * time.Minute)

	walletRepo := repository.NewWalletRepository(db)
	credentialsRepo := repository.NewCredentialsRepository(db)
	proofHandler := api.NewProofHandler(walletRepo)
	credentialsHandler := api.NewCredentialsHandler(credentialsRepo, walletRepo)
	userRepo := user.NewUser(db)
	googleHandler := api.NewGoogleHandler(userRepo)
	challengeHandler := api.NewChallengeHandler(newSigner, db, userRepo, sessionStore)

	router := mux.NewRouter()

	// credentials
	router.HandleFunc("/credentials", credentialsHandler.GenerateCredential).Methods("POST")
	router.HandleFunc("/merkle-root", credentialsHandler.GetMerkleRoot).Methods("GET")
	router.HandleFunc("/merkle-proof/{credential}", credentialsHandler.GenerateMerkleProof).Methods("GET")
	router.HandleFunc("/generate-tree-object", credentialsHandler.GenerateTreeObject).Methods("POST")
	// zk proof
	router.HandleFunc("/proof", proofHandler.GenerateProof).Methods("POST")
	// auth
	router.HandleFunc("/auth/google/login", googleHandler.HandleGoogleLogin).Methods("GET")
	router.HandleFunc("/auth/google/callback", googleHandler.HandleGoogleCallback).Methods("GET")

	router.HandleFunc("/challenge", challengeHandler.SendChallengeHandler).Methods("POST")
	router.HandleFunc("/challenge/verify", challengeHandler.VerifyChallengeHandler).Methods("POST")
	serverAddr := fmt.Sprintf(":%s", getEnvOrDefault("APP_PORT", "8081"))
	log.Printf("Server starting on %s", serverAddr)
	log.Fatal(http.ListenAndServe(serverAddr, router))
}
