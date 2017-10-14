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

func copyGrid(grid [][]float64) [][]float64 {
	gridCopy := make([][]float64, len(grid))
	for i, _ := range grid {
		gridCopy[i] = make([]float64, len(grid[i]))
		copy(gridCopy[i], grid[i])
	}

	return gridCopy
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
	gridNormal := applyGridDistribution(copyGrid(gridResized), convertToNormal)
	fmt.Println("gridNormal: ")
	printGrid(gridNormal)

	convertToExp := getExpConverter(0.5)
	gridExp := applyGridDistribution(copyGrid(gridResized), convertToExp)
	fmt.Println("gridExp: ")
	printGrid(gridExp)

	convertToLogNorm := getLogNormalConverter(1.0, 0.5, 1000)
	gridLogNorm := applyGridDistribution(copyGrid(gridResized), convertToLogNorm)
	fmt.Println("gridLogNorm: ")
	printGrid(gridLogNorm)

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

func TestGetMagnitudes(t *testing.T) {
	vectors := [][][2]float64{
		{{-3.0, 4.0}, {-1.5, -2.0}},
		{{1.0, 0.0}, {0.0, 1.0}},
	}
	expected := [][]float64{
		{5.0, 2.5},
		{1.0, 1.0},
	}

	magnitudes := getMagnitudes(vectors)

	for i, _ := range expected {
		for j, _ := range expected[i] {
			if math.Abs((magnitudes[i][j]-expected[i][j])/expected[i][j]) > 0.01 {
				fmt.Println("magnitude off: ")
				fmt.Printf("  vector: <%f, %f>\n",
					vectors[i][j][0], vectors[i][j][1])
				fmt.Printf("  expected magnitude: %f\n", expected[i][j])
				fmt.Printf("  returned magnitude: %f\n", magnitudes[i][j])
				t.Fail()
			}
		}
	}
}

func TestSetMagnitudes(t *testing.T) {
	vectors := [][][2]float64{
		{{-3.0, 4.0}, {-1.5, -2.0}},
		{{1.0, 0.0}, {0.0, 1.0}},
	}
	magnitudes := [][]float64{
		{0.0, -5.0},
		{5.0, -0.5},
	}
	expected := [][][2]float64{
		{{0.0, 0.0}, {3.0, 4.0}},
		{{5.0, 0.0}, {0.0, -0.5}},
	}

	scaled := setMagnitudes(vectors, magnitudes)

	for i, _ := range expected {
		for j, _ := range expected[i] {
			delta_x := scaled[i][j][0] - expected[i][j][0]
			delta_y := scaled[i][j][1] - expected[i][j][1]
			if math.Abs(delta_x) > 0.01 || math.Abs(delta_y) > 0.01 {
				fmt.Println("incorrectly scaled: ")
				fmt.Printf("  expected: <%f, %f>\n",
					expected[i][j][0], expected[i][j][1])
				fmt.Printf("  magnitude: %f\n", magnitudes[i][j])
				fmt.Printf("  scaled: <%f, %f>\n",
					scaled[i][j][0], scaled[i][j][1])
				t.Fail()
			}
		}
	}
}

func TestGridPercentile(t *testing.T) {
	values := make([]float64, 100)
	for i, x := range rand.Perm(100) {
		values[i] = float64(x)
	}
	grid := make([][]float64, 10)
	for i, _ := range grid {
		grid[i] = values[10*i:10*(i+1)]
	}
	split := gridPercentile(grid, 0.355)
	if split < 34 || split > 35 {
		fmt.Printf("expected 34 < split < 35, got %f\n", split)
		t.Fail()
	}
}
