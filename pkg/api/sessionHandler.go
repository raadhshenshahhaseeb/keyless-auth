package api

import (
	"crypto/ecdsa"
	"sync"
	"time"
)

// Session holds ephemeral information for a single challenge–response flow.
type Session struct {
	ID               string            // SharedKey of session's shared key
	EphemeralKey     *ecdsa.PrivateKey // Round key
	Challenge        string            // challenge
	CreatedAt        time.Time
	TTL              time.Duration
	Verified         bool             // exists
	ChallengerPubKey *ecdsa.PublicKey // pubK of challenger
}

// IsExpired checks if the session has outlived its TTL.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.CreatedAt.Add(s.TTL))
}

// SessionStore manages ephemeral sessions by their ID (the ephemeral public key).
type SessionStore struct {
	mu         sync.RWMutex
	sessions   map[string]*Session
	defaultTTL time.Duration
}

// NewSessionStore creates a SessionStore with a default TTL for sessions.
func NewSessionStore(defaultTTL time.Duration) *SessionStore {
	if defaultTTL <= 0 {
		defaultTTL = 10 * time.Minute
	}
	return &SessionStore{
		sessions:   make(map[string]*Session),
		defaultTTL: defaultTTL,
	}
}

// CreateSession inserts a new session keyed by the ephemeral public key (ID).
// If session.TTL is zero, it uses defaultTTL. Also sets CreatedAt if zero.
func (ss *SessionStore) CreateSession(sess *Session) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if sess.TTL == 0 {
		sess.TTL = ss.defaultTTL
	}
	if sess.CreatedAt.IsZero() {
		sess.CreatedAt = time.Now()
	}
	ss.sessions[sess.ID] = sess
}

// GetSession returns the session with the given ephemeral public key, or nil if
// no valid session exists (not found or expired).
func (ss *SessionStore) GetSession(ephemeralPubKey string) *Session {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	sess, exists := ss.sessions[ephemeralPubKey]
	if !exists {
		return nil
	}
	if sess.IsExpired() {
		return nil
	}
	return sess
}

func (ss *SessionStore) GetSessionByChallengerKey(key *ecdsa.PublicKey) *Session {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	for _, sess := range ss.sessions {
		if sess.ChallengerPubKey == key {
			return sess
		}
	}
	return nil
}

// DeleteSession removes the session for the given ephemeral public key.
func (ss *SessionStore) DeleteSession(ephemeralPubKey string) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	delete(ss.sessions, ephemeralPubKey)
}

// PruneExpired iterates over sessions and deletes any that are expired.
func (ss *SessionStore) PruneExpired() {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	for id, sess := range ss.sessions {
		if sess.IsExpired() {
			delete(ss.sessions, id)
		}
	}
}

// Count returns the number of currently valid (non-expired) sessions.
func (ss *SessionStore) Count() int {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	count := 0
	for _, sess := range ss.sessions {
		if !sess.IsExpired() {
			count++
		}
	}
	return count
}
