package randoms

import (
	"errors"
	"math/rand"

	"golang.org/x/exp/constraints"
)

var ErrNegativeWeight = errors.New("negative weight")
var ErrZeroWeight = errors.New("zero weight")

// 加权随机选择
func WeightedRandomSelect[T constraints.Integer](weights []T) (int, error) {
	var totalWeight int
	for _, weight := range weights {
		// 负权重是不允许的
		if weight < 0 {
			return 0, ErrNegativeWeight
		}
		totalWeight += int(weight)
	}

	// 总权重为 0 也是不允许的
	if totalWeight <= 0 {
		return 0, ErrZeroWeight
	}

	r := rand.Intn(totalWeight)

	for i, weight := range weights {
		totalWeight -= int(weight)
		if r >= totalWeight {
			return i, nil
		}
	}
	return len(weights) - 1, nil
}
