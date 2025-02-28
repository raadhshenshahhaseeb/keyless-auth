package circuit

import (
	"fmt"
	"github.com/consensys/gnark/std/hash"

	"github.com/consensys/gnark/frontend"
	poseidon2 "github.com/consensys/gnark/std/permutation/poseidon2"
)

/**
TODO: Create API for proof generation
const proof = tree.getProof(hashedCredential);
const positions = proof.map(p => p.position === 'left' ? 0 : 1);

const data = {
    leaf: hashedCredential,
    root: merkleRoot,
    siblings: proof.map(p => p.data.toString('hex')),
    positions: positions
};

fetch('/generate-proof', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data)
});
**/

// Circuits defines the zk-SNARK circuit
type Circuits struct {
	Leaf      frontend.Variable   `gnark:",public"` // Public input: leaf
	Root      frontend.Variable   `gnark:",public"` // Public input: Merkle root
	Siblings  []frontend.Variable // Private input: sibling hashes
	Positions []frontend.Variable // Private input: positions (0 = left, 1 = right)
}

func NewMerkleDamgardHasher(api frontend.API) (hash.FieldHasher, error) {
	f, err := poseidon2.NewPoseidon2(api)
	if err != nil {
		return nil, fmt.Errorf("could not create poseidon2 hasher: %w", err)
	}
	return hash.NewMerkleDamgardHasher(api, f, 0), nil
}

// Define declares the circuit's constraints
func (circuit *Circuits) Define(api frontend.API) error {
	hash, _ := NewMerkleDamgardHasher(api)
	// Initialize the computed hash with the leaf
	computedHash := circuit.Leaf

	// Iterate through the proof (sibling hashes and positions)
	for i := 0; i < len(circuit.Siblings); i++ {
		// Check if the sibling is on the left or right
		if circuit.Positions[i] == 0 {
			// Sibling is on the left: hash(sibling, computedHash)
			computedHash = hash.Hash(circuit.Siblings[i], computedHash)
		} else {
			// Sibling is on the right: hash(computedHash, sibling)
			computedHash = hash.Hash(computedHash, circuit.Siblings[i])
		}
	}

	// Ensure the computed hash matches the given Merkle root
	api.AssertIsEqual(computedHash, circuit.Root)

	return nil
}
