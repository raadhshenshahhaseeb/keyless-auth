package api

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	"keyless-auth/repository/user"
	"keyless-auth/service/signer"
	"keyless-auth/services"
)

// [WIP]
type PubKeyLoginHandler struct {
	SignerObject    signer.Signer
	UserPubKeyStore user.Repo
	db              *services.RedisClient
	ss              *SessionStore
}

func (s *PubKeyLoginHandler) Metamask() string {
	return ""
}

func NewChallengeHandler(s signer.Signer, db *services.RedisClient, u user.Repo, ss *SessionStore) *PubKeyLoginHandler {
	return &PubKeyLoginHandler{
		SignerObject:    s,
		UserPubKeyStore: u,
		db:              db,
		ss:              ss,
	}
}

type VerifyChallengeRequest struct {
	DecryptedChallenge string `json:"decrypted_challenge"`
	Signature          string `json:"signature"`
	EphemeralPK        string `json:"ephemeral_pk"`
	ChallengerPK       string `json:"challenger_pk"`
}

type VerifyChallengeResponse struct {
	ID           string `json:"id"`
	SessionToken string `json:"session_token"`
	PublicKey    string `json:"public_key"`
}

func (s *PubKeyLoginHandler) VerifyChallengeHandler(w http.ResponseWriter, r *http.Request) {
	var req VerifyChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	session := s.ss.GetSession(req.EphemeralPK)
	if session == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userPK, err := s.SignerObject.PublicKeyFromBytes(req.ChallengerPK)
	if err != nil {
		http.Error(w, "Invalid public key", http.StatusBadRequest)
		return
	}

	if userPK != session.ChallengerPubKey || strings.Compare(req.EphemeralPK, session.ID) != 0 {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// Verify the signature using the decrypted challenge as the message.
	verified := s.SignerObject.VerifySignature(*userPK, req.DecryptedChallenge, req.Signature)
	if !verified {
		http.Error(w, "Invalid signature", http.StatusBadRequest)
		return
	}

	newUID := uuid.New().String()
	err = s.UserPubKeyStore.SavePubKeyUser(&user.UserWithPubKey{
		PubKey: string(crypto.FromECDSAPub(userPK)),
		ID:     newUID,
	})
	if err != nil {
		http.Error(w, "Failed to save user", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(VerifyChallengeResponse{
		ID:           newUID,
		SessionToken: "",
	})
}

func (s *PubKeyLoginHandler) SendChallengeHandler(w http.ResponseWriter, r *http.Request) {
	var req ChallengePayloadReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	userPK, err := s.SignerObject.PublicKeyFromBytes(req.PubKey)
	if err != nil {
		http.Error(w, "invalid public key", http.StatusBadRequest)
	}

	randomSigner, err := signer.New()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	session := s.ss.GetSessionByChallengerKey(userPK)
	if session != nil && session.Verified {

	}

	sharedKey := randomSigner.GetSharedKey(*userPK)

	hashedKey, ciphered, err := randomSigner.EncryptAndGetChallengeHash(sharedKey, uuid.NewString())
	if err != nil {
		http.Error(w, "unable to send challenge", http.StatusInternalServerError)
		return
	}

	serverSig, err := s.SignerObject.Sign(hex.EncodeToString(crypto.FromECDSAPub(randomSigner.GetPublicKey())))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	s.ss.CreateSession(&Session{
		ID:               sharedKey,
		Challenge:        ciphered,
		CreatedAt:        time.Now(),
		TTL:              10,
		Verified:         false,
		ChallengerPubKey: userPK,
		EphemeralKey:     randomSigner.GetPrivateKey(),
	})

	json.NewEncoder(w).Encode(ChallengePayloadResp{
		Challenge:       ciphered,
		HashedSharedKey: hashedKey,
		EphemeralPubKey: hex.EncodeToString(crypto.FromECDSAPub(randomSigner.GetPublicKey())),
		Signature:       serverSig,
	})
}

type ChallengePayloadReq struct {
	PubKey string `json:"public_key"`
}

type ChallengePayloadResp struct {
	Challenge       string `json:"challenge,"`
	HashedSharedKey string `json:"hashed_signature,"`
	EphemeralPubKey string `json:"ephemeral_pub_key"`
	Signature       string `json:"signature"`
}
