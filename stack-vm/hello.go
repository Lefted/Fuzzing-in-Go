package main

import (
	"fmt"
	"strings"
)

const (
	PLUS = iota // iota is a special constant that starts at 0 and increments by 1 for each const
	MULT
	ONE
	TWO
)

func showToken(code int) string {
	switch code {
	case PLUS:
		return "+"
	case MULT:
		return "*"
	case ONE:
		return "1"
	case TWO:
		return "2"
	default:
		return "unknown"
	}
}
func show(code []int) string {
	var builder strings.Builder
	for _, token := range code {
		builder.WriteString(showToken(token) + " ")
	}
	return builder.String()
}

////////////////////
// Expressions

type Exp interface {
	eval() int      // evaluate the expression (interpreter)
	convert() []int // convert to "reverse polish notation" (compiler)
}

type IntExp struct {
	value int
}

func NewIntExp(x int) Exp { // New function instead of constructor to return a pointer
	if x == 1 || x == 2 {
		return &IntExp{value: x}
	} else {
		fmt.Printf("Invalid value: %d. Must be 1 or 2\n", x)
		return &IntExp{value: 1}
	}
}

func (exp *IntExp) eval() int {
	return exp.value
}

func (exp *IntExp) convert() []int {
	n := ONE
	if exp.value == 2 {
		n = TWO
	}
	return []int{n}
}

type PlusExp struct {
	left  Exp
	right Exp
}

func NewPlusExp(left Exp, right Exp) Exp {
	return &PlusExp{left: left, right: right}
}

func (exp *PlusExp) eval() int {
	return exp.left.eval() + exp.right.eval()
}

func (exp *PlusExp) convert() []int {
	var v1 = exp.left.convert()
	var v2 = exp.right.convert()
	v1 = append(v1, v2...) // append v2 to v1 ... unpacks the slice
	v1 = append(v1, PLUS)
	return v1
}

type MultExp struct {
	left  Exp
	right Exp
}

func NewMultExp(left Exp, right Exp) Exp {
	return &MultExp{left: left, right: right}
}

func (exp MultExp) eval() int {
	return exp.left.eval() * exp.right.eval()
}

func (exp MultExp) convert() []int {
	var v1 = exp.left.convert()
	var v2 = exp.right.convert()
	v1 = append(v1, v2...)
	v1 = append(v1, MULT)
	return v1
}

// //////////////////
// VM run-time
type VM struct {
	codes []int
}

type VMRunnable interface {
	run() int
	showRunConvert()
	convert() Exp
}

func NewVM(codes []int) *VM {
	return &VM{codes: codes}
}

func (vm *VM) showRunConvert() {
	fmt.Println("VM code: ", show(vm.codes))
	fmt.Println("=> ", vm.run())
	fmt.Println("Exp: ", vm.convert().eval())
}

func (vm *VM) run() int {
	stack := []int{}
	for _, code := range vm.codes {
		switch code {
		case ONE:
			stack = append(stack, 1)
		case TWO:
			stack = append(stack, 2)
		case MULT:
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, left*right)
		case PLUS:
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, left+right)
		}
	}
	return stack[0]
}

func (vm *VM) convert() Exp {
	stack := []Exp{}
	for _, code := range vm.codes {
		switch code {
		case ONE:
			stack = append(stack, NewIntExp(1))
		case TWO:
			stack = append(stack, NewIntExp(2))
		case MULT:
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, NewMultExp(left, right))
		case PLUS:
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, NewPlusExp(left, right))
		}
	}
	return stack[0]
}

////////////////////
// Examples

func testVM() {
	{
		code := []int{ONE, TWO, TWO, PLUS, MULT}
		vm := NewVM(code)
		vm.showRunConvert()
	}

	{
		code := []int{ONE, TWO, PLUS, TWO, MULT}
		vm := NewVM(code)
		vm.showRunConvert()
	}
}

func testExp() {
	var e = NewPlusExp((NewMultExp(NewIntExp(1), NewIntExp(2))), NewIntExp(1))

	var run = func(e Exp) {
		fmt.Println("Exp yields ", e.eval())
		var vm = NewVM(e.convert())
		vm.showRunConvert()
	}

	run(e)
}

func main() {
	// code := []int{ONE, TWO, PLUS, ONE, TWO, MULT}
	// fmt.Println(show(code))

	// Create a new IntExp instance
	// exp1, err := NewIntExp(3)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// Evaluate the expression
	// fmt.Println("Eval: ", exp1.eval())
	// fmt.Println("Convert: ", show(exp1.convert()))

	// testVM()
	testExp()
}
