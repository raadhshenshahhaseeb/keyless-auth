package repository

import (
	"context"
	"keyless-auth/storage"
	"log"
)

type CredentialsRepository struct {
	db *storage.Redis
}

func NewCredentialsRepository(db *storage.Redis) *CredentialsRepository {
	return &CredentialsRepository{db: db}
}

func (cred *CredentialsRepository) SaveCredential(credential string) error {
	ctx := context.Background()
	// Add leaf to redis set (for fast membership check)
	if err := cred.db.Client.SAdd(ctx, "merkle:credentials:set", credential).Err(); err != nil {
		log.Printf("Failed to add credential to redis set: %v", err)
		return err
	}

	// Add leaf to redis list (for ordered retrieval)
	if err := cred.db.Client.RPush(ctx, "merkle:credentials:list", credential).Err(); err != nil {
		log.Printf("Failed to add credential to redis list: %v", err)
		return err
	}
	return nil
}

func (cred *CredentialsRepository) DoesCredentialExist(credential string) (bool, error) {
	ctx := context.Background()
	return cred.db.Client.SIsMember(ctx, "merkle:credentials:set", credential).Result()
}

func (cred *CredentialsRepository) GetCredentials() ([]string, error) {
	ctx := context.Background()
	// Get all credentials from redis list to build merkle tree
	creds, err := cred.db.Client.LRange(ctx, "merkle:credentials:list", 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return creds, nil
}
