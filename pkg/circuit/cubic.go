package circuit

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"os"
)

type Circuit struct {
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

type cubic struct {
}

func (c *Circuit) Define(api frontend.API) error {
	x3 := api.Mul(c.X, c.X, c.X)
	api.AssertIsEqual(x3, c.Y)
	return nil
}

type Cubic interface {
}

func New(circuit *Circuit) (Cubic, error) {
	if circuit == nil {
		return nil, fmt.Errorf("circuit cannot be nil")
	}

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &Circuit{})
	if err != nil {
		return nil, fmt.Errorf("constraint system error: %w", err)
	}

	witness, err := frontend.NewWitness(circuit, ecc.BN254.ScalarField())
	if err != nil {
		return nil, fmt.Errorf("witness error: %w", err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return nil, fmt.Errorf("pub witness error: %w", err)
	}

	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		return nil, fmt.Errorf("setup error: %w", err)
	}

	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		return nil, fmt.Errorf("proof error: %w", err)
	}

	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}

	verifySolidityPath := fmt.Sprintf("../%ctransaction%conchain%ccontracts%ccubic_groth16.sol", os.PathSeparator, os.PathSeparator, os.PathSeparator, os.PathSeparator)
	f, err := os.OpenFile(verifySolidityPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open verify solution file: %w", err)
	}

	defer f.Close()
	vk.ExportSolidity(f)

	return &cubic{}, nil
}
