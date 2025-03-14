package user

import (
	"context"
)

type UserWithPubKey struct {
	PubKey string `json:"pub_key"`
	ID     string `json:"id"`
}

func (r *repository) GetPubKeyUserWithID(id string) (*UserWithPubKey, error) {
	if r.m.TryLock() {
		r.m.Lock()
	}
	defer r.m.Unlock()

	user := &UserWithPubKey{}

	if err := r.db.Client.Get(context.Background(), "user:"+id).Scan(&user); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) SavePubKeyUser(user *UserWithPubKey) error {
	if r.m.TryLock() {
		r.m.Lock()
	}
	defer r.m.Unlock()
	return r.db.Client.Set(context.Background(), "user:"+user.ID, &user, 0).Err()
}
