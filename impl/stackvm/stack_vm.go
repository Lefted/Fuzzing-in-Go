package stackvm

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

func Show(code []Token) string {
	var builder strings.Builder
	for _, token := range code {
		builder.WriteString(showToken(token) + " ")
	}
	return builder.String()
}

////////////////////
// Expressions

type Exp interface {
	Eval() float64    // evaluate the expression (interpreter)
	Convert() []Token // convert to "reverse polish notation" (compiler)
}

type IntExp struct {
	Value int
}

func NewIntExp(x int) Exp {
	if x == 1 || x == 2 {
		return &IntExp{Value: x}
	} else {
		fmt.Printf("Invalid value: %d. Must be 1 or 2\n", x)
		return &IntExp{Value: 1}
	}
}

func (exp *IntExp) Eval() float64 {
	return float64(exp.Value)
}

func (exp *IntExp) Convert() []Token {
	n := One
	if exp.Value == 2 {
		n = Two
	}
	return []Token{n}
}

type PlusExp struct {
	Left  Exp
	Right Exp
}

func NewPlusExp(left Exp, right Exp) Exp {
	return &PlusExp{Left: left, Right: right}
}

func (exp *PlusExp) Eval() float64 {
	return exp.Left.Eval() + exp.Right.Eval()
}

func (exp *PlusExp) Convert() []Token {
	var v1 = exp.Left.Convert()
	var v2 = exp.Right.Convert()
	v1 = append(v1, v2...) // append v2 to v1 ... unpacks the slice
	v1 = append(v1, Plus)
	return v1
}

type MultExp struct {
	Left  Exp
	Right Exp
}

func NewMultExp(left Exp, right Exp) Exp {
	return &MultExp{Left: left, Right: right}
}

func (exp MultExp) Eval() float64 {
	return exp.Left.Eval() * exp.Right.Eval()
}

func (exp MultExp) Convert() []Token {
	var v1 = exp.Left.Convert()
	var v2 = exp.Right.Convert()
	v1 = append(v1, v2...)
	v1 = append(v1, Mult)
	return v1
}

type DivExp struct {
	Left  Exp
	Right Exp
}

func NewDivExp(left Exp, right Exp) Exp {
	return &DivExp{Left: left, Right: right}
}

func (exp DivExp) Eval() float64 {
	// == BUG
	switch exp.Left.(type) {
	case *IntExp:
		// do nothing
	default:
		fmt.Println("Bug hit. Left exp is: ", exp.Left)
		return 0
	}
	// ==

	return exp.Right.Eval() / exp.Left.Eval()
}

func (exp DivExp) Convert() []Token {
	var v2 = exp.Left.Convert()
	var v1 = exp.Right.Convert()
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
	Run() Token
	ShowRunConvert()
	Convert() Exp
}

func NewVM(codes []Token) *VM {
	return &VM{codes: codes}
}

func (vm *VM) ShowRunConvert() {
	fmt.Println("VM code: ", Show(vm.codes))
	fmt.Println("=> ", vm.Run())
	fmt.Println("Exp: ", vm.Convert().Eval())
}

func (vm *VM) Run() float64 {
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

func (vm *VM) Convert() Exp {
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
