package encoding

import (
	"math/rand"
	"testing"

	"project/impl/stackvm"
	gen "project/impl/stackvm/generator"
)

func TestEncodeDecode(t *testing.T) {
	rand := rand.New(rand.NewSource(1))
	exp := gen.RandomExp(rand, 1)
	var vm = stackvm.NewVM(exp.Convert())
	vm.ShowRunConvert()
	encodedExp, err := EncodeWithDepth(exp, 4, 0)
	if err != nil {
		t.Fatalf("Failed to encode expression: %v", err)
	}

	t.Logf("Encoded expression: %v", encodedExp)

	decodedExp, err := Decode(encodedExp, 4)
	if err != nil {
		t.Fatalf("Failed to decode expression: %v", err)
	}

	vmCode := decodedExp.Convert() // convert the expression to VM code
	vm2 := stackvm.NewVM(vmCode)   // create a new VM instance
	expression2 := vm2.Convert()   // convert the VM code back to an expression

	resultFromExp := decodedExp.Eval()
	resultFromVM := expression2.Eval()
	t.Logf("Result from VM: %g", resultFromVM)
	t.Logf("Result from Expression: %g", resultFromExp)

	// assert that Exp.eval == Exp2.eval
	if resultFromExp != resultFromVM {
		t.Errorf("Mismatch: original evaluation = %g, VM evaluation = %g", resultFromExp, resultFromVM)
	}
	t.Logf("Decoded expression: %s", decodedExp)
}
