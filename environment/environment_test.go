package environment

import (
	"fmt"
	"testing"
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
