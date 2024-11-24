package generator

import (
	"math/rand"

	"project/impl/stackvm"
)

func randomIntExp(rand *rand.Rand) stackvm.Exp {
	value := rand.Intn(2) + 1
	return stackvm.NewIntExp(value)
}

func randomPlusExp(rand *rand.Rand, depth int) stackvm.Exp {
	left := RandomExp(rand, depth-1)
	right := RandomExp(rand, depth-1)
	return stackvm.NewPlusExp(left, right)
}

func randomMultExp(rand *rand.Rand, depth int) stackvm.Exp {
	left := RandomExp(rand, depth-1)
	right := RandomExp(rand, depth-1)
	return stackvm.NewMultExp(left, right)
}

func randomDivExp(rand *rand.Rand, depth int) stackvm.Exp {
	left := RandomExp(rand, depth-1)
	right := RandomExp(rand, depth-1)
	return stackvm.NewDivExp(left, right)
}

func RandomExp(rand *rand.Rand, depth int) stackvm.Exp {
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
