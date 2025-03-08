package repository

import (
	"context"
	"keyless-auth/domain"
	"keyless-auth/storage"
)

type UserRepository struct {
	db *storage.Redis
}

func NewUserRepository(db *storage.Redis) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SaveUser(user *domain.User) error {
	ctx := context.Background()
	return r.db.Client.Set(ctx, "user:"+user.ID, user, 0).Err()
}

func (r *UserRepository) GetUser(id string) (*domain.User, error) {
	ctx := context.Background()

	var user domain.User
	if err := r.db.Client.Get(ctx, "user:"+id).Scan(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
