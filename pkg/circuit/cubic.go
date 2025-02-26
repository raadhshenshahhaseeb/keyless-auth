package circuit

import (
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"math/big"
	"os"
)

type Circuit struct {
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

type cubic struct {
	System        constraint.ConstraintSystem
	Circuit       Circuit
	compileConfig frontend.CompileConfig
}

func (c *Circuit) Define(api frontend.API) error {
	x3 := api.Mul(c.X, c.X, c.X)
	api.AssertIsEqual(x3, c.Y)
	return nil
}

type Cubic interface {
}

func New(circuit *Circuit, field *big.Int) (Cubic, error) {
	if &circuit == nil {
		return nil, fmt.Errorf("circuit cannot be nil")
	}

	if field == nil {
		field = ecc.BN254.ScalarField()
	}

	builder := r1cs.NewBuilder

	system, _ := frontend.Compile(field, builder, circuit)
	witness, _ := frontend.NewWitness(circuit, ecc.BN254.ScalarField())
	publicWitness, _ := witness.Public()

	pk, vk, _ := groth16.Setup(system)
	proof, _ := groth16.Prove(system, pk, witness)
	err := groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}

	verifySolidityPath := fmt.Sprintf("..%conchain%ccontracts%ccubic_groth16.sol", os.PathSeparator, os.PathSeparator, os.PathSeparator)
	f, _ := os.OpenFile(verifySolidityPath, os.O_CREATE|os.O_WRONLY, 0666)
	defer f.Close()
	vk.ExportSolidity(f)

	return &cubic{
		System:  system,
		Circuit: *circuit,
	}, nil
}
