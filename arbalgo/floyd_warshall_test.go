package arbalgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloydWarshall(t *testing.T) {
	weights := map[int]map[int]float64{
		0: {2: -2},
		1: {0: 4, 2: 3},
		2: {3: 2},
		3: {1: -1},
	}

	type testCase struct {
		i    int
		j    int
		path []int
		dist float64
	}
	wants := []testCase{
		{0, 1, []int{0, 2, 3, 1}, -1},
		{0, 2, []int{0, 2}, -2},
		{0, 3, []int{0, 2, 3}, 0},
		{1, 0, []int{1, 0}, 4},
		{1, 2, []int{1, 0, 2}, 2},
		{1, 3, []int{1, 0, 2, 3}, 4},
		{2, 0, []int{2, 3, 1, 0}, 5},
		{2, 1, []int{2, 3, 1}, 1},
		{2, 3, []int{2, 3}, 2},
		{3, 0, []int{3, 1, 0}, 3},
		{3, 1, []int{3, 1}, -1},
		{3, 2, []int{3, 1, 0, 2}, 1},
	}

	shortestPaths := FloydWarshall(weights)
	for _, want := range wants {
		gotPath, gotDist := shortestPaths.Get(want.i, want.j)
		assert.Equal(t, want.path, gotPath)
		assert.Equal(t, want.dist, gotDist)
	}
}
