package usecase

import "math"

// hungarian решает задачу назначения методом венгерского алгоритма.
// cost[i][j] — стоимость назначения курьера i на заказ j.
// Возвращает срез assignRow, где assignRow[i] = j означает: курьер i назначен на заказ j.
// Если курьеров больше чем заказов — лишние не назначаются (assignRow[i] = -1).
func hungarian(cost [][]float64) []int {
	n := len(cost)
	if n == 0 {
		return nil
	}
	m := len(cost[0])

	size := max(m, n)

	inf := math.MaxFloat64 / 2

	c := make([][]float64, size)
	for i := range c {
		c[i] = make([]float64, size)
		for j := range c[i] {
			if i < n && j < m {
				c[i][j] = cost[i][j]
			} else {
				c[i][j] = 0
			}
		}
	}

	u := make([]float64, size+1)
	v := make([]float64, size+1)
	p := make([]int, size+1)
	way := make([]int, size+1)

	for i := 1; i <= size; i++ {
		p[0] = i
		j0 := 0
		minV := make([]float64, size+1)
		used := make([]bool, size+1)

		for j := range minV {
			minV[j] = inf
		}

		for {
			used[j0] = true
			i0 := p[j0]
			delta := inf
			var j1 int

			for j := 1; j <= size; j++ {
				if !used[j] {
					cur := c[i0-1][j-1] - u[i0] - v[j]
					if cur < minV[j] {
						minV[j] = cur
						way[j] = j0
					}
					if minV[j] < delta {
						delta = minV[j]
						j1 = j
					}
				}
			}

			for j := 0; j <= size; j++ {
				if used[j] {
					u[p[j]] += delta
					v[j] -= delta
				} else {
					minV[j] -= delta
				}
			}

			j0 = j1
			if p[j0] == 0 {
				break
			}
		}

		for j0 != 0 {
			p[j0] = p[way[j0]]
			j0 = way[j0]
		}
	}

	result := make([]int, n)
	for i := range result {
		result[i] = -1
	}

	for j := 1; j <= size; j++ {
		i := p[j] - 1
		jReal := j - 1
		if i >= 0 && i < n && jReal < m {
			result[i] = jReal
		}
	}

	return result
}

func euclidean(lat1, lon1, lat2, lon2 float64) float64 {
	dlat := lat1 - lat2
	dlon := lon1 - lon2
	return math.Sqrt(dlat*dlat + dlon*dlon)
}
