package repository

import (
	"fmt"
	"time"

	"github.com/wealdtech/go-merkletree"
)

// wallet address is mapped to a root
// each root has l_sibling and r_sibling
// r_sibling is always a credential or sub_root
// l_sibling is always a proof
// each root is always public

type NodeType int

const (
	Root NodeType = iota
	SubRoot
	Credential
	Proof
)

// MerkleNode describes a single node in the Merkle tree.
type MerkleNode struct {
	ID           string // e.g. a UUID
	NodeType     NodeType
	Hash         string // Hash of this leaf/node
	ProofIndex   uint64
	ProofHashes  [][]byte
	TreeRoot     []byte // Merkle root after insertion
	PrevRoot     []byte // Optional: the previous root
	CreatedAt    time.Time
	CredentialID string // which credential this node belongs to
}

var NodeTypeNames = map[NodeType]string{
	Root:       "root",
	SubRoot:    "sroot",
	Credential: "credential",
	Proof:      "proof",
}

func (t NodeType) String() string {
	if name, ok := NodeTypeNames[t]; ok {
		return name
	}
	return fmt.Sprintf("Unknown NodeType (%d)", t)
}

// -----------------TODO

// GlobalMerkleObject is for future reference
type GlobalMerkleObject struct {
	Node *MerkleNode            `json:"node"`
	Tree *merkletree.MerkleTree `json:"tree"`
}

func (o *GlobalMerkleObject) ToChildren() (*MerkleNode, *merkletree.MerkleTree, error) {
	return nil, nil, nil
}

func (o *GlobalMerkleObject) ToParent(*MerkleNode, *merkletree.MerkleTree, error) error {
	return nil
}
