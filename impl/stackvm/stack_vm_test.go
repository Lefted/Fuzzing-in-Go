package stackvm_test

import (
	"math/rand"
	"testing"

	"project/impl/fuzzplus"
	"project/impl/stackvm"
	"project/impl/stackvm/encoding"
	gen "project/impl/stackvm/generator"
)

func TestVM(t *testing.T) {
	{
		code := []stackvm.Token{stackvm.One, stackvm.Two, stackvm.Two, stackvm.Plus, stackvm.Mult}
		vm := stackvm.NewVM(code)
		vm.ShowRunConvert()
	}

	{
		code := []stackvm.Token{stackvm.One, stackvm.Two, stackvm.Plus, stackvm.Two, stackvm.Mult}
		vm := stackvm.NewVM(code)
		vm.ShowRunConvert()
	}
	{
		code := []stackvm.Token{stackvm.One, stackvm.Two, stackvm.Div, stackvm.Two, stackvm.Mult, stackvm.Two, stackvm.Div, stackvm.Two, stackvm.Div}
		vm := stackvm.NewVM(code)
		vm.ShowRunConvert()
	}
	{
		code := []stackvm.Token{stackvm.Two, stackvm.One, stackvm.Div, stackvm.Two, stackvm.Mult, stackvm.Two, stackvm.Div, stackvm.Two, stackvm.Div}
		vm := stackvm.NewVM(code)
		vm.ShowRunConvert()
	}
}

func TestExp(t *testing.T) {
	var e = stackvm.NewPlusExp((stackvm.NewMultExp(stackvm.NewIntExp(1), stackvm.NewIntExp(2))), stackvm.NewIntExp(1))

	var run = func(e stackvm.Exp) {
		t.Logf("Exp yields %g", e.Eval())
		var vm = stackvm.NewVM(e.Convert())
		vm.ShowRunConvert()
	}

	run(e)
}

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
		expression := gen.RandomExp(rand, 3)

		// act
		// Property 1: Exp => convert => VMCode => convert Exp2
		vmCode := expression.Convert() // convert the expression to VM code
		vm := stackvm.NewVM(vmCode)    // create a new VM instance
		expression2 := vm.Convert()    // convert the VM code back to an expression

		resultFromExp := expression.Eval()
		resultFromVM := expression2.Eval()
		t.Logf("Result from VM: %g", resultFromVM)
		t.Logf("Result from Expression: %g", resultFromExp)

		// assert that Exp.eval == Exp2.eval
		if resultFromExp != resultFromVM {
			t.Errorf("Mismatch: original evaluation = %g, VM evaluation = %g", resultFromExp, resultFromVM)
		}
	})
}

func FuzzPlusExp(f *testing.F) {
	ff := fuzzplus.NewFuzzPlus(f)

	rand := rand.New(rand.NewSource(1))

	for i := 0; i < 500; i++ { // Problem: only 1 works because we generating results in different array lengths/ different number of arguments for the fuzz function
		exp := gen.RandomExp(rand, 1)
		encodedExp, err := encoding.EncodeWithDepth(exp, 2, 0)
		if err != nil {
			f.Fatalf("Failed to encode expression: %v", err)
		}
		// f.Logf("Length of encoded expression: %d", len(encodedExp))
		ff.Add2(encodedExp)
	}

	ff.Fuzz(func(t *testing.T, in []encoding.EncodedExp) {
		// t.Logf("Length of encoded expression: %d", len(in))

		// decode the encoded expression
		expression, err := encoding.Decode(in, 2)
		if err != nil {
			// t.Fatalf("Failed to decode expression: %v", err)
			return
		}

		// act
		// Property 1: Exp => convert => VMCode => convert Exp2
		vmCode := expression.Convert() // convert the expression to VM code
		vm := stackvm.NewVM(vmCode)    // create a new VM instance

		resultFromExp := expression.Eval()
		resultFromVM := vm.Run()

		// assert that Exp.eval == VM.run
		if resultFromExp != resultFromVM {
			t.Log(stackvm.Show(vmCode))
			t.Logf("Result from VM: %g", resultFromVM)
			t.Logf("Result from Expression: %g", resultFromExp)
			t.Errorf("Mismatch: original evaluation = %g, VM evaluation = %g", resultFromExp, resultFromVM)
		}
	})
}
