package stackvm_test

import (
	"math"
	"math/rand"
	"project/impl/fuzzplus"
	"project/impl/stackvm"
	"project/impl/stackvm/encoding"
	gen "project/impl/stackvm/generator"
	"testing"
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

func TestDivBug(t *testing.T) {
	// arrange
	exp := stackvm.NewDivExp( //
		stackvm.NewPlusExp( //
			stackvm.NewIntExp(2), stackvm.NewIntExp(1)), //
		stackvm.NewIntExp(1))

	// act
	vmCode := exp.Convert()
	vm := stackvm.NewVM(vmCode)
	resultFromExp := exp.Eval()
	resultFromVM := vm.Run()

	// assert that Exp.eval == VM.run
	if resultFromExp != resultFromVM {
		t.Log(stackvm.Show(vmCode))
		t.Logf("VM yields: %g", resultFromVM)
		t.Logf("Exp yields: %g", resultFromExp)
		t.Errorf("Mismatch, delta is %g", math.Abs(resultFromExp-resultFromVM))
	}
}

func FuzzWithGenerator(f *testing.F) {
	// No need to add seed inputs as the generator will generate random expressions

	f.Fuzz(func(t *testing.T, seed int) {
		// use seed for reproducibility
		rand := rand.New(rand.NewSource(int64(seed)))

		// use generator to create a random expression
		exp := gen.RandomExp(rand, 3)

		// act
		vmCode := exp.Convert()
		vm := stackvm.NewVM(vmCode)
		resultFromExp := exp.Eval()
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

func FuzzPlusExpNonResilient(f *testing.F) {
	ff := fuzzplus.NewFuzzPlus(f)

	rand := rand.New(rand.NewSource(1))

	for i := 0; i < 1; i++ {
		exp := gen.RandomExp(rand, 1)
		encodedExp, err := encoding.EncodeWithDepth(exp, 2, 0)
		if err != nil {
			f.Fatalf("Failed to encode expression: %v", err)
		}
		ff.Add2(encodedExp)
	}

	ff.Fuzz(func(t *testing.T, in []encoding.EncodedExp) {
		// decode the encoded expression
		expression, err := encoding.Decode(in, 2, false)
		if err != nil {
			return
		}

		// act
		vmCode := expression.Convert()
		vm := stackvm.NewVM(vmCode)

		resultFromExp := expression.Eval()
		resultFromVM := vm.Run()

		// assert that Exp.eval == VM.run
		if resultFromExp != resultFromVM {
			t.Log(stackvm.Show(vmCode))
			t.Logf("Result from VM: %g", resultFromVM)
			t.Logf("Result from Expression: %g", resultFromExp)
			t.Errorf("Mismatch: original evaluation = %g, VM evaluation = %g", resultFromExp, resultFromVM)
		}

		t.Log("Test")
	})
}

func FuzzPlusExpResilient(f *testing.F) {
	ff := fuzzplus.NewFuzzPlus(f)

	rand := rand.New(rand.NewSource(1))

	for i := 0; i < 500; i++ {
		exp := gen.RandomExp(rand, 1)
		encodedExp, err := encoding.EncodeWithDepth(exp, 2, 0)
		if err != nil {
			f.Fatalf("Failed to encode expression: %v", err)
		}
		ff.Add2(encodedExp)
	}

	ff.Fuzz(func(t *testing.T, in []encoding.EncodedExp) {
		// decode the encoded expression
		expression, err := encoding.Decode(in, 2, true)
		if err != nil {
			return
		}

		// act
		vmCode := expression.Convert()
		vm := stackvm.NewVM(vmCode)

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
