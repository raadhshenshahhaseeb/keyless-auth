package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"keyless-auth/pkg/circuit"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

type ProofRequest struct {
	Leaf      string   `json:"leaf"`      // Leaf hash
	Root      string   `json:"root"`      // Merkle root
	Siblings  []string `json:"siblings"`  // Sibling hashes
	Positions []int    `json:"positions"` // Positions (0 = left, 1 = right)
}

// API response payload
type ProofResponse struct {
	Proof []byte `json:"proof"` // Generated proof
}

func generateProofHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var req ProofRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Compile the circuit
	var circuit circuit.Circuits
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		http.Error(w, "Failed to compile circuit", http.StatusInternalServerError)
		return
	}

	// Generate proving and verification keys
	pk, _, err := groth16.Setup(r1cs)
	if err != nil {
		http.Error(w, "Failed to generate keys", http.StatusInternalServerError)
		return
	}

	// Create the witness
	witness := circuit.Circuit{
		Leaf:      frontend.Variable(req.Leaf),
		Root:      frontend.Variable(req.Root),
		Siblings:  make([]frontend.Variable, len(req.Siblings)),
		Positions: make([]frontend.Variable, len(req.Positions)),
	}
	for i := 0; i < len(req.Siblings); i++ {
		witness.Siblings[i] = frontend.Variable(req.Siblings[i])
		witness.Positions[i] = frontend.Variable(req.Positions[i])
	}

	// Generate the proof
	proof, err := groth16.Prove(r1cs, pk, &witness)
	if err != nil {
		http.Error(w, "Failed to generate proof", http.StatusInternalServerError)
		return
	}

	// Return the proof
	var proofBytes []byte
	proof.WriteTo(bytes.NewBuffer(proofBytes))
	resp := ProofResponse{Proof: proofBytes}
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/generate-proof", generateProofHandler)
	http.ListenAndServe(":8080", nil)
}
