package api

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
	"github.com/wealdtech/go-merkletree"

	"keyless-auth/repository"
	"keyless-auth/service"
)

type GenerateTreeRequest struct {
	HashedCredential string `json:"hashed_credential"`
}

type GenerateTreeResponse struct {
	Object repository.Object `json:"object,inline"`
}

type CredentialRequest struct {
	HashedCredential string `json:"hashed_credential"`
}

type CredentialResponse struct {
	MerkleRoot    string            `json:"merkle_root"`
	WalletAddress string            `json:"wallet_address"`
	Proof         *merkletree.Proof `json:"proof"`
	Leaf          string            `json:"leaf"`
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
	// TODO: with existing credential
	tree, node, proof, err := h.merkleTree.WithNewCredential(credential)
	if err != nil {
		http.Error(w, "failed to generate merkle proof", http.StatusInternalServerError)
		return
	}

	// "SaveCredentialAndNode" if you want to store the root. We can also store node only.
	err = h.credRepo.SaveCredentialAndNode(context.Background(), credential, hex.EncodeToString(tree.Root()), node)
	if err != nil {
		log.Println("failed to store credential and node: ", err)
		http.Error(w, "failed to store credential and node", http.StatusInternalServerError)
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
	if exists, err := h.credRepo.Exists(req.HashedCredential); err != nil || exists {
		// TODO: better return 200 with message to fetch the wallet address or credential on a different endpoint
		http.Error(w, "Credential already exists", http.StatusBadRequest)
		return
	}

	// generate wallet address
	walletAddress, privKey, err := GenerateWalletAddress()
	if err != nil {
		http.Error(w, "Failed to generate wallet address", http.StatusInternalServerError)
		return
	}

	tree, node, proof, err := h.merkleTree.WithNewCredential(req.HashedCredential)
	if err != nil {
		http.Error(w, "Failed to generate merkle tree root", http.StatusInternalServerError)
		return
	}

	// "SaveCredentialAndNode" if you want to store the root. We can also store node only.
	err = h.credRepo.SaveCredentialAndNode(context.Background(), node.ActualCredential, tree.String(), node)
	if err != nil {
		log.Println("failed to store credential and node: ", err)
		http.Error(w, "failed to store credential and node", http.StatusInternalServerError)
		return
	}

	err = h.credRepo.SetMostRecentMerkleNode(context.Background(), node)
	if err != nil {
		log.Println("failed to store recent node: ", err)
		http.Error(w, "failed to store recent node", http.StatusInternalServerError)
		return
	}

	// store wallet
	if err := h.walletRepo.Save(walletAddress, privKey, node.Hash, tree.String()); err != nil {
		log.Printf("Failed to save wallet: %v", err)
		http.Error(w, "Failed to save wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CredentialResponse{
		MerkleRoot:    tree.String(),
		WalletAddress: walletAddress,
		Proof:         proof,
		Leaf:          node.Hash,
	})
}

func (h *CredentialsHandler) GetWalletByCredential(w http.ResponseWriter, r *http.Request) {
	credential := mux.Vars(r)["credential"]

	wallet, err := h.walletRepo.GetWalletByCredential(credential)
	if err != nil {
		http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CredentialResponse{MerkleRoot: wallet.MerkleRoot})
}

func (h *CredentialsHandler) GetWalletIfExists(w http.ResponseWriter, r *http.Request) {
	credential := mux.Vars(r)["credential"]

	wallet, err := h.walletRepo.GetWalletByCredential(credential)
	if err != nil {
		http.Error(w, "Failed to get wallet", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(CredentialResponse{MerkleRoot: wallet.MerkleRoot})
}

func (h *CredentialsHandler) GenerateTreeObject(w http.ResponseWriter, r *http.Request) {
	var req GenerateTreeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.HashedCredential == "" {
		http.Error(w, "invalid or empty request object", http.StatusBadRequest)
		return
	}

	// check if credential already exists
	if exists, err := h.credRepo.Exists(req.HashedCredential); err != nil || exists {
		http.Error(w, "Credential already exists", http.StatusUnprocessableEntity)
		return
	}

	// TODO: get user from session state
	// TODO: get user wallet

	treeObj, err := h.merkleTree.GenerateTree(req.HashedCredential)
	if err != nil {
		http.Error(w, "Failed to generate merkle tree root", http.StatusInternalServerError)
		return
	}

	object := &repository.Object{
		Tree:      treeObj,
		User:      nil,
		CreatedAt: time.Now().UTC(),
	}

	err = h.credRepo.AddToTree(context.Background(), object)
	if err != nil {
		log.Println("failed to store credential and node: ", err)
		http.Error(w, "failed to store credential and node", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(&object)
}

// Register with a credential
// We generate a wallet
// We generate merkle tree -> proof
// merkle tree/proof
