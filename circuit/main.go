package main

import (
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"
)

var (
	MAX_DEPTH = 256
)

// ZKAuthCircuit proves credential hash inclusion in a Merkle tree
type ZKAuthCircuit struct {
	Root          frontend.Variable   `gnark:",public"`
	ProofElements []frontend.Variable // private
	ProofIndex    frontend.Variable   // private
	Leaf          frontend.Variable   // private
}

// Define the zk circuit
func (circuit *ZKAuthCircuit) Define(api frontend.API) error {
	h, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	// Hash leaf
	h.Reset()
	h.Write(circuit.Leaf)
	hashed := h.Sum()

	depth := len(circuit.ProofElements)
	if depth == 0 {
		depth = MAX_DEPTH
	}
	proofIndices := api.ToBinary(circuit.ProofIndex, depth)

	// // Continuously hash with the proof elements
	for i := 0; i < len(circuit.ProofElements); i++ {
		element := circuit.ProofElements[i]
		// 0 = left, 1 = right
		index := proofIndices[i]

		d1 := api.Select(index, element, hashed)
		d2 := api.Select(index, hashed, element)

		h.Reset()
		h.Write(d1, d2)
		hashed = h.Sum()
	}

	// Verify calculates hash is equal to root
	api.AssertIsEqual(hashed, circuit.Root)
	return nil
}

func GenerateGroth16() error {
	var circuit ZKAuthCircuit

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return err
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return err
	}
	{
		f, err := os.Create("mt.g16.vk")
		if err != nil {
			return err
		}
		_, err = vk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}
	{
		f, err := os.Create("mt.g16.pk")
		if err != nil {
			return err
		}
		_, err = pk.WriteRawTo(f)
		if err != nil {
			return err
		}
	}

	{
		f, err := os.Create("contract_mt.g16.sol")
		if err != nil {
			return err
		}
		err = vk.ExportSolidity(f)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	err := GenerateGroth16()
	if err != nil {
		log.Fatalf("Failed to generate Groth16 keys: %v", err)
	}
}
