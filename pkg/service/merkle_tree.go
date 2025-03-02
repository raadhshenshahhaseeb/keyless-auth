package service

import (
	"keyless-auth/repository"

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

func (s *MerkleTreeService) GenerateMerkleTree(newCredential string) ([]byte, error) {
	// fetch all credentials
	credentials, err := s.credRepo.GetCredentials()
	if err != nil {
		return nil, err
	}

	var data [][]byte
	for _, credential := range credentials {
		data = append(data, []byte(credential))
	}

	// add new credential
	data = append(data, []byte(newCredential))

	// build tree
	tree, err := merkletree.NewUsing(data, keccak256.New(), []byte{})
	if err != nil {
		return nil, err
	}

	// return root
	return tree.Root(), nil
}

func (s *MerkleTreeService) GenerateMerkleProof(credential string) (*merkletree.Proof, error) {
	// fetch all credentials
	credentials, err := s.credRepo.GetCredentials()
	if err != nil {
		return nil, err
	}

	var data [][]byte
	for _, credential := range credentials {
		data = append(data, []byte(credential))
	}

	// build tree
	tree, err := merkletree.NewUsing(data, keccak256.New(), []byte{})
	if err != nil {
		return nil, err
	}

	// generate proof
	proof, err := tree.GenerateProof([]byte(credential))
	if err != nil {
		return nil, err
	}

	return proof, nil
}
