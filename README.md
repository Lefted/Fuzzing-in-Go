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

## Top-down parser for regular expressions from Softwareproject
[Slide](https://sulzmann.github.io/SoftwareProjekt/lec-cpp-advanced-syntax.html#(5))

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
