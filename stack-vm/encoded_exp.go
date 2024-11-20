package main

import (
	"fmt"
	"reflect"
)

type ExpType int

const (
	TokenInt  ExpType = 0
	TokenPlus ExpType = 1
	TokenMult ExpType = 2
	TokenDiv  ExpType = 3
)

type EncodedExp struct {
	Type  int
	Value int // Only used for TokenInt
}

func Encode(exp Exp) ([]EncodedExp, error) {
	var tokens []EncodedExp

	var helper func(e Exp) error
	helper = func(e Exp) error {
		switch v := e.(type) {
		case *IntExp:
			tokens = append(tokens, EncodedExp{Type: 0, Value: v.value})
		case *PlusExp:
			tokens = append(tokens, EncodedExp{Type: 1})
			if err := helper(v.left); err != nil {
				return err
			}
			if err := helper(v.right); err != nil {
				return err
			}
		case *MultExp:
			tokens = append(tokens, EncodedExp{Type: 2})
			if err := helper(v.left); err != nil {
				return err
			}
			if err := helper(v.right); err != nil {
				return err
			}
		case *DivExp:
			tokens = append(tokens, EncodedExp{Type: 3})
			if err := helper(v.left); err != nil {
				return err
			}
			if err := helper(v.right); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported expression type: %v", reflect.TypeOf(e))
		}
		return nil
	}

	if err := helper(exp); err != nil {
		return nil, err
	}

	return tokens, nil
}

func Decode(tokens []EncodedExp) (Exp, error) {
	var parse func(*int) (Exp, error)

	parse = func(pos *int) (Exp, error) {
		if *pos >= len(tokens) {
			return nil, fmt.Errorf("unexpected end of tokens")
		}

		token := tokens[*pos]
		*pos++ // Advance to the next token

		switch token.Type {
		case 0:
			// ensure that the value is 1 or 2
			if token.Value != 1 && token.Value != 2 {
				return nil, fmt.Errorf("invalid value for IntExp: %v", token.Value)
			}

			// Create a terminal node (IntExp)
			return &IntExp{value: token.Value}, nil

		case 1, 2, 3:
			// Create an operator node
			// Parse left and right subtrees recursively
			left, err := parse(pos)
			if err != nil {
				return nil, err
			}

			right, err := parse(pos)
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

		return nil, fmt.Errorf("unknown token type: %v", token.Type)
	}

	// Start parsing from the first token
	pos := 0
	exp, err := parse(&pos)
	if err != nil {
		return nil, err
	}

	// Ensure all tokens were consumed
	if pos != len(tokens) {
		return nil, fmt.Errorf("unused tokens remain after parsing")
	}

	return exp, nil
}
