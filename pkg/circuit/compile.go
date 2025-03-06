package circuit

import (
	"keyless-auth/domain"

	"github.com/consensus-shipyard/go-gnark/frontend"
	"github.com/consensus-shipyard/go-gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
)

func Compile() (*groth16.ProvingKey, *r1cs.R1CS, error) {
	var ckt ZKAuthCircuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &ckt)
	if err != nil {
		return nil, nil, err
	}

	pk, _, err := groth16.Setup(r1cs)
	if err != nil {
		return nil, err
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

func main() {
	Compile()
}
