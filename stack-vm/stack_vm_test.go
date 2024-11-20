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
		// Problem: Guided fuzzer can only control seed and thus coverage exploration will be limited
		rand := rand.New(rand.NewSource(int64(seed)))

		// Generate a random expression
		expression := randomExp(rand, 3)

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

func TestGenerator(t *testing.T) {
	rand := rand.New(rand.NewSource(2))
	exp := randomExp(rand, 4)
	t.Logf("Random expression: %s", exp)
	var vm = NewVM(exp.convert())
	vm.showRunConvert()
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

func TestEncodeDecode(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	exp := randomExp(rand, 3)
	var vm = NewVM(exp.convert())
	vm.showRunConvert()
	encodedExp, err := Encode(exp)
	if err != nil {
		t.Fatalf("Failed to encode expression: %v", err)
	}

	t.Logf("Encoded expression: %v", encodedExp)

	decodedExp, err := Decode(encodedExp)
	if err != nil {
		t.Fatalf("Failed to decode expression: %v", err)
	}

	vmCode := decodedExp.convert() // convert the expression to VM code
	vm2 := NewVM(vmCode)           // create a new VM instance
	expression2 := vm2.convert()   // convert the VM code back to an expression

	resultFromExp := decodedExp.eval()
	resultFromVM := expression2.eval()
	t.Logf("Result from VM: %g", resultFromVM)
	t.Logf("Result from Expression: %g", resultFromExp)

	// assert that Exp.eval == Exp2.eval
	if resultFromExp != resultFromVM {
		t.Errorf("Mismatch: original evaluation = %g, VM evaluation = %g", resultFromExp, resultFromVM)
	}
	t.Logf("Decoded expression: %s", decodedExp)
}

func FuzzPlusExp(f *testing.F) {
	ff := NewFuzzPlus(f)

	rand := rand.New(rand.NewSource(1))

	for i := 0; i < 1; i++ { // Problem: only 1 works because we generating results in different array lengths/ different number of arguments for the fuzz function
		exp := randomExp(rand, 4)
		encodedExp, err := Encode(exp)
		if err != nil {
			f.Fatalf("Failed to encode expression: %v", err)
		}
		ff.Add2(encodedExp)
	}

	ff.Fuzz(func(t *testing.T, in []EncodedExp) {
		// decode the encoded expression
		expression, err := Decode(in)
		if err != nil {
			// t.Fatalf("Failed to decode expression: %v", err)
			return
		}
		t.Logf("Decoded expression: %s", expression)

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
