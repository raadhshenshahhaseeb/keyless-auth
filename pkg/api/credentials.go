package api

import (
	"encoding/hex"
	"encoding/json"
	"keyless-auth/repository"
	"keyless-auth/service"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/wealdtech/go-merkletree"
)

type CredentialRequest struct {
	HashedCredential string `json:"hashed_credential"`
}

type CredentialResponse struct {
	WalletAddress string `json:"wallet_address"`
}

type MerkleRootResponse struct {
	MerkleRoot string `json:"merkle_root"`
	NumLeaves  int    `json:"num_leaves"`
}

type MerkleProofResponse struct {
	Proof *merkletree.Proof `json:"proof"`
}

type CredentialsHandler struct {
	credRepo   *repository.CredentialsRepository
	walletRepo *repository.WalletRepository
	merkleTree *service.MerkleTreeService
}

func NewCredentialsHandler(credRepo *repository.CredentialsRepository, walletRepo *repository.WalletRepository) *CredentialsHandler {
	return &CredentialsHandler{
		credRepo:   credRepo,
		walletRepo: walletRepo,
		merkleTree: service.NewMerkleTreeService(credRepo),
	}
}

func (h *CredentialsHandler) GetMerkleRoot(w http.ResponseWriter, r *http.Request) {
	tree, numLeaves, err := h.merkleTree.GetMerkleTree()
	if err != nil {
		http.Error(w, "Failed to get merkle root", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(MerkleRootResponse{
		MerkleRoot: "0x" + hex.EncodeToString(tree.Root()),
		NumLeaves:  numLeaves,
	})
}

func (h *CredentialsHandler) GenerateMerkleProof(w http.ResponseWriter, r *http.Request) {
	credential := mux.Vars(r)["credential"]

	proof, err := h.merkleTree.GenerateMerkleProof(credential)
	if err != nil {
		http.Error(w, "Failed to generate merkle proof", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(MerkleProofResponse{
		Proof: proof,
	})
}

func GenerateWalletAddress() (string, []byte, error) {
	privKey, err := crypto.GenerateKey()
	if err != nil {
		return "", nil, err
	}

	return crypto.PubkeyToAddress(privKey.PublicKey).Hex(), privKey.D.Bytes(), nil
}

func (h *CredentialsHandler) GenerateCredential(w http.ResponseWriter, r *http.Request) {
	var req CredentialRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// check if credential already exists
	if exists, err := h.credRepo.DoesCredentialExist(req.HashedCredential); err != nil || exists {
		// TODO: better return 200 with message to fetch the wallet address or credential on a different endpoint
		http.Error(w, "Credential already exists", http.StatusBadRequest)
		return
	}

	// store leaf in merkle tree
	if err := h.credRepo.SaveCredential(req.HashedCredential); err != nil {
		log.Printf("Failed to save credential: %v", err)
		http.Error(w, "Failed to save credential", http.StatusInternalServerError)
		return
	}

	// generate wallet address
	walletAddress, privKey, err := GenerateWalletAddress()
	if err != nil {
		http.Error(w, "Failed to generate wallet address", http.StatusInternalServerError)
		return
	}

	// store wallet
	if err := h.walletRepo.Save(walletAddress, privKey, req.HashedCredential); err != nil {
		log.Printf("Failed to save wallet: %v", err)
		http.Error(w, "Failed to save wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CredentialResponse{WalletAddress: walletAddress})
}

func (h *CredentialsHandler) GetWalletAddressByCredential(w http.ResponseWriter, r *http.Request) {
	credential := mux.Vars(r)["credential"]

	wallet, err := h.walletRepo.GetWalletByCredential(credential)
	if err != nil {
		http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CredentialResponse{WalletAddress: wallet.Address})
}
