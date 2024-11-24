package generator

import (
	"math/rand"
	"project/impl/stackvm"
	"testing"
)

func TestGenerator(t *testing.T) {
	rand := rand.New(rand.NewSource(2))
	exp := RandomExp(rand, 4)
	t.Logf("Random expression: %s", exp)
	var vm = stackvm.NewVM(exp.Convert())
	vm.ShowRunConvert()
}
