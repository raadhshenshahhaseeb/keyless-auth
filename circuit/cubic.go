package main

import (
	"github.com/consensys/gnark/frontend"
)

// Circuit defines a simple circuit
// x**3 + x + 5 == y
type Circuit struct {
	// struct tags on a variable is optional
	// default uses variable name and secret visibility.
	X frontend.Variable `gnark:"x"`
	Y frontend.Variable `gnark:",public"`
}

// Define declares the circuit constraints
// x**3 + x + 5 == y
func (cubicCircuit *Circuit) DefineX(api frontend.API) error {
	x3 := api.Mul(cubicCircuit.X, cubicCircuit.X, cubicCircuit.X)
	api.AssertIsEqual(cubicCircuit.Y, api.Add(x3, cubicCircuit.X, 5))
	return nil
}
