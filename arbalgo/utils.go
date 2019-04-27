package arbalgo

import (
	"reflect"
	"sort"
)

func FindUniqueCycles(p *ShortestPaths, V int) [][]int {
	cycles := [][]int{}
	for i := 0; i < V; i++ {
		for j := 0; j < V; j++ {
			_, cycle, _ := p.Get(i, j)
			if len(cycle) == 0 {
				continue
			}
			contains := false
			for _, c := range cycles {
				if SameCycles(c, cycle) {
					contains = true
					break
				}
			}
			if !contains {
				cycles = append(cycles, cycle)
			}
		}
	}
	return cycles
}

func SameCycles(first []int, second []int) bool {
	firstCopy := first[:]
	sort.Ints(firstCopy)
	secondCopy := second[:]
	sort.Ints(secondCopy)
	return reflect.DeepEqual(firstCopy, secondCopy)
}
