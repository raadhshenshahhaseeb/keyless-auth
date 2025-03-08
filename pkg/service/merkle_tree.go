package service

import (
	"keyless-auth/repository"
	"log"

	"github.com/wealdtech/go-merkletree"
	"github.com/wealdtech/go-merkletree/keccak256"
)

type MerkleTreeService struct {
	credRepo *repository.CredentialsRepository
}

func NewMerkleTreeService(credRepo *repository.CredentialsRepository) *MerkleTreeService {
	return &MerkleTreeService{
		credRepo: credRepo,
	}
}

func (s *MerkleTreeService) GetMerkleTree() (*merkletree.MerkleTree, int, error) {
	// fetch all credentials
	credentials, err := s.credRepo.GetCredentials()
	if err != nil {
		return nil, 0, err
	}

	var data [][]byte
	for _, credential := range credentials {
		data = append(data, []byte(credential))
	}

	// build tree
	tree, err := merkletree.NewUsing(data, keccak256.New(), []byte{})
	if err != nil {
		return nil, 0, err
	}

	// return root
	return tree, len(credentials), nil
}

func (s *MerkleTreeService) GenerateMerkleProof(credential string) (*merkletree.Proof, error) {
	// fetch all credentials
	tree, _, err := s.GetMerkleTree()
	if err != nil {
		log.Printf("Failed to get merkle tree: %v", err)
		return nil, err
	}

	// generate proof
	proof, err := tree.GenerateProof([]byte(credential))
	if err != nil {
		log.Printf("Failed to generate merkle proof: %v", err)
		return nil, err
	}

	return proof, nil
}
