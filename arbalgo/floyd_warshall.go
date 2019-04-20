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

	for k := 0; k < V; k++ {
		for i := 0; i < V; i++ {
			for j := 0; j < V; j++ {
				if dist[i][j] > dist[i][k]+dist[k][j] {
					dist[i][j] = dist[i][k] + dist[k][j]
					next[i][j] = next[i][k]
				}
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

func (p ShortestPaths) Get(u, v int) ([]int, float64) {
	dist := p.dist[u][v]
	if p.next[u][v] == nil {
		return nil, 0
	}
	path := []int{u}
	for u != v {
		u = *p.next[u][v]
		path = append(path, u)
	}
	return path, dist
}
