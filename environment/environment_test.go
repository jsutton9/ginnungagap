package environment

import (
	"fmt"
	"testing"
)

func TestDiamondSquare(t *testing.T) {
	gridOld := diamondSquare(10)
	fmt.Println("gridOld: ")
	for _, row := range gridOld {
		for _, x := range row {
			fmt.Printf("%5.2f ", x);
		}
		fmt.Println("");
	}
	fmt.Println("")

	gridNew := resizeGrid(gridOld, 22, 7)
	fmt.Println("gridNew: ")
	for _, row := range gridNew {
		for _, x := range row {
			fmt.Printf("%5.2f ", x);
		}
		fmt.Println("");
	}

	t.Fail()
}
