package gridvis

import (
	"bytes"
	"fmt"
	"html/template"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
)

const MAX_SIZE = 400

type Server struct {
	Grids map[string]Grid
	indexTemplate *template.Template
	gridTemplate *template.Template
}

type Grid struct {
	Name string
	Values [][]float64
	Scale []ScalePoint
}

type ScalePoint struct {
	Value float64
	Color [3]uint32
	ColorString string
}

func NewServer(templateDir string) Server {
	indexTemplate, err := template.ParseFiles(templateDir+"index.html")
	if err != nil {
		panic(err)
	}

	gridTemplate, err := template.ParseFiles(templateDir+"grid.html")
	if err != nil {
		panic(err)
	}

	f, err := os.Open(templateDir+"style.css")
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header()["Content-Type"] = []string{"text/css"}
		w.Write(b)
	})

	return Server{
		Grids: make(map[string]Grid),
		indexTemplate: indexTemplate,
		gridTemplate: gridTemplate,
	}
}

func (s Server) Serve() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (s Server) AddGrid(name string, values [][]float64, scaleSize int) {
	grid := Grid{
		Name: name,
		Values: values,
		Scale: make([]ScalePoint, scaleSize),
	}
	grid.setGrayScale(scaleSize)
	s.Grids[name] = grid

	var b bytes.Buffer
	err := s.gridTemplate.Execute(&b, grid)
	if err != nil {
		panic(err)
	}
	page := make([]byte, b.Len())
	_, err = b.Read(page)
	if err != nil {
		panic(err)
	}
	image := grid.imageBytes(MAX_SIZE)
	http.HandleFunc("/"+name+".html", func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	})
	http.HandleFunc("/"+name+".png", func(w http.ResponseWriter, r *http.Request) {
		w.Write(image)
	})

	s.renderIndex()
}

func (s Server) renderIndex() {
	var b bytes.Buffer
	err := s.indexTemplate.Execute(&b, s)
	if err != nil {
		panic(err)
	}
	page := make([]byte, b.Len())
	_, err = b.Read(page)
	if err != nil {
		panic(err)
	}
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Write(page)
	}
	http.HandleFunc("/index.html", handler)
}

func (g Grid) setGrayScale(scaleSize int) {
	valueList := make([]float64, len(g.Values)*len(g.Values[0]))
	for i,r := range g.Values {
		for j,c := range r {
			valueList[i*len(g.Values)+j] = c
		}
	}
	sort.Float64s(valueList)

	for i,_ := range g.Scale {
		r := i*len(valueList)/(scaleSize-1)
		l := r - 1
		if l < 0 {
			g.Scale[i] = newScalePoint(valueList[0], [3]uint32{0, 0, 0})
		} else if r >= len(valueList) {
			g.Scale[i] = newScalePoint(valueList[len(valueList)-1], [3]uint32{255, 255, 255})
		} else {
			partition := 0.5*(valueList[l]+valueList[r])
			level := uint32(float64(255)*float64(r)/float64(len(valueList)))
			g.Scale[i] = newScalePoint(partition, [3]uint32{level, level, level})
		}
	}
}

func newScalePoint(value float64, color [3]uint32) ScalePoint {
	return ScalePoint{
		Value: value,
		Color: color,
		ColorString: fmt.Sprintf("%06X", 65536*color[0]+256*color[1]+color[2]),
	}
}

func (g Grid) imageBytes(maxSize int) []byte {
	image := g.image(maxSize)
	var b bytes.Buffer
	err := png.Encode(&b, image)
	if err != nil {
		panic(err)
	}
	out := make([]byte, b.Len())
	_, err = b.Read(out)
	if err != nil {
		panic(err)
	}
	return out
}

func (g Grid) image(maxSize int) *image.RGBA {
	var cellSize int
	if len(g.Values[0]) > len(g.Values) {
		cellSize = maxSize/len(g.Values[0])
	} else {
		cellSize = maxSize/len(g.Values)
	}
	height := cellSize*len(g.Values)
	width := cellSize*len(g.Values[0])
	canvas := image.NewRGBA(image.Rect(0, 0, height, width))
	for i, row := range g.Values {
		for j, cell := range row {
			var ci int
			for ci=0; ci<len(g.Scale)-2 && g.Scale[ci+1].Value <= cell; ci++ {}
			blend := (cell-g.Scale[ci].Value)/(g.Scale[ci+1].Value-g.Scale[ci].Value)
			c := color.RGBA{
				uint8((1-blend)*float64(g.Scale[ci].Color[0]) + blend*float64(g.Scale[ci+1].Color[0])),
				uint8((1-blend)*float64(g.Scale[ci].Color[1]) + blend*float64(g.Scale[ci+1].Color[1])),
				uint8((1-blend)*float64(g.Scale[ci].Color[2]) + blend*float64(g.Scale[ci+1].Color[2])),
				255,
			}
			x := cellSize*j
			y := height - cellSize*(i+1)
			r := image.Rect(x, y, x+cellSize, y+cellSize)
			draw.Draw(canvas, r, &image.Uniform{c}, image.ZP, draw.Src)
		}
	}

	return canvas
}
