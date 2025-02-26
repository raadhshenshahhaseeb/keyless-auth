package main

import (
	"fmt"
	"os"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type Circuit struct {
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

func (c *Circuit) Define(api frontend.API) error {
	x3 := api.Mul(c.X, c.X, c.X)
	api.AssertIsEqual(x3, c.Y)
	return nil
}

func main() {
	assignment := Circuit{
		X: 2,
		Y: 8,
	}

	r1cs, _ := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit{})
	witness, _ := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	pk, vk, _ := groth16.Setup(r1cs)
	proof, _ := groth16.Prove(r1cs, pk, witness)
	err := groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}

	verifySolidityPath := fmt.Sprintf("..%conchain%ccontracts%ccubic_groth16.sol", os.PathSeparator, os.PathSeparator, os.PathSeparator)
	f, _ := os.OpenFile(verifySolidityPath, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	vk.ExportSolidity(f)
}
