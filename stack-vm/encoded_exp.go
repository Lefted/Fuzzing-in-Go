package main

import (
	"fmt"
	"reflect"
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

func EncodeWithDepth(exp Exp, maxDepth, currentDepth int) ([]EncodedExp, error) {
	var tokens []EncodedExp

	var helper func(e Exp, depth int) error
	helper = func(e Exp, depth int) error {
		if depth >= maxDepth {
			switch v := e.(type) {
			case *IntExp:
				tokens = append(tokens, EncodedExp{Type: 0, Value: v.value})
			default:
				return fmt.Errorf("unexpected non-terminal at max depth: %v", reflect.TypeOf(e))
			}
			return nil
		}

		switch v := e.(type) {
		case *IntExp:
			tokens = append(tokens, EncodedExp{Type: 0, Value: v.value})
			// Add padding for remaining depth
			padding := calc_padding(maxDepth - depth) // TODO check this, maybe needs -1 but I doubt it
			for i := 0; i < padding; i++ {
				tokens = append(tokens, EncodedExp{Type: 4})
			}

		case *PlusExp:
			tokens = append(tokens, EncodedExp{Type: 1})
			if err := helper(v.left, depth+1); err != nil {
				return err
			}
			if err := helper(v.right, depth+1); err != nil {
				return err
			}

		case *MultExp:
			tokens = append(tokens, EncodedExp{Type: 2})
			if err := helper(v.left, depth+1); err != nil {
				return err
			}
			if err := helper(v.right, depth+1); err != nil {
				return err
			}

		case *DivExp:
			tokens = append(tokens, EncodedExp{Type: 3})
			if err := helper(v.left, depth+1); err != nil {
				return err
			}
			if err := helper(v.right, depth+1); err != nil {
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

func Decode(tokens []EncodedExp, maxDepth int) (Exp, error) {
	var parse func(*int, int) (Exp, error)

	parse = func(pos *int, currentDepth int) (Exp, error) {
		if *pos >= len(tokens) {
			return nil, fmt.Errorf("unexpected end of tokens")
		}

		token := tokens[*pos]

		*pos++ // Consume the current token

		switch token.Type {
		case 0: // Handle terminal nodes (IntExp)
			// ensure that the value is 1 or 2
			if token.Value != 1 && token.Value != 2 {
				return nil, fmt.Errorf("invalid value for IntExp: %v", token.Value)
			}

			// Create a terminal node (IntExp)
			intExp := &IntExp{value: token.Value}

			// Consume padding for remaining depth
			expectedPadding := calc_padding(maxDepth - currentDepth) // todo check this, maybe needs -1 but I doubt it
			for i := 0; i < expectedPadding; i++ {
				if *pos >= len(tokens) || tokens[*pos].Type != 4 {
					return nil, fmt.Errorf("unexpected token at padding position %v", *pos)
				}
				*pos++ // Consume the DontCare token
			}

			return intExp, nil

		case 1, 2, 3: // Handle non-terminal nodes (PlusExp, MultExp, DivExp)
			if currentDepth >= maxDepth {
				return nil, fmt.Errorf("unexpected non-terminal token at max depth")
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
				return &PlusExp{left: left, right: right}, nil
			case 2:
				return &MultExp{left: left, right: right}, nil
			case 3:
				return &DivExp{left: left, right: right}, nil
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
	if pos != len(tokens) {
		return nil, fmt.Errorf("unused tokens remaining: %d", len(tokens)-pos)
	}

	return exp, nil
}
