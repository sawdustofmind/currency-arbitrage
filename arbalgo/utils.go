package arbalgo

import "reflect"

func FindUniqueCycles(p *ShortestPaths, V int) [][]int {
	cycles := [][]int{}
	for i := 0; i < V; i++ {
		for j := 0; j < V; j++ {
			_, cycle, _ := p.Get(i, j)
			if len(cycle) == 0 {
				continue
			}
			cycle = ArrangeCycle(cycle)
			contains := false
			for _, c := range cycles {
				if reflect.DeepEqual(c, cycle) {
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

// ArrangeCycle sorts cycle that minimum element moves to the start of cycle path
func ArrangeCycle(cycle []int) []int {
	minI := 0
	for i, v := range cycle {
		if v < cycle[minI] {
			minI = i
		}
	}
	return append(cycle[minI:], cycle[0:minI]...)
}
