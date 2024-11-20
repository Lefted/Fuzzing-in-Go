package main

import (
	"fmt"
	"strings"
)

type Token int

const (
	Plus Token = iota // iota is a special constant that starts at 0 and increments by 1 for each const
	Mult
	Div
	One
	Two
)

func showToken(code Token) string {
	switch code {
	case Plus:
		return "+"
	case Mult:
		return "*"
	case Div:
		return "/"
	case One:
		return "1"
	case Two:
		return "2"
	default:
		return "Unknown"
	}
}

func show(code []Token) string {
	var builder strings.Builder
	for _, token := range code {
		builder.WriteString(showToken(token) + " ")
	}
	return builder.String()
}

////////////////////
// Expressions

type Exp interface {
	eval() float64    // evaluate the expression (interpreter)
	convert() []Token // convert to "reverse polish notation" (compiler)
}

type IntExp struct {
	value int
}

func NewIntExp(x int) Exp {
	if x == 1 || x == 2 {
		return &IntExp{value: x}
	} else {
		fmt.Printf("Invalid value: %d. Must be 1 or 2\n", x)
		return &IntExp{value: 1}
	}
}

func (exp *IntExp) eval() float64 {
	return float64(exp.value)
}

func (exp *IntExp) convert() []Token {
	n := One
	if exp.value == 2 {
		n = Two
	}
	return []Token{n}
}

type PlusExp struct {
	left  Exp
	right Exp
}

func NewPlusExp(left Exp, right Exp) Exp {
	return &PlusExp{left: left, right: right}
}

func (exp *PlusExp) eval() float64 {
	return exp.left.eval() + exp.right.eval()
}

func (exp *PlusExp) convert() []Token {
	var v1 = exp.left.convert()
	var v2 = exp.right.convert()
	v1 = append(v1, v2...) // append v2 to v1 ... unpacks the slice
	v1 = append(v1, Plus)
	return v1
}

type MultExp struct {
	left  Exp
	right Exp
}

func NewMultExp(left Exp, right Exp) Exp {
	return &MultExp{left: left, right: right}
}

func (exp MultExp) eval() float64 {
	return exp.left.eval() * exp.right.eval()
}

func (exp MultExp) convert() []Token {
	var v1 = exp.left.convert()
	var v2 = exp.right.convert()
	v1 = append(v1, v2...)
	v1 = append(v1, Mult)
	return v1
}

type DivExp struct {
	left  Exp
	right Exp
}

func NewDivExp(left Exp, right Exp) Exp {
	return &DivExp{left: left, right: right}
}

func (exp DivExp) eval() float64 {
	return exp.right.eval() / exp.left.eval()
}

func (exp DivExp) convert() []Token {
	var v2 = exp.left.convert()
	var v1 = exp.right.convert()
	v1 = append(v1, v2...)
	v1 = append(v1, Div)
	return v1
}

// //////////////////
// VM run-time
type VM struct {
	codes []Token
}

type VMRunnable interface {
	run() Token
	showRunConvert()
	convert() Exp
}

func NewVM(codes []Token) *VM {
	return &VM{codes: codes}
}

func (vm *VM) showRunConvert() {
	fmt.Println("VM code: ", show(vm.codes))
	fmt.Println("=> ", vm.run())
	fmt.Println("Exp: ", vm.convert().eval())
}

func (vm *VM) run() float64 {
	stack := []float64{}
	for _, code := range vm.codes {
		switch code {
		case One:
			stack = append(stack, 1)
		case Two:
			stack = append(stack, 2)
		case Mult:
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, left*right)
		case Plus:
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, left+right)
		case Div:
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, right/left)
		}
	}
	return stack[0]
}

func (vm *VM) convert() Exp {
	stack := []Exp{}
	for _, code := range vm.codes {
		switch code {
		case One:
			stack = append(stack, NewIntExp(1))
		case Two:
			stack = append(stack, NewIntExp(2))
		case Mult:
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, NewMultExp(left, right))
		case Plus:
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, NewPlusExp(left, right))
		case Div:
			var left = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			var right = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, NewDivExp(left, right))
		}
	}
	return stack[0]
}

////////////////////
// Examples

func testVM() {
	{
		code := []Token{One, Two, Two, Plus, Mult}
		vm := NewVM(code)
		vm.showRunConvert()
	}

	{
		code := []Token{One, Two, Plus, Two, Mult}
		vm := NewVM(code)
		vm.showRunConvert()
	}
	{
		code := []Token{One, Two, Div, Two, Mult, Two, Div, Two, Div}
		vm := NewVM(code)
		vm.showRunConvert()
	}
	{
		code := []Token{Two, One, Div, Two, Mult, Two, Div, Two, Div}
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
	//  code := []Token{One, Two, Plus, Two, Mult}
	//  fmt.Println(show(code))

	//  Create a new IntExp instance
	//  exp1 := NewIntExp(2)

	// Evaluate the expression
	//  fmt.Println("Eval: ", exp1.eval())
	//  fmt.Println("Convert: ", show(exp1.convert()))

	testVM()
	//testExp()
}
