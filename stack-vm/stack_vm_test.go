package main

import (
	"math/rand"
	"testing"
)

// Fuzz test to generate random expressions and verify their evaluations
func FuzzRandomExp(f *testing.F) {
	// Add some seed inputs for the fuzz testing
	// f.Add(1) // Seed input can be expanded as needed
	// f.Add(2)
	// f.Add(3)

	testInputs := []int{5, 0, 50}

	for _, tc := range testInputs {
		f.Add(tc) // Use f.Add to provide a seed corpus
	}

	f.Fuzz(func(t *testing.T, seed int) {
		// Random number generator with set seed for reproducibility
		rand := rand.New(rand.NewSource(int64(seed)))

		// Generate a random expression
		expression := generateRandomExp(3, rand) // Adjust the depth/complexity as needed

		// act
		vmCode := expression.convert() // convert the expression to VM code
		vm := NewVM(vmCode)            // create a new VM instance
		expression2 := vm.convert()    // convert the VM code back to an expression

		resultFromExp := expression.eval()
		resultFromVM := expression2.eval()
		t.Logf("Result from VM: %d", resultFromVM)
		t.Logf("Result from Expression: %d", resultFromExp)

		// assert that the results are the same
		if resultFromExp != resultFromVM {
			t.Errorf("Mismatch: original evaluation = %d, VM evaluation = %d", resultFromExp, resultFromVM)
		}
	})
}

// Function to generate a random expression
func generateRandomExp(depth int, rand *rand.Rand) Exp {
	if depth <= 0 {
		// Randomly choose between 1 and 2
		return NewIntExp(rand.Intn(2) + 1)
	}

	// Randomly choose an operator: 0 for Plus, 1 for Mult
	operator := rand.Intn(2)

	var left, right Exp
	left = generateRandomExp(depth-1, rand)
	right = generateRandomExp(depth-1, rand)

	if operator == 0 {
		return NewPlusExp(left, right)
	}
	return NewMultExp(left, right)
}
