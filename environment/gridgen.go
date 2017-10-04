package environment

import (
	"math"
	"math/rand"
	"sort"
	"time"
)

func diamondSquare(minSize int) [][]float64 {
	var exp uint = 0
	size := 1
	for size < minSize {
		exp++
		size = 1 << exp
	}
	size += 1

	grid := make([][]float64, size)
	for i := range grid {
		grid[i] = make([]float64, size)
	}
	rand.Seed(time.Now().UnixNano())
	randomCorner := 2*rand.Float64() - 1
	grid[0][0] = randomCorner
	grid[0][size-1] = randomCorner
	grid[size-1][0] = randomCorner
	grid[size-1][size-1] = randomCorner

	randRange := 2.0
	for subSize:=size-1; subSize > 1; subSize/=2 {
		// diamond step
		for i0 := 0; i0 < size-1; i0 += subSize {
			for j0 := 0; j0 < size-1; j0 += subSize {
				total := grid[i0][j0] +
					grid[i0][j0+subSize] +
					grid[i0+subSize][j0] +
					grid[i0+subSize][j0+subSize]
				r := randRange*(rand.Float64() - 0.5)
				grid[i0+subSize/2][j0+subSize/2] = total/4 + r
			}
		}

		randRange /= math.Sqrt2

		// square step
		for i0 := 0; i0 < size; i0 += subSize {
			for j0 := 0; j0 < size; j0 += subSize {
				if i0 < size-1 {
					jLeft := j0 - subSize/2
					if jLeft < 0 {
						jLeft += size - 1
					}
					jRight := j0 + subSize/2
					if jRight >= size {
						jRight -= size - 1
					}

					total := grid[i0][j0] +
						grid[i0+subSize][j0] +
						grid[i0+subSize/2][jLeft] +
						grid[i0+subSize/2][jRight]
					r := randRange*(rand.Float64() - 0.5)
					grid[i0+subSize/2][j0] = total/4 + r
				}
				if j0 < size-1 {
					iUp := i0 - subSize/2
					if iUp < 0 {
						iUp += size - 1
					}
					iDown := i0 + subSize/2
					if iDown >= size {
						iDown -= size - 1
					}

					total := grid[i0][j0] +
						grid[i0][j0+subSize] +
						grid[iUp][j0+subSize/2] +
						grid[iDown][j0+subSize/2]
					r := randRange*(rand.Float64() - 0.5)
					grid[i0][j0+subSize/2] = total/4 + r
				}
			}
		}

		randRange /= math.Sqrt2
	}

	return grid
}

func resizeGrid(gridOld [][]float64, hNew int, wNew int) [][]float64 {
	hOld := len(gridOld)
	wOld := len(gridOld[0])
	iScale := float64(hOld-1) / float64(hNew-1)
	jScale := float64(wOld-1) / float64(wNew-1)

	gridNew := make([][]float64, hNew)
	for i := range gridNew {
		gridNew[i] = make([]float64, wNew)
	}

	for iNew:=0; iNew<hNew; iNew++ {
		iOld := iScale*float64(iNew)
		iTop := int(math.Floor(iOld))
		iBottom := int(math.Ceil(iOld))
		weightTop := float64(iBottom) - iOld
		weightBottom := 1.0 - weightTop
		for jNew:=0; jNew<wNew; jNew++ {
			jOld := jScale*float64(jNew)
			jLeft := int(math.Floor(jOld))
			jRight := int(math.Ceil(jOld))
			weightLeft := float64(jRight) - jOld
			weightRight := 1.0 - weightLeft

			gridNew[iNew][jNew] = weightTop*weightLeft*gridOld[iTop][jLeft] +
				weightTop*weightRight*gridOld[iTop][jRight] +
				weightBottom*weightLeft*gridOld[iBottom][jLeft] +
				weightBottom*weightRight*gridOld[iBottom][jRight]
		}
	}

	return gridNew
}

func randomSquare(size int) [][]float64 {
	grid := diamondSquare(size)
	return resizeGrid(grid, size, size)
}

func zeroSquare(size int) [][]float64 {
	grid := make([][]float64, size)
	for i:=0; i<size; i++ {
		grid[i] = make([]float64, size)
	}
	return grid
}

func interleaveGrids(a [][]float64, b [][]float64) [][][2]float64 {
	v := make([][][2]float64, len(a))
	for i, _ := range a {
		v[i] = make([][2]float64, len(a[0]))
		for j, _ := range a[i] {
			v[i][j][0] = a[i][j]
			v[i][j][1] = b[i][j]
		}
	}

	return v
}

// for sorting grid positions by value
type cell struct {
	I int
	J int
	X float64
}
type cellSlice []cell
func (cs cellSlice) Less(i, j int) bool {
	return cs[i].X < cs[j].X
}
func (cs cellSlice) Len() int {
	return len(cs)
}
func (cs cellSlice) Swap(i, j int) {
	cs[i], cs[j] = cs[j], cs[i]
}

// "convert" should convert from [0, 1) uniform distribution to target distribution
func applyGridDistribution(grid [][]float64, convert func(float64) float64) [][]float64 {
	h := len(grid)
	w := len(grid[0])

	// get list of positions sorted by value
	cells := make(cellSlice, h*w)
	for i:=0; i<h; i++ {
		for j:=0; j<w; j++ {
			cells[i*w+j] = cell{
				I: i,
				J: j,
				X: grid[i][j],
			}
		}
	}
	sort.Sort(cells)

	// convert grid to "convert" distribution
	delta := 1.0/float64(w*h)
	for k:=0; k<w*h; k++ {
		c := cells[k]
		grid[c.I][c.J] = convert(float64(k)*delta)
	}

	return grid
}

// TODO: use inverse-CDF instead of sampling
func getNormalConverter(mean float64, std float64, precision int) func(float64) float64 {
	samples := make([]float64, precision)
	for i, _ := range samples {
		samples[i] = rand.NormFloat64()*std + mean
	}
	sort.Float64s(samples)

	return func(u float64) float64 {
		pos := u*float64(precision)
		i := int(pos)
		floor := samples[i]
		ceil := floor
		if (i+1 < precision) {
			ceil = samples[i+1]
		}
		weight := 1 - pos + float64(i)

		return floor*weight + ceil*(1.0-weight)
	}
}

func getExpConverter(lambda float64) func(float64) float64 {
	return func(u float64) float64 {
		return -math.Log(1 - u)/lambda
	}
}

func getLogNormalConverter(mu float64, sigma float64, precision int) func(float64) float64 {
	normalConverter := getNormalConverter(mu, sigma, precision)
	return func(u float64) float64 {
		return math.Exp(normalConverter(u))
	}
}

func uniformConverter(u float64) float64 {
	return u
}

// A two-vector with two independently randomly generated coordinates will be biased toward higher
// magnitudes along the diagonal over orthogonal directions.
func removeDiagonalBias(grid [][][2]float64) [][][2]float64 {
	h := len(grid)
	w := len(grid[0])

	for i:=0; i<h; i++ {
		for j:=0; j<w; j++ {
			x := grid[i][j][0]
			y := grid[i][j][1]
			if x == 0 || y == 0 {
				continue
			}
			var ratio float64;
			if x < y {
				ratio = x/y
			} else {
				ratio = y/x
			}
			bias := math.Sqrt(1 + ratio*ratio)
			grid[i][j][0] /= bias
			grid[i][j][1] /= bias
		}
	}

	return grid
}

// Given a grid of 2-vectors, return a grid of their magnitudes.
func getMagnitudes(vectors [][][2]float64) [][]float64 {
	magnitudes := make([][]float64, len(vectors))
	for i, _ := range vectors {
		magnitudes[i] = make([]float64, len(vectors[i]))
		for j, _ := range vectors[i] {
			x := vectors[i][j][0]
			y := vectors[i][j][1]
			magnitudes[i][j] = math.Sqrt(x*x + y*y)
		}
	}

	return magnitudes
}

// Given a grid of 2-vectors and a grid of magnitudes, scale the vectors to have those magnitudes.
func setMagnitudes(vectors [][][2]float64, magnitudes [][]float64) [][][2]float64 {
	for i, _ := range vectors {
		for j, _ := range vectors[i] {
			x := vectors[i][j][0]
			y := vectors[i][j][1]
			r := magnitudes[i][j]
			if x == 0.0 && y == 0.0 {
				continue
			}
			scale_factor := r/math.Sqrt(x*x+y*y)
			vectors[i][j][0] = scale_factor*x
			vectors[i][j][1] = scale_factor*y
		}
	}

	return vectors
}

func randomVectorSquare(size int) float[][][2] {
	a := randomSquare(size)
	b := randomSquare(size)
	v := interleaveGrids(a, b)
	return removeDiagonalBias(v)
}

func zeroVectorSquare(size int) float[][][2] {
	grid := make([][][2]float64, size)
	for i:=0; i<size; i++ {
		grid[i] = make([][2]float64, size)
	}
}

func applyVectorMagnitudeDistribution(grid [][][2]float64, convert func(float64) float64)
		[][][2]float64 {
	magnitudes := applyGridDistribution(getMagnitudes(grid), convert)
	return setMagnitudes(grid, magnitudes)
}
