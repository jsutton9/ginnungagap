package gridvis

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"net/http"
	"os"
)

type Server struct {
	grids map[string]Grid
}

type Grid struct {
	Name string
	Values [][]float64
	Scale []ScalePoint
}

type ScalePoint struct {
	Value float64
	Color [3]uint32
}

func (s Server) AddGrid(name string, values [][]float64, scaleSize int) {
	grid := Grid{
		Name: name,
		Values: values,
		Scale: make([][2]float64, scaleSize),
	}
	grid.setGrayScale(scaleSize)
	s.grids[name] = grid
}

func (g Grid) setGrayScale(scaleSize int) {
	valueList := make([]float64, len(values)*len(values[0]))
	for i,r := range g.Values {
		for j,c := range r {
			valueList[i*len(g.Values)+j] = c
		}
	}
	sort.Float64s(valueList)

	for i,_ := range grid.scale {
		r = i*len(valueList)/(scaleSize-1)
		l = r - 1
		if l < 0 {
			grid.Scale[i] = {
				Value: valueList[0],
				Color: {0, 0, 0}
			}
		} else if r >= len(valueList) {
			grid.Scale[i] = {
				Value: valueList[len(valueList)-1],
				Color: {255, 255, 255}
		} else {
			partition := 0.5*(valueList[l]+valueList[r])
			level := int(255*float64(r)/len(valueList))
			grid.Scale[i] = {
				Value: partition,
				Color: {level, level, level},
			}
		}
	}
}

func (g Grid) image(maxSize int) *image.RGBA {
	cellSize := maxSize/math.max(len(g.Values), len(g.Values[0]))
	height := cellSize*len(g.Values)
	width := cellSize*len(g.Values[0])
	canvas := image.NewRGBA(image.Rect(0, 0, height, width))
	for i, row := range g.Values {
		for j, cell := range row {
			var ci int
			for ci=0; ci<len(g.Scale)-2 && g.Scale[ci+1].Value <= cell; ci++ {}
			blend := (cell-g.Scale[ci].Value)/(g.Scale[ci+1].Value-g.Scale[ci].Value)
			c := image.Color{
				uint32((1-blend)*g.Scale[ci].Color[0] + blend*g.Scale[ci+1].Color[0]),
				uint32((1-blend)*g.Scale[ci].Color[1] + blend*g.Scale[ci+1].Color[1]),
				uint32((1-blend)*g.Scale[ci].Color[2] + blend*g.Scale[ci+1].Color[2]),
			}
			x := cellSize*j
			y := cellSize*(height-i-1)
			r := image.Rect(x, y, x+cellSize, y+cellSize))
			draw.Draw(canvas, r, &image.Uniform{c}, image.ZP, draw.Src)
		}
	}

	return &canvas
}
