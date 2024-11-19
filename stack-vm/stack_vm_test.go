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
		// Property 1: Exp => convert => VMCode => convert Exp2
		vmCode := expression.convert() // convert the expression to VM code
		vm := NewVM(vmCode)            // create a new VM instance
		expression2 := vm.convert()    // convert the VM code back to an expression

		resultFromExp := expression.eval()
		resultFromVM := expression2.eval()
		t.Logf("Result from VM: %g", resultFromVM)
		t.Logf("Result from Expression: %g", resultFromExp)

		// assert that Exp.eval == Exp2.eval
		if resultFromExp != resultFromVM {
			t.Errorf("Mismatch: original evaluation = %g, VM evaluation = %g", resultFromExp, resultFromVM)
		}
	})
}

func TestGenerateRandomExp(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	exp := generateRandomExp(5, rand)
	t.Logf("Random expression: %s", exp)
	var vm = NewVM(exp.convert())
	vm.showRunConvert()
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

func randomIntExp(rand *rand.Rand) Exp {
	value := rand.Intn(2) + 1
	return NewIntExp(value)
}

func randomPlusExp(rand *rand.Rand, depth int) Exp {
	left := randomExp(rand, depth-1)
	right := randomExp(rand, depth-1)
	return NewPlusExp(left, right)
}

func randomMultExp(rand *rand.Rand, depth int) Exp {
	left := randomExp(rand, depth-1)
	right := randomExp(rand, depth-1)
	return NewMultExp(left, right)
}

func randomDivExp(rand *rand.Rand, depth int) Exp {
	left := randomExp(rand, depth-1)
	right := randomExp(rand, depth-1)
	return NewDivExp(left, right)
}

func randomExp(rand *rand.Rand, depth int) Exp {
	if depth <= 0 {
		return randomIntExp(rand)
	}

	operator := rand.Intn(4)
	switch operator {
	case 0:
		return randomIntExp(rand)
	case 1:
		return randomPlusExp(rand, depth)
	case 2:
		return randomMultExp(rand, depth)
	case 3:
		return randomDivExp(rand, depth)
	default:
		return randomIntExp(rand)
	}
}
