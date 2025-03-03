package repository

import (
	"context"
	"keyless-auth/domain"
	"keyless-auth/storage"
	"time"
)

type WalletRepository struct {
	db *storage.Redis
}

func NewWalletRepository(db *storage.Redis) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) Save(address string, privKey []byte, credential string) error {
	wallet := &domain.Wallet{
		Address:    address,
		PrivateKey: privKey,
		Credential: credential,
	}
	return r.db.Save(context.Background(), storage.GenerateCacheKey("wallet", credential), wallet, time.Hour*24)
}

func (r *WalletRepository) GetWalletByCredential(hashedCredential string) (*domain.Wallet, error) {
	value, err := r.db.Get(context.Background(), storage.GenerateCacheKey("wallet", hashedCredential))
	if err != nil {
		return nil, err
	}
	var wallet domain.Wallet
	err = storage.Deserialize(string(value), &wallet)
	return &wallet, err
}
