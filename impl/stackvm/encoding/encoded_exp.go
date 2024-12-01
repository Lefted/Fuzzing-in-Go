package encoding

import (
	"fmt"
	"reflect"

	"project/impl/stackvm"
)

type ExpType int

const (
	TokenInt      ExpType = 0
	TokenPlus     ExpType = 1
	TokenMult     ExpType = 2
	TokenDiv      ExpType = 3
	TokenDontCare ExpType = 4
)

type EncodedExp struct {
	Type  int
	Value int // Only used for TokenInt
}

func EncodeWithDepth(exp stackvm.Exp, maxDepth, currentDepth int) ([]EncodedExp, error) {
	var tokens []EncodedExp

	var helper func(e stackvm.Exp, depth int) error
	helper = func(e stackvm.Exp, depth int) error {
		if depth >= maxDepth {
			switch v := e.(type) {
			case *stackvm.IntExp:
				tokens = append(tokens, EncodedExp{Type: 0, Value: v.Value})
			default:
				return fmt.Errorf("unexpected non-terminal at max depth: %v", reflect.TypeOf(e))
			}
			return nil
		}

		switch v := e.(type) {
		case *stackvm.IntExp:
			tokens = append(tokens, EncodedExp{Type: 0, Value: v.Value})
			// Add padding for remaining depth
			padding := calc_padding(maxDepth - depth)
			for i := 0; i < padding; i++ {
				tokens = append(tokens, EncodedExp{Type: 4})
			}

		case *stackvm.PlusExp:
			tokens = append(tokens, EncodedExp{Type: 1})
			if err := helper(v.Left, depth+1); err != nil {
				return err
			}
			if err := helper(v.Right, depth+1); err != nil {
				return err
			}

		case *stackvm.MultExp:
			tokens = append(tokens, EncodedExp{Type: 2})
			if err := helper(v.Right, depth+1); err != nil {
				return err
			}
			if err := helper(v.Right, depth+1); err != nil {
				return err
			}

		case *stackvm.DivExp:
			tokens = append(tokens, EncodedExp{Type: 3})
			if err := helper(v.Left, depth+1); err != nil {
				return err
			}
			if err := helper(v.Right, depth+1); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unsupported expression type: %v", reflect.TypeOf(e))
		}
		return nil
	}

	if err := helper(exp, currentDepth); err != nil {
		return nil, err
	}

	return tokens, nil
}

func calc_padding(remaining_layers int) int {
	if remaining_layers <= 0 {
		return 0
	}

	total_padding := 0
	for i := 0; i < remaining_layers; i++ {
		total_padding += 2 << i
	}
	return total_padding
}

func Decode(tokens []EncodedExp, maxDepth int, resilient bool) (stackvm.Exp, error) {
	var parse func(*int, int) (stackvm.Exp, error)

	parse = func(pos *int, currentDepth int) (stackvm.Exp, error) {
		if *pos >= len(tokens) {
			return nil, fmt.Errorf("unexpected end of tokens")
		}

		token := tokens[*pos]

		*pos++ // Consume the current token

		switch token.Type {
		case 0: // Handle terminal nodes (IntExp)

			// ensure that the value is 1 or 2
			if token.Value != 1 && token.Value != 2 {
				if resilient {
					token.Value = 1
				} else {
					return nil, fmt.Errorf("invalid value for IntExp: %v", token.Value)
				}
			}

			// Create a terminal node (IntExp)
			intExp := &stackvm.IntExp{Value: token.Value}

			// Consume padding for remaining depth
			expectedPadding := calc_padding(maxDepth - currentDepth)
			for i := 0; i < expectedPadding; i++ {

				if *pos >= len(tokens) || tokens[*pos].Type != 4 {
					if !resilient {
						return nil, fmt.Errorf("unexpected token at padding position %v", *pos)
					}
				}
				*pos++ // Consume the DontCare token
			}

			return intExp, nil

		case 1, 2, 3: // Handle non-terminal nodes (PlusExp, MultExp, DivExp)
			if currentDepth >= maxDepth {
				if resilient {
					return &stackvm.IntExp{Value: 1}, nil
				} else {
					return nil, fmt.Errorf("unexpected non-terminal token at max depth")
				}
			}

			// Create an operator node
			// Parse left and right subtrees recursively
			left, err := parse(pos, currentDepth+1)
			if err != nil {
				return nil, err
			}

			right, err := parse(pos, currentDepth+1)
			if err != nil {
				return nil, err
			}

			// Build the appropriate expression node
			switch token.Type {
			case 1:
				return &stackvm.PlusExp{Left: left, Right: right}, nil
			case 2:
				return &stackvm.MultExp{Left: left, Right: right}, nil
			case 3:
				return &stackvm.DivExp{Left: left, Right: right}, nil
			}
		}
		return nil, fmt.Errorf("unknown token type to decode: %v", token.Type)
	}

	// Start parsing from the first token
	pos := 0
	exp, err := parse(&pos, 0)
	if err != nil {
		return nil, err
	}

	// Ensure all tokens were consumed
	if pos != len(tokens) && !resilient {
		return nil, fmt.Errorf("unused tokens remaining: %d", len(tokens)-pos)
	}

	return exp, nil
}
