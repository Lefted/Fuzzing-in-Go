# Main Repo

https://github.com/sulzmann/Seminar/blob/main/winter24-25.md

# Topic Selection

## Interested in topic T1: Fuzzing in Go

> 1. Build your own "larger" Go application. For example, you could (re)program some of the `Softwareprojekt' exercises in Go.
>
> 2. Introduce some bugs
>
> 3. See how effective fuzzing is for bug finding
>
> 4. Report your experiences

# Possible Code Examples

## Stack-based virttual machine from Softwareproject

[Slide](https://sulzmann.github.io/SoftwareProjekt/lec-cpp-advanced-vm.html)

### Idea:

- Generate random expressions
  - Randomly generate IntExpression with 1 or 2 value
  - Combine randomly
  - Control depth
- Expression Evaluation Consistency

  For each generated expression, convert it to VM code and back to an expression.
  Compare the evaluation results of the original and the converted expressions.

- VM Code Execution Consistency

  Convert generated VM code to an expression and back to VM code.
  Ensure that running both the original and the converted VM code yields the same result.

## Top-down parser for regular expressions from Softwareproject

[Slide](<https://sulzmann.github.io/SoftwareProjekt/lec-cpp-advanced-syntax.html#(5)>)

[OnlineGDB](https://www.onlinegdb.com/)

### Idea:

- Generate random expressions
  - Generate with random operators and symbols
  - Begin with valid expressions. Then mutate by randomly insering operators or symbols
  - Control number of nested expressions
  - Control length
- Compare with other parsers

## Json Path Parser

### Idea

Build my own json path parser project and then fuzz test it.

### Core-functionality

1. Basic Path Navigation
   - Root symbol (`$`) to define the starting point of the JSON path
   - Dot Notation (`.`) for accessing keys or properies in a JSON object. E.g. `$.store.books`
   - Bracket Noaion (`[]`) for accessing properties wih dynamic or non-sanadard keys and for handling arrays. E.g. `$.store.books[0]`
2. Wildcard Selection
   - For properties (`*`) to select all properties in an object.
     E.g. `$-store.*.id`
   - For arrays (`[*]`) to select all elements of an array.
     E.g. `$-store.books[*]`
3. Recursive Descent
   - For traversing all levels of the JSON tree to search for a property
     E.g. `$..auhor` to retrive all auhor fields in the document
4. Filtering?
   - Filter Operator `?()` to only return elements matching a specific condition
     E.g. `$.store.books[?(@.price < 10)]`
   - Existance Operator `@` which refers to the current node
     E.g. `$.store.books[?(@.isbn)]` to return all books with an isbn

### Optional-functionality

1. Slice Operaion `[start:end]` to select a specific range of entries from an array
   E.g. `$.store.books[0:2]`

### Fuzz-testing

- Generate random jsons using an existing generator
- Generate json paths
  - Generate json paths using the basic syntax
  - Generate malformed paths
    - Unclosed brackets
    - Missing dots
    - Invalid characters
  - Edge casaes
    - Empty pahs
    - Extremly long paths
    - Large indices
    - Deeply nested recursion

# Lets start with the Stack VM

## Introduction

We use the Reverse Polish Notation. This means we write the operands first and then the operation to be performed.

`1 * (2 + 3)`

becomes

`Push 1; Push 2; Push 3; Plus; Mult;`

Let's start with the same instruction set we have in Softwareprojekt:

- ONE: Push 1
- TWO: Push 2
- PLUS: Addition
- MULT: Multiplication

`1 * (2 + 1)`

becomes

`ONE; TWO; ONE; PLUS; MULT`

or as Syntaxtree

`Mult(Int(1), Plus(Int(2), Int(1)))`

we will use recursion to convert an expression to a tree

## Some findings

- Usage of `ioata` instead of enums
- No function overloading

# Types

We can use the `type` keyword to define a new types for our tokens. Then we can use these types in our functions instead of `int`.

```go
type Token int

const (
	Plus Token = iota
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
```

Switch case is unfortunately still not exhaustive

# Div operator

```go
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
```

```go

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
```

# Generators

```go
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
```

# Fuzzing with guided fuzzer

## Idea

Convert our expressions into a format that FuzzPlus can understand. Then we can use the fuzzer to generate random expressions and test them.

### Encoding

Convert an expressions into a list of tokens and insert padding to make the list a fixed size.

E.g. for the random expression with a depth of 1

```
                      +
                      / \
                     2  2
```

padded to a depth of 2 using DontCare tokens (x)

```
                       +
                      / \
                     2    2
                  /  |   |   \
                 x  x    x   x
```

we can encode this as

```
[
 {Type: plus, Value: 0}
 {Type: int, Value: 2}
 {Type: dontCare, Value: 0}
 {Type: dontCare, Value: 2}
 {Type: int, Value: 2}
 {Type: dontCare, Value: 0}
 {Type: dontCare, Value: 0}
]
```

Now we can:

- Generate random expressions
- Convert them to the encoded format (which will result in arrays of tokens of the same size)
- Use FuzzPlus's `Add2` function to create a corpus of encoded expressions
- Use FuzzPlus's `Fuzz` function
- Inside it: decode the encoded expression
- Run our test case

## Testing it

The bug I introduced was that any 'Div' expression would not handle child expressions other than 'Int' expressions.
This means that the expression `Div(Plus(Int(1), Int(2)), Int(2))` would not be evaluated correctly and instead return 0.

```
fuzz: elapsed: 16m0s, execs: 46156948 (51444/sec), new interesting: 1 (total: 32)
fuzz: elapsed: 16m3s, execs: 46315931 (53035/sec), new interesting: 1 (total: 32)
fuzz: elapsed: 16m6s, execs: 46475106 (53062/sec), new interesting: 1 (total: 32)
fuzz: elapsed: 16m9s, execs: 46637031 (53973/sec), new interesting: 1 (total: 32)
fuzz: elapsed: 16m12s, execs: 46792053 (51663/sec), new interesting: 2 (total: 33)
fuzz: elapsed: 16m15s, execs: 46925088 (49735/sec), new interesting: 2 (total: 33)
--- FAIL: FuzzPlusExp (974.85s)
    --- FAIL: FuzzPlusExp/214440cc69e9949e (0.00s)
        stack_vm_test.go:176: 1 1 1 * /
        stack_vm_test.go:177: Result from VM: 1
        stack_vm_test.go:178: Result from Expression: 0
        stack_vm_test.go:179: Mismatch: original evaluation = 0, VM evaluation = 1

    Failing input written to testdata\fuzz\FuzzPlusExp\214440cc69e9949e
    To re-run:
    go test -run=FuzzPlusExp/214440cc69e9949e
=== NAME
FAIL
exit status 1
FAIL    stack-vm/stack-vm       975.220s
```

As we can see, the fuzzer found the bug. But it took a long time to find it. This may be because I only provided seeds for the fuzzer with a depth of 1. This means that the fuzzer had to find that it could use a depth of 2 in the expression to find the bug.

What I noticed is that the output of the saved test is not very useful. This is what the fuzzer saves:

```
go test fuzz v1
int(3)
int(12)
int(2)
int(2)
int(0)
int(1)
int(0)
int(1)
int(0)
int(1)
int(4)
int(-66)
int(4)
int(0)
```

But it is not strictly the same we receive in our FuzzPlus function

```
[0] = stack-vm/stack-vm.EncodedExp {Type: 1, Value: 0}
[1] = stack-vm/stack-vm.EncodedExp {Type: 0, Value: 2}
[2] = stack-vm/stack-vm.EncodedExp {Type: 4, Value: 0}
[3] = stack-vm/stack-vm.EncodedExp {Type: 4, Value: 0}
[4] = stack-vm/stack-vm.EncodedExp {Type: 0, Value: 2}
[5] = stack-vm/stack-vm.EncodedExp {Type: 4, Value: 0}
[6] = stack-vm/stack-vm.EncodedExp {Type: 4, Value: 0}
```

FuzzPlus seems to reconstruct the input fine, but the original output of the fuzzer is not really readable.
Nevertheless, if we just run the function with the input again using
`go test -run=FuzzPlusExp/214440cc69e9949e`
we can inspect the input.
