package circuit

import (
	"os"

	"keyless-auth/domain"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

var (
	MAX_DEPTH = 256
)

func Compile() (*groth16.ProvingKey, constraint.ConstraintSystem, error) {
	var ckt ZKAuthCircuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &ckt)
	if err != nil {
		return nil, nil, err
	}

	pk := groth16.NewProvingKey(ecc.BN254)
	{
		f, err := os.Open("mt.g16.pk")
		if err != nil {
			return nil, nil, err
		}
		_, err = pk.ReadFrom(f)
		if err != nil {
			return nil, nil, err
		}
		f.Close()
	}
	return &pk, r1cs, nil
}

func CompileCircuit(proof domain.Proof) (*groth16.Proof, error) {
	pk, r1cs, err := Compile()
	if err != nil {
		return nil, err
	}

	assignment := ZKAuthCircuit{
		Leaf:          frontend.Variable(proof.Leaf),
		Root:          frontend.Variable(proof.Root),
		ProofElements: make([]frontend.Variable, len(proof.Siblings)),
		ProofIndex:    frontend.Variable(proof.Positions),
	}

	for i := 0; i < len(proof.Siblings); i++ {
		assignment.ProofElements[i] = frontend.Variable(proof.Siblings[i])
	}
	assignment.ProofIndex = frontend.Variable(1)

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}

	prf, err := groth16.Prove(r1cs, *pk, witness)
	if err != nil {
		return nil, err
	}

	return &prf, nil
}
