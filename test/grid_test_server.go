package main

import (
	"github.com/jsutton9/ginnungagap/environment"
	"github.com/jsutton9/ginnungagap/environment/gridvis"
)

func main() {
	s := gridvis.NewServer(
		"/home/jake/code/go/src/github.com/jsutton9/ginnungagap/environment/gridvis/")
	gridOld := environment.DiamondSquare(13)
	s.AddGrid("diamond-square", gridOld, 20)
	s.Serve()
}
