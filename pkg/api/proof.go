package api

import (
	"bytes"
	"encoding/json"
	"keyless-auth/circuit"
	"keyless-auth/repository"
	"net/http"
<<<<<<< HEAD
=======

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
>>>>>>> bbd94c8 (WIP: Circuit and merkle tree logic (#7))
)

type ProofRequest struct {
	Leaf      string   `json:"leaf"`      // Leaf hash
	Root      string   `json:"root"`      // Merkle root
	Siblings  []string `json:"siblings"`  // Sibling hashes
	Positions []int    `json:"positions"` // Positions (0 = left, 1 = right)
}

type ProofResponse struct {
	Proof []byte `json:"proof"`
}

type ProofHandler struct {
	walletRepo *repository.WalletRepository
}

func NewProofHandler(walletRepo *repository.WalletRepository) *ProofHandler {
	return &ProofHandler{
		walletRepo: walletRepo,
	}
}

func (h *ProofHandler) GenerateProof(w http.ResponseWriter, r *http.Request) {
	var req ProofRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

<<<<<<< HEAD
	proof, err := circuit.CompileCircuit(req)
=======
	var ckt circuit.ZKAuthCircuit
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &ckt)
>>>>>>> bbd94c8 (WIP: Circuit and merkle tree logic (#7))
	if err != nil {
		http.Error(w, "Failed to compile circuit", http.StatusInternalServerError)
		return
	}

<<<<<<< HEAD
=======
	pk, _, err := groth16.Setup(r1cs)
	if err != nil {
		http.Error(w, "Failed to generate keys", http.StatusInternalServerError)
		return
	}
	assignment := circuit.ZKAuthCircuit{
		Leaf:          frontend.Variable(req.Leaf),
		Root:          frontend.Variable(req.Root),
		ProofElements: make([]frontend.Variable, len(req.Siblings)),
		ProofIndex:    frontend.Variable(req.Positions),
	}

	for i := 0; i < len(req.Siblings); i++ {
		assignment.ProofElements[i] = frontend.Variable(req.Siblings[i])
	}
	assignment.ProofIndex = frontend.Variable(1)

	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		http.Error(w, "Failed to create witness", http.StatusInternalServerError)
		return
	}

	proof, err := groth16.Prove(r1cs, pk, witness)
	if err != nil {
		http.Error(w, "Failed to generate proof", http.StatusInternalServerError)
		return
	}

>>>>>>> bbd94c8 (WIP: Circuit and merkle tree logic (#7))
	var proofBytes []byte
	proof.WriteTo(bytes.NewBuffer(proofBytes))
	json.NewEncoder(w).Encode(ProofResponse{Proof: proofBytes})
}
