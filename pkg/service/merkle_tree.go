package service

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"

	"keyless-auth/repository"
)

type MerkleTreeService struct {
	credRepo *repository.CredentialsRepository
}

func NewMerkleTreeService(credRepo *repository.CredentialsRepository) *MerkleTreeService {
	return &MerkleTreeService{
		credRepo: credRepo,
	}
}

// TODO
// func (s *MerkleTreeService) GetMerkleTree() (*repository.MerkleNode, *merkletree.MerkleTree, error) {
// 	obj, err := s.credRepo.GetGlobalMerkleObject()
// 	if err != nil {
// 		return nil, nil, fmt.Errorf("merkle: failed to get most recent merkle node: %w", err)
// 	}
//
// 	node, tree, err :=
//
// 	return node, tree, nil
// }

// WithNewCredential returns merkle tree, node and an error.
func (s *MerkleTreeService) WithNewCredential(newCredential string) (*merkletree.MerkleTree, *repository.MerkleNode, *merkletree.Proof, error) {
	if newCredential == "" {
		return nil, nil, nil, errors.New("credential must not be empty")
	}

	ctx := context.Background()

	exists, err := s.credRepo.Exists(newCredential)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to check for duplicates: %w", err)
	}
	if exists {
		return nil, nil, nil, errors.New("credential already exists")
	}

	err = s.credRepo.AddGlobalCredential(ctx, newCredential)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to append to global credentials: %w", err)
	}

	credentials, err := s.credRepo.GetAllGlobalCredentials(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, nil, nil, fmt.Errorf("failed to retrieve global credentials: %w", err)
	}

	var data [][]byte
	for _, credential := range credentials {
		data = append(data, []byte(credential))
	}

	tree, err := merkletree.NewUsing(data, keccak256.New(), []byte{})
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to build Merkle tree: %w", err)
	}

	proofIndex := uint64(len(data) - 1)
	proof, err := tree.GenerateProof(data[proofIndex])
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to generate proof: %w", err)
	}

	prevRecentNode, err := s.credRepo.GetMostRecentMerkleNode(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, nil, nil, fmt.Errorf("failed to get most recent merkle node: %w", err)
	}

	var proofHashes []string
	for _, proofHash := range proof.Hashes {
		proofHashes = append(proofHashes, hex.EncodeToString(proofHash))
	}

	prevRecentRoot := ""
	if prevRecentNode != nil && len(prevRecentNode.TreeRoot) != 0 {
		prevRecentRoot = prevRecentNode.TreeRoot
	}

	newNode := &repository.MerkleNode{
		ID:               uuid.New().String(),
		NodeType:         repository.Credential,
		Hash:             newCredential,
		Position:         proof.Index,
		ProofHashes:      proofHashes,
		TreeRoot:         tree.String(),
		PrevRoot:         prevRecentRoot,
		CreatedAt:        time.Now(),
		ActualCredential: newCredential,
	}

	return tree, newNode, proof, nil
}

func (s *MerkleTreeService) GetMerkleTree() (*merkletree.MerkleTree, int, error) {
	// fetch all credentials
	credentials, err := s.credRepo.GetAllGlobalCredentials(context.Background())
	if err != nil {
		return nil, 0, err
	}

	var data [][]byte
	for _, cHex := range credentials {
		cBytes, decodeErr := hex.DecodeString(cHex)
		if decodeErr != nil {
			log.Printf("Skipping invalid hex credential %q: %v", cHex, decodeErr)
			continue
		}
		data = append(data, cBytes)
	}

	tree, err := merkletree.NewUsing(data, keccak256.New(), []byte{})
	if err != nil {
		return nil, 0, err
	}

	// return root
	return tree, len(credentials), nil
}

func HashCredential(cred string) string {
	salt := []byte{0x1c, 0x9d, 0x3c, 0x4f}
	h := keccak256.New()
	credHash := h.Hash([]byte(cred))
	saltedHash := h.Hash(append(credHash, salt...))
	return hex.EncodeToString(saltedHash)
}

func (s *MerkleTreeService) GenerateTree(nCredential string) (*repository.Tree, error) {
	err := s.credRepo.AddGlobalCredential(context.Background(), hex.EncodeToString([]byte(nCredential)))
	if err != nil {
		return nil, fmt.Errorf("failed to append to global credentials: %w", err)
	}

	credentials, err := s.credRepo.GetAllGlobalCredentials(context.Background())
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to retrieve global credentials: %w", err)
	}

	var data [][]byte
	for _, credential := range credentials {
		hexCredential, err := hex.DecodeString(credential)
		if err != nil {
			return nil, fmt.Errorf("failed to decode credential: %w", err)
		}

		data = append(data, hexCredential)
	}

	tree, err := merkletree.NewUsing(data, keccak256.New(), []byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to build Merkle tree: %w", err)
	}

	proofIndex := uint64(len(data) - 1)
	proof, err := tree.GenerateProof(data[proofIndex])
	if err != nil {
		return nil, fmt.Errorf("failed to generate proof: %w", err)
	}

	prevRecentNode, err := s.credRepo.RecentMerkleNode(context.Background())
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get most recent merkle node: %w", err)
	}

	var proofHashes []string
	for _, proofHash := range proof.Hashes {
		proofHashes = append(proofHashes, hex.EncodeToString(proofHash))
	}

	prevRecentRoot := ""
	if prevRecentNode != nil {
		prevRecentRoot = prevRecentNode.Tree.Root
	}

	return &repository.Tree{
		Leaf:          nCredential,
		Index:         proof.Index,
		ProofElements: proofHashes,
		Root:          hex.EncodeToString(tree.Root()),
		PrevRoot:      prevRecentRoot,
	}, nil
}
