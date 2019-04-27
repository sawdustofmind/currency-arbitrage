package arbalgo

import (
	"math"
)

// FloydWarshall finds shortest paths for given edges
func FloydWarshall(edges map[int]map[int]float64) *ShortestPaths {
	V := len(edges)
	dist := make([][]float64, V)
	next := make([][]*int, V)
	for u := 0; u < V; u++ {
		dist[u] = make([]float64, V)
		next[u] = make([]*int, V)
		for v := 0; v < V; v++ {
			weight, ok := edges[u][v]
			if ok {
				dist[u][v] = weight
				vert := v
				next[u][v] = &vert
				continue
			}
			if u == v {
				dist[u][v] = 0
				vert := v
				next[u][v] = &vert
			} else {
				dist[u][v] = math.Inf(1)
			}
		}
	}

LOOP:
	for k := 0; k < V; k++ {
		for i := 0; i < V; i++ {
			for j := 0; j < V; j++ {
				if dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
					next[i][j] = next[i][k]
				}
			}
		}
		for i := 0; i < V; i++ {
			if dist[i][i] < 0 {
				break LOOP
			}
		}
	}
	return &ShortestPaths{
		dist: dist,
		next: next,
	}
}

type ShortestPaths struct {
	dist [][]float64
	next [][]*int
}

func (p ShortestPaths) HasCycles() bool {
	V := len(p.dist)
	for i := 0; i < V; i++ {
		if p.dist[i][i] < 0 {
			return true
		}
	}
	return false
}

// Get returns shortest path with distance if graph has no negative cycles.
// Otherwise, path to any vert of cycle and all vertices of cycle
func (p ShortestPaths) Get(u, v int) ([]int, []int, float64) {
	dist := p.dist[u][v]
	if p.next[u][v] == nil {
		return nil, nil, 0
	}
	path := []int{u}
	if u == v {
		u = *p.next[u][v]
		if u != v {
			path = append(path, u)
		}
	}
	for u != v {
		u = *p.next[u][v]
		if pathToCycle, cycle, ok := p.checkCycle(path, u); ok {
			return pathToCycle, cycle, 0
		}
		path = append(path, u)
	}
	return path, nil, dist
}

func (p ShortestPaths) checkCycle(path []int, next int) ([]int, []int, bool) {
	for i, prev := range path {
		if prev == next {
			pathToCycle := path[:i+1]
			cycle := path[i:]
			return pathToCycle, cycle, true
		}
	}
	return nil, nil, false
}
