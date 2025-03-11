package main

import (
	"errors"
	"flag"
	"log"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/hash/mimc"

	"github.com/PolyhedraZK/ExpanderCompilerCollection/ecgo"
	"github.com/PolyhedraZK/ExpanderCompilerCollection/ecgo/test"
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

func GenerateGroth16(assignment *ZKAuthCircuit) error {
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, assignment)
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

func GenerateExpander() error {
	assignment := &ZKAuthCircuit{
		Root: "123456789012345678901234567890123456789012345678901234567890abcd",
		ProofElements: []frontend.Variable{
			"234567890123456789012345678901234567890123456789012345678901dcba",
			"345678901234567890123456789012345678901234567890123456789012efab",
		},
		ProofIndex: 0,
		Leaf:       "123456789012345678901234567890123456789012345678901234567890fedc",
	}
	circuit, err := ecgo.Compile(ecc.BN254.ScalarField(), assignment)
	if err != nil {
		return err
	}

	c := circuit.GetLayeredCircuit()
	os.WriteFile("expander_circuit.txt", c.Serialize(), 0o644)
	inputSolver := circuit.GetInputSolver()
	witness, err := inputSolver.SolveInputAuto(nil)
	if err != nil {
		return err
	}
	os.WriteFile("expander_witness.txt", witness.Serialize(), 0o644)
	if !test.CheckCircuit(c, witness) {
		return errors.New("witness is not valid")
	}
	return nil
}

func main() {
	// Command line flags:
	// -expander: Generate expander circuit and witness files
	//   Example: ./main -expander
	//
	// -groth16: Generate Groth16 proving/verification keys and Solidity verifier
	//   Example: ./main -groth16
	expander := flag.Bool("expander", false, "Generate expander")
	groth16 := flag.Bool("groth16", false, "Generate Groth16 keys")
	flag.Parse()
	if !*expander && !*groth16 {
		log.Fatal("Please provide a command: 'expander' or 'groth16'")
	}

	if *expander {
		err := GenerateExpander()
		if err != nil {
			log.Fatalf("Failed to generate expander: %v", err)
		}
	} else if *groth16 {
		err := GenerateGroth16(&ZKAuthCircuit{
			Root: "ed6cfd06d0c37b1f964a2e63e5dce2bab107297c4a5518392d08a2cea24794dc",
			ProofElements: []frontend.Variable{
				make([]frontend.Variable, 2)},
			ProofIndex: 2,
			Leaf:       "aae457593db8c6ab406c81939d4ffb39e9ac16aeceeb6e109d1593aff3b91ecd",
		})
		if err != nil {
			log.Fatalf("Failed to generate Groth16 keys: %v", err)
		}
	} else {
		log.Fatalf("Invalid argument: %v", os.Args[1])
	}
}
