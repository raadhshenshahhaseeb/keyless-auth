package service

import (
	"context"
	"errors"
	"fmt"
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

	// TODO: hex and then encode to string before checking and storing
	exists, err := s.credRepo.Exists(newCredential)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to check for duplicates: %w", err)
	}
	if exists {
		return nil, nil, nil, errors.New("credential already exists")
	}

	// TODO: hash credentials, encode them and then store them to repo
	// hashedCred := hashCredential(newCredential)
	// hashedHex := hex.EncodeToString(hashedCred)
	//
	// err = s.credRepo.AddGlobalCredential(ctx, hashedHex)
	// if err != nil {
	// 	return nil, nil, nil, fmt.Errorf("failed to append to global credentials: %w", err)
	// }

	err = s.credRepo.AddGlobalCredential(ctx, newCredential)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("failed to append to global credentials: %w", err)
	}

	credentials, err := s.credRepo.GetAllGlobalCredentials(ctx)
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, nil, nil, fmt.Errorf("failed to retrieve global credentials: %w", err)
	}

	// TODO: storing hexed or encoded credentials and then decode them, refer to commented snippet below
	var data [][]byte
	for _, credential := range credentials {
		data = append(data, []byte(credential))
	}

	// data := make([][]byte, 0, len(credentials))
	// for _, cHex := range credentials {
	// 	cBytes, decodeErr := hex.DecodeString(cHex)
	// 	if decodeErr != nil {
	// 		log.Printf("Skipping invalid hex credential %q: %v", cHex, decodeErr)
	// 		continue
	// 	}
	// 	data = append(data, cBytes)
	// }

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
	if err != nil && errors.Is(err, redis.Nil) {
		return nil, nil, nil, fmt.Errorf("failed to get most recent merkle node: %w", err)
	}

	var prevRoot []byte
	if prevRecentNode != nil {
		prevRoot = prevRecentNode.TreeRoot
	}

	newNode := &repository.MerkleNode{
		ID:           uuid.New().String(),
		NodeType:     repository.Credential,
		Hash:         newCredential,
		ProofIndex:   proofIndex,
		ProofHashes:  proof.Hashes,
		TreeRoot:     tree.Root(),
		PrevRoot:     prevRoot,
		CreatedAt:    time.Now(),
		CredentialID: newCredential,
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
	for _, credential := range credentials {
		data = append(data, []byte(credential))
	}

	// TODO
	// data := make([][]byte, 0, len(credentials))
	// for _, cHex := range credentials {
	// 	cBytes, decodeErr := hex.DecodeString(cHex)
	// 	if decodeErr != nil {
	// 		log.Printf("Skipping invalid hex credential %q: %v", cHex, decodeErr)
	// 		continue
	// 	}
	// 	data = append(data, cBytes)
	// }

	// build tree
	tree, err := merkletree.NewUsing(data, keccak256.New(), []byte{})
	if err != nil {
		return nil, 0, err
	}

	// return root
	return tree, len(credentials), nil
}

func hashCredential(cred string) []byte {
	salt := []byte{0x1c, 0x9d, 0x3c, 0x4f}
	h := keccak256.New()
	credHash := h.Hash([]byte(cred))
	saltedHash := h.Hash(append(credHash, salt...))
	return saltedHash
}
