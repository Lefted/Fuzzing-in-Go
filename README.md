# Fuzzing in Go

## Main Repo

https://github.com/sulzmann/Seminar/blob/main/winter24-25.md

# Outline

1. Seminar Preparation - Topic Selection for T1
   - Stack VM
   - Json Path Parser
   - Regular Expression Parser
2. Meeting 1 - Implementation of the Stack VM
   - Introduction
   - Implementation
   - Some caveats about go that occurred during the implementation
3. Meeting 2 - Fuzzing the Stack VM
   - Types
   - Div operator
   - Generators
   - Fuzzing with Generators
   - Fuzzing with the guided fuzzer
4. Meeting 3 - Improving the Fuzzing
   - Resilient decoding
   - Summary

# Seminar Preparation

## Topic Selection

### Interested in topic T1: Fuzzing in Go

> 1. Build your own "larger" Go application. For example, you could (re)program some of the `Softwareprojekt' exercises in Go.
>
> 2. Introduce some bugs
>
> 3. See how effective fuzzing is for bug finding
>
> 4. Report your experiences

The following outline possible projects that could be used for this topic.

### Stack-based virtual machine from Softwareproject

The original code from Softwareproject can be found on [these Slides](https://sulzmann.github.io/SoftwareProjekt/lec-cpp-advanced-vm.html).
And if you want to run it you can use [OnlineGDB](https://www.onlinegdb.com/).

Idea:

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

### Top-down parser for regular expressions from Softwareproject

The original code from Softwareproject can be found on [these Slides](<https://sulzmann.github.io/SoftwareProjekt/lec-cpp-advanced-syntax.html#(5)>).
And if you want to run it you can use [OnlineGDB](https://www.onlinegdb.com/).

Idea:

- Generate random expressions
  - Generate with random operators and symbols
  - Begin with valid expressions. Then mutate by randomly inserting operators or symbols
  - Control number of nested expressions
  - Control length
- Compare with other parsers

### Json Path Parser

Idea:

Build my own json path parser project and then fuzz test it.

**Core-functionality**

1. Basic Path Navigation
   - Root symbol (`$`) to define the starting point of the JSON path
   - Dot Notation (`.`) for accessing keys or poperies in a JSON object. E.g. `$.store.books`
   - Bracket Notion (`[]`) for accessing properties wih dynamic or non-standard keys and for handling arrays. E.g. `$.store.books[0]`
2. Wildcard Selection
   - For properties (`*`) to select all properties in an object.
     E.g. `$-store.*.id`
   - For arrays (`[*]`) to select all elements of an array.
     E.g. `$-store.books[*]`
3. Recursive Descent
   - For traversing all levels of the JSON tree to search for a property
     E.g. `$..author` to retrieve all author fields in the document
4. Filtering?
   - Filter Operator `?()` to only return elements matching a specific condition
     E.g. `$.store.books[?(@.price < 10)]`
   - Existence Operator `@` which refers to the current node
     E.g. `$.store.books[?(@.isbn)]` to return all books with an isbn

**Optional-functionality**

1. Slice Operation `[start:end]` to select a specific range of entries from an array
   E.g. `$.store.books[0:2]`

**Fuzz testing ideas**

- Generate random JSONs using an existing generator
- Generate json paths
  - Generate json paths using the basic syntax
  - Generate malformed paths
    - Unclosed brackets
    - Missing dots
    - Invalid characters
  - Edge cases
    - Empty paths
    - Extremely long paths
    - Large indices
    - Deeply nested recursion

# Meeting 1 - Implementation of the Stack VM

We decided to use the stack-based virtual machine as our test project.
First we will implement the stack-based virtual machine in Go.
Then we will introduce a bug in the implementation.
Finally, we will use the fuzzer to find the bug.

## Introduction

The virtual machine can use the Reverse Polish Notation. This means we write the operands first and then the operation to be performed.

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

using recursion it will be easy to convert an expression to a tree.

## Implementation

The original source code be found at [sulzmann.github.io](<https://sulzmann.github.io/SoftwareProjekt/lec-cpp-advanced-vm.html#(6)>).

Our go implementation is located in `impl/stackvm/stack_vm.go`.

### Some caveats about go that occurred during the implementation

- Usage of `ioata` instead of enums
- No function overloading

# Meeting 2 - Fuzzing the Stack VM

## Types

Our previous code used `int` to represent the tokens.
We can improve our code by using types instead.

To do this we can use the `type` keyword to define a new types for our tokens.
Then we can use these types in our functions instead of `int`.

```go
// File: impl/stackvm/stack_vm.go

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

The compiler will now check that we are using the correct types in our functions.
This will help us to avoid bugs where we accidentally use the wrong token.

For the following code :

```go
	var token stackvm.Token
	var value int
	token = stackvm.Plus
	value = token
```

The compiler will give us an error for the last line:

```
cannot use token (variable of type stackvm.Token) as int value in assignment
```

### Some caveats

It would still be possible to call the function with an `int` instead of a `Token`.
The compiler will not complain about `showToken(5)`. Even though `5` is not a valid token.

This also means that we still need to expect `int` values in our functions.
Our switch case from above can therefore not be exhaustive.

## Div operator

In order to demonstrate a more complex bug we introduce a new operator `Div`.
The `Div` operator will divide the second operand by the first operand.

```go
// File: impl/stackvm/stack_vm.go

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

Our virtual machine also needs to be updated to handle the new operator:

```go
// File: impl/stackvm/stack_vm.go

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

Later we will introduce a bug in the `Div` operator to demonstrate the fuzzer.

## Generators

Generators are a concept of haskell's quickcheck library.
They make it easy to generate random test data for more complex data structures.
An explanation can be found in the [docs](https://hackage.haskell.org/package/QuickCheck-2.15.0.1/docs/Test-QuickCheck-Gen.html)
or in [this quickcheck paper](https://www.cs.tufts.edu/~nr/cs257/archive/john-hughes/quick.pdf).
To get an idea how generators could work in C++ see [these slides](<https://sulzmann.github.io/SoftwareProjekt/lec-cpp-advanced-quick-check.html#(1)>).

We will try to implement a quickcheck-like generator for our expressions.

```go
// File: impl/stackvm/generator/exp_generator.go

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

## Fuzzing with Generators

We can use the generator to create a random expression and then run the test case with the generated expression.

```go
func FuzzWithGenerator(f *testing.F) {
	// No need to add seed inputs as the generator will generate random expressions

	f.Fuzz(func(t *testing.T, seed int) {
		// use seed for reproducibility
		rand := rand.New(rand.NewSource(int64(seed)))

		// use generator to create a random expression
		exp := gen.RandomExp(rand, 3)

		// ... run the test case
	})
}
```

**The Bug**

The bug is that any `Div` expression will not handle left child expression other than 'Int' expressions.

```go
// File: impl/stackvm/stack_vm.go

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
```

```
Works:
                           DIV
                         /     \
                        1       2

Fails:
                           DIV
                         /     \
                      PLUS      2
                     /    \
		    1      2
```

The second example will run into our bug since the left child is not an `Int` expression.

### Why not simply swap the operands?

Initially the idea was to swap the operands in the `Div` operator.
But since we want to compare the result with the results from using the guided fuzzer we decided to introduce a bug that is not as easy to find.

Here's why:

1. The guided fuzzer receives a seed of randomly generated expressions
2. It will then first try to find any bug using the provided seed inputs
3. After that it will try to mutate inputs to find new bugs

Swapping the operands:

- The bug is likely to be found with the initial seed inputs
  <br> → The guided fuzzer won't have anything to do

Our complex bug:

- We can only provide randomly generated expressions with a depth of 1 as seed input
  <br> → The seed input will not trigger the bug
  <br> → The guided fuzzer will have to 'find out' that it can use a depth of 2 to find the bug

## Testing using the generators

**Command**

```bash
go test -timeout 30s -run ^FuzzWithGenerator$ project/impl/stackvm -test.v --fuzz=FuzzWithGenerator
```

**Output**

```
=== RUN   FuzzWithGenerator
warning: starting with empty corpus
fuzz: elapsed: 0s, execs: 0 (0/sec), new interesting: 0 (total: 0)
fuzz: elapsed: 0s, execs: 1 (4/sec), new interesting: 0 (total: 0)
--- FAIL: FuzzWithGenerator (0.27s)
    --- FAIL: FuzzWithGenerator (0.00s)
        stack_vm_test.go:67: 1 2 + 1 1 * / 2 *
        stack_vm_test.go:68: Result from VM: 6
        stack_vm_test.go:69: Result from Expression: 0
        stack_vm_test.go:70: Mismatch: original evaluation = 0, VM evaluation = 6

    Failing input written to testdata\fuzz\FuzzWithGenerator\122121052aa033c3
    To re-run:
    go test -run=FuzzWithGenerator/122121052aa033c3
=== NAME
FAIL
exit status 1
FAIL    project/impl/stackvm    0.664s
```

The bug was found in `0.27s`. This is already a very good result.
<br>
The downside is that the _go fuzzer_ will not be able to guide the fuzzing process. This is because the fuzzer can only control the seed and not the generated expressions. This means that the fuzzer will not be able to explore the coverage space systematically as it can not guide the generation of the expressions.
<br>
Furthermore there is no process where the faulty test input is minimized. Thus we as a developer are left with the input `1 2 + 1 1 * / 2 *`. We might need to spend some time to understand what is the problem with this input.

## Fuzzing with the guided fuzzer

Since the built in go fuzzer does not support fuzzing structs, we will use the [_FuzzPlus_](https://github.com/MaxiLambda/go-seminar) implementation from the seminar by Maximilian Lincks.

### Initial idea

We generate random expression with our previously implemented generator. Then we can use _FuzzPlus_ to create a corpus of these expressions.
FuzzPlus will flatten our expressions into an array of 'fuzzable' values. The go fuzzer will mutate the array of values in each iteration and FuzzPlus will convert it back to our expression.
Our test function will then receive the expression and we can run the test case:

```
exp struct -> flattened array -> mutated array -> exp struct
```

```go
func FuzzTest(f *testing.F) {
	ff := fuzzplus.NewFuzzPlus(f)

	rand := rand.New(rand.NewSource(1))

	// create a test corpus
	for i := 0; i < 500; i++ {
		exp := gen.RandomExp(rand, 2) // max-depth 2
		ff.Add2(exp)
	}

	ff.Fuzz(func(t *testing.T, in stackvm.Exp) {
		// in: the generated expression

		// .. run the test case
	})
}
```

### Problem

Let's say we have the following expressions `A` and `B`:

```
A:
                       +
                      /  \
                     2    2

B:
		       +
                      /  \
                     2    *
                          |  \
                          1   2
```

If we flatten these expressions we get something similar to the following arrays:

```
A:
[+, 2, 2]

B:
[+, 2, *, 1, 2]
```

As you can see, the arrays have different lengths. This is a problem because the go fuzzer requires the corpus to have the same length in each iteration. This also means that the generated expression we receive as input will have the same length in each iteration.

### Solution

Before handing our expressions to FuzzPlus we will convert them into a middleman format that guarantees a fixed length. We will receive the middleman format as test input in each iteration. Thus before we run our test case we will convert the middleman format back to our expression.

```
exp struct -> fixed size format -> flattened array -> mutated array -> fixed size format -> exp struct
```

```go
// File: impl/stackvm/stack_vm_test.go

func FuzzPlusExp(f *testing.F) {
	ff := fuzzplus.NewFuzzPlus(f)

	rand := rand.New(rand.NewSource(1))

	// create a test corpus
	for i := 0; i < 500; i++ {
		exp := gen.RandomExp(rand, 2) // max-depth 2
		encodedExp, err := encoding.EncodeWithDepth(exp, 3, 0) // fixed-size-depth: 3
		if err != nil {
			f.Fatalf("Failed to encode expression: %v", err)
		}
		ff.Add2(encodedExp)
	}

	ff.Fuzz(func(t *testing.T, in []encoding.EncodedExp) {
		// decode the fixed-size format
		expression, err := encoding.Decode(in, 3)
		if err != nil {
			return
		}

		// ... run the test case
	})
}
```

**Encoding**

Convert an expressions into a list of tokens and insert padding to make the list a fixed size.

For our example expressions `A` and `B` and a fixed-depth of 3 we can encode them as follows:

```
A:

                       +
                      /  \
                     2    2

Encoded:
                           +
                        /      \
                     2           2
                   /     |       |    \
		  x      x      x      x
	         / \    / \    / \    / \
	         x x    x x    x x    x x
B:
		       +
	              /  \
                     2    *
                          |  \
                          1   2

Encoded:
                           +
                        /      \
                     2           *
                   /     |       |    \
		  x      x      1      2
	         / \    / \    / \    / \
                 x x    x x    x x    x x

```

If we now flatten these expressions we get something similar to the following arrays:

```
A:
[+, 2, x, x, x, x, x, x, 2, x, x, x, x, x, x]

B:
[+, 2, x, x, x, x, x, x, *, 1, x, x, 2, x, x]
```

For our middleman format `x` will be a `DontCare` token. We will use it as padding to make the list a fixed size.

Example of encoding `A` as list of basic tokens:

```
[
 {Type: Plus, Value: 0}
 {Type: Int, Value: 2}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: Int, Value: 2}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
 {Type: DontCare, Value: 0}
]
```

This solves the problem of having expressions with different lengths.

### Testing it

**Command**

```bash
go test -timeout 30s -run ^FuzzPlusExpNonResilient$ project/impl/stackvm -test.v --fuzz=FuzzPlusExpNonResilient
```

**Output**

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

As we can see, the fuzzer found the bug. But it took `16m 15s` to find it. This may be because we only provided seeds for the fuzzer with a depth of 1. This means that the fuzzer had to find that it could use a depth of 2 in the expression to find the bug.

The test input that caused the bug is written to a file:

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

But to be able to understand it we need to convert it back to our expression. Which can be tedious:

```
[0] = stack-vm/stack-vm.EncodedExp {Type: 3, Value: 12}
[1] = stack-vm/stack-vm.EncodedExp {Type: 2, Value: 2}
[2] = stack-vm/stack-vm.EncodedExp {Type: 0, Value: 1}
[3] = stack-vm/stack-vm.EncodedExp {Type: 0, Value: 1}
[4] = stack-vm/stack-vm.EncodedExp {Type: 0, Value: 1}
[5] = stack-vm/stack-vm.EncodedExp {Type: 4, Value: -66}
[6] = stack-vm/stack-vm.EncodedExp {Type: 4, Value: 0}
```

To make this process easier it is best to add logging in the test case for the failing input. This way we can easily see the failing input in the console.
We can also do this later on and run the test with the same input again:

`go test -run=FuzzPlusExpNonResilient/214440cc69e9949e`

Now it is easy to understand that the failing expression was `1 1 1 * /`.
This is a pretty minimal expression and it is easy to see that the result should be `1` and not `0`.

# Meeting 3 - Improving the Fuzzing

## Resilient decoding

In our previous run we found that the fuzzer was able to find the bug but it took a long time to do so.
We can also see that it took `46925088` executions to find it.
<br>
This may be because during decoding of the fixed-size format we returned early if we encountered an error. For example if the fuzzer inserted an invalid amount of `DontCare` tokens or if the value of an `Int` token was invalid.

A more resilient decoding would be to continue decoding even if we encounter an error. This way we can still run the test case and the fuzzer can continue to mutate the input. This may allow the fuzzer to find the bug faster.

In `resilient` mode we will:

- interpret any invalid `Int` value as `1`
- if we encounter a `DontCare` token skip the next _n_ tokens instead of insisting that the next _n_ tokens are valid
- not insist that all tokens must be consumed
- interpret non-terminal tokens at max depth as 'Int' with value '1'

Let's check if this can help the fuzzer to find the bug faster.

**Command**

```bash
go test -timeout 30s -run ^FuzzPlusExpResilient$ project/impl/stackvm -test.v --fuzz=FuzzPlusExpResilient
```

**Output**

```

fuzz: elapsed: 0s, gathering baseline coverage: 0/36 completed
fuzz: elapsed: 0s, gathering baseline coverage: 36/36 completed, now fuzzing with 8 workers
fuzz: elapsed: 3s, execs: 101343 (33757/sec), new interesting: 0 (total: 36)
fuzz: elapsed: 4s, execs: 139145 (38269/sec), new interesting: 0 (total: 36)
--- FAIL: FuzzPlusExpResilient (4.19s)
--- FAIL: FuzzPlusExpResilient (0.00s)
stack_vm_test.go:169: 1 1 1 + /
stack_vm_test.go:170: Result from VM: 0.5
stack_vm_test.go:171: Result from Expression: 0
stack_vm_test.go:172: Mismatch: original evaluation = 0, VM evaluation = 0.5

    Failing input written to testdata\fuzz\FuzzPlusExpResilient\8fce845b0985bb2a
    To re-run:
    go test -run=FuzzPlusExpResilient/8fce845b0985bb2a

=== NAME
FAIL
exit status 1
FAIL project/impl/stackvm 4.447s
```

Huh? The fuzzer found the bug in **4.19s**. This is a lot faster than the 16 minutes it took before. But was this just luck? Let's run the fuzzer again to see if it can find the bug again.

**Output 2**

```

=== RUN FuzzPlusExpResilient
...
fuzz: elapsed: 0s, gathering baseline coverage: 0/36 completed
fuzz: elapsed: 0s, gathering baseline coverage: 36/36 completed, now fuzzing with 8 workers
fuzz: elapsed: 2s, execs: 33400 (14271/sec), new interesting: 0 (total: 36)
--- FAIL: FuzzPlusExpResilient (2.55s)
--- FAIL: FuzzPlusExpResilient (0.00s)
stack_vm_test.go:163: 1 1 1 + /
stack_vm_test.go:164: Result from VM: 0.5
stack_vm_test.go:165: Result from Expression: 0
stack_vm_test.go:166: Mismatch: original evaluation = 0, VM evaluation = 0.5

    Failing input written to testdata\fuzz\FuzzPlusExpResilient\01ea708fd5625eb9
    To re-run:
    go test -run=FuzzPlusExpResilient/01ea708fd5625eb9

=== NAME
FAIL
exit status 1
FAIL project/impl/stackvm 3.114s

```

Again the same bug was found in **2.55s**. This seems to indicate that the fuzzer was not just lucky in the first run.
Interestingly both times the fuzzer found the bug with the same input.

## Pros and Cons of Generators vs Guided Fuzzing

|      | Generators                                                                | Guided Fuzzing                                                                                                                                                    |
| ---- | ------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Pros | - easy to implement <br> - very fast                                      | - guided approach observing code coverage may find more edge-cases <br> - tries to minimize the failing input                                                     |
| Cons | - only as good as randomness <br> - failing inputs are likely to be large | - cumbersome to implement for complex data-structures because structs are not natively supported <br> - not as fast for complex data-structures due to 'decoding' |

## Summary

What can we learn from this?

Using the _go fuzzer_ with for complex structure like our expressions is not trivial.
With some tinkering using a middleman format and _FuzzPlus_ we were able to use the _go fuzzer_ to find the bug in our implementation.
<br>
It is important to convert the middleman format back to the original structure in a resilient way.
This way the fuzzing process doesn't waste time on invalid inputs.
<br>
One may also see differences in time using different amount of seeds.
The difficulty lies in finding the right balance between the amount of seeds.
Using too many seeds may lead to overfitting where the fuzzer may not divate from the structure of the seeds enough.
Using too few seeds may lead to underfitting where the fuzzer may not recognize the general structure of the inputs.
<br>
But for this project I didn't see a significant difference in time using 1 or 500 seeds for the corpus.

Even though the guided fuzzing eventually worked fine using our own generators to generate random expressions was still faster and most importantly easier.
If writing generators seems more reasonable than creating an encoding/decoding step for the fixed-size format, it is probably the better choice.

Nevertheless, the _go fuzzer_ is a powerful tool and can be used to find bugs in complex structures.
<br>
I could imagine that it could maybe find more edge case bugs. For example in our expression project it could maybe decide to do addition operations until we reach the maximum float value.
This would be hard to find with a random generator since the probability of generating such a large expression approaches zero.
Especially if we also had a subtraction operation.
Then on average we would generate an equal amount of addition and subtraction operations and the probability of reaching the maximum float value would be even lower.
<br>
But even with the guided fuzzer such edge cases wouldn't be able to be found within a few minutes.
This process might be better suited for long fuzzing sessions using specific project versions and not for running in a CI pipeline.

Future projects could explore how to create a fixed-size format for unknown structures.
With further reflection on the encoding/decoding process it might be possible to create a more general solution that can be used for any structure.
Then using the guided fuzzer would require less steps to set up and could be used for more projects.
