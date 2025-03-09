package repository

import (
	"context"
	"log"
	"time"

	"keyless-auth/domain"
	"keyless-auth/storage"
)

type WalletRepository struct {
	db *storage.Redis
}

func NewWalletRepository(db *storage.Redis) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Save(address string, privKey []byte, credential string, merkleRoot string) error {
	wallet := &domain.Wallet{
		Address:    address,
		PrivateKey: privKey,
		Credential: credential,
		MerkleRoot: merkleRoot,
	}
	log.Printf("Saving wallet: %v", wallet)

	serializedWallet, err := storage.Serialize(wallet)
	if err != nil {
		log.Printf("Failed to serialize wallet: %v", err)
		return err
	}

	err = r.db.Save(context.Background(), storage.GenerateCacheKey("wallet", credential), serializedWallet, time.Hour*24)
	if err != nil {
		log.Printf("Failed to save wallet: %v", err)
		return err
	}
	return nil
}

func (r *WalletRepository) GetWalletByCredential(hashedCredential string) (*domain.Wallet, error) {
	value, err := r.db.Get(context.Background(), storage.GenerateCacheKey("wallet", hashedCredential))
	if err != nil {
		log.Printf("Failed to get wallet by credential: %v", err)
		return nil, err
	}
	var wallet domain.Wallet
	err = storage.Deserialize(string(value), &wallet)
	return &wallet, err
}
