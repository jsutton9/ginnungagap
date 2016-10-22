package environment

import (
	"math"
	"math/rand"
	"time"
)

func diamondSquare(minSize int) *[][]float64 {
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

	return &grid
}

func resizeGrid(gridOld *[][]float64, hNew int, wNew int) *[][]float64 {
	hOld := len(*gridOld)
	wOld := len((*gridOld)[0])
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

			gridNew[iNew][jNew] = weightTop*weightLeft*(*gridOld)[iTop][jLeft] +
				weightTop*weightRight*(*gridOld)[iTop][jRight] +
				weightBottom*weightLeft*(*gridOld)[iBottom][jLeft] +
				weightBottom*weightRight*(*gridOld)[iBottom][jRight]
		}
	}

	return &gridNew
}
