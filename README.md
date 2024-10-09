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
* Generate random expressions
  * Randomly generate IntExpression with 1 or 2 value
  * Combine randomly
  * Control depth
* Expression Evaluation Consistency
  
  For each generated expression, convert it to VM code and back to an expression.
  Compare the evaluation results of the original and the converted expressions.
* VM Code Execution Consistency
  
  Convert generated VM code to an expression and back to VM code.
  Ensure that running both the original and the converted VM code yields the same result.

## Top-down parser for regular expressions from Softwareproject
[Slide](https://sulzmann.github.io/SoftwareProjekt/lec-cpp-advanced-syntax.html#(5))

[OnlineGDB](https://www.onlinegdb.com/)
 

 ### Idea:
 * Generate random expressions
   * Generate with random operators and symbols
   * Begin with valid expressions. Then mutate by randomly insering operators or symbols
   * Control number of nested expressions
   * Control length
 * Compare with other parsers

 ## Json Path Parser

 ### Idea

Build my own json path parser project and then fuzz test it.

### Core-functionality

1. Basic Path Navigation
    * Root symbol (`$`) to define the starting point of the JSON path
    * Dot Notation (`.`) for accessing keys or properies in a JSON object. E.g. `$.store.books`
    * Bracket Noaion (`[]`) for accessing properties wih dynamic or non-sanadard keys and for handling arrays. E.g. `$.store.books[0]`
2. Wildcard Selection
    * For properties (`*`) to select all properties in an object.
    E.g. `$-store.*.id`
    * For arrays (`[*]`) to select all elements of an array.
    E.g. `$-store.books[*]`
3. Recursive Descent
    * For traversing all levels of the JSON tree to search for a property
    E.g. `$..auhor` to retrive all auhor fields in the document
4. Filtering?
    * Filter Operator `?()` to only return elements matching a specific condition
    E.g. `$.store.books[?(@.price < 10)]`
    * Existance Operator `@` which refers to the current node
    E.g. `$.store.books[?(@.isbn)]` to return all books with an isbn

### Optional-functionality
1. Slice Operaion `[start:end]` to select a specific range of entries from an array
E.g. `$.store.books[0:2]`

### Fuzz-testing

* Generate random jsons using an existing generator
* Generate json paths
  * Generate json paths using the basic syntax
  * Generate malformed paths
    * Unclosed brackets
    * Missing dots
    * Invalid characters
  * Edge casaes
    * Empty pahs
    * Extremly long paths
    * Large indices
    * Deeply nested recursion