package environment

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"
)

func printGrid(grid [][]float64) {
	for _, row := range grid {
		for _, x := range row {
			fmt.Printf("%5.2f ", x);
		}
		fmt.Println("");
	}
	fmt.Println("")
}

// This fails intentionally. The output is meant to be visually inspected.
func TestDiamondSquare(t *testing.T) {
	gridOld := diamondSquare(13)
	fmt.Println("gridOld: ")
	printGrid(gridOld)

	gridResized := resizeGrid(gridOld, 22, 11)
	fmt.Println("gridResized: ")
	printGrid(gridResized)

	convertToNormal := getNormalConverter(50.0, 10.0, 1000)
	gridNormal := applyGridDistribution(gridResized, convertToNormal)
	fmt.Println("gridNormal: ")
	printGrid(gridNormal)

	convertToExp := getExpConverter(0.5)
	gridExp := applyGridDistribution(gridResized, convertToExp)
	fmt.Println("gridExp: ")
	printGrid(gridExp)

	t.Fail()
}

func TestRemoveDiagonalBias(t *testing.T) {
	size := 100
	rand.Seed(time.Now().UnixNano())
	grid := make([][][2]float64, size)
	directions := make([][]float64, size)
	for i := range grid {
		grid[i] = make([][2]float64, size)
		directions[i] = make([]float64, size)
		for j := range grid[i] {
			grid[i][j][0] = 2*rand.Float64() - 1
			grid[i][j][1] = 2*rand.Float64() - 1
			directions[i][j] = math.Atan2(grid[i][j][1], grid[i][j][0])
		}
	}

	grid = removeDiagonalBias(grid)

	magnitudeMean := 0.0
	for i := range grid {
		for j := range grid {
			x := grid[i][j]
			magnitudeMean += math.Sqrt(x[0]*x[0] + x[1]*x[1])
		}
	}
	magnitudeMean /= float64(size*size)
	if math.Abs(magnitudeMean) > 0.51 {
		fmt.Printf("mean magnitude after debiasing: %f\n", magnitudeMean)
		t.Fail()
	}
	for i := range grid {
		for j := range grid[i] {
			direction := math.Atan2(grid[i][j][1], grid[i][j][0])
			delta := direction - directions[i][j]
			if delta > 0.01*math.Pi && delta < 1.98*math.Pi {
				fmt.Printf("direction changed; delta = %f\n", delta)
				fmt.Printf("  old direction: %f\n", directions[i][j])
				fmt.Printf("  new direction: %f\n", direction)
				fmt.Printf("  new vector: <%f, %f>\n", grid[i][j][0], grid[i][j][1])
				t.Fail()
			}
		}
	}
}
