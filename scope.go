package main

import (
	"fmt"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/png"
	"os"
)

func colorify(path string, Sc [][]int, C []colorful.Color, V spn.VarSet, idx int) {
	I := image.NewRGBA(image.Rect(0, 0, Width, Height))
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			p := x + y*Width
			q := float64(V[p]) / float64(Max)
			I.Set(x, y, colorful.Color{q, q, q})
		}
	}
	for i, sc := range Sc {
		for _, p := range sc {
			x, y := p%Width, p/Width
			I.Set(x, y, C[i])
		}
	}
	f, err := os.Create(fmt.Sprintf("%s%d.png", path, idx))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	if err = png.Encode(f, I); err != nil {
		fmt.Println(err)
	}
}

func ColorScope(S spn.SPN, mdl string, is int, T []map[int]int) {
	path := fmt.Sprintf("results/%s/%d/", mdl, is)
	spath, ppath := path+"sums/", path+"products/"
	if err := os.MkdirAll(spath, 0700); err != nil {
		panic(err)
	} else if err = os.MkdirAll(ppath, 0700); err != nil {
		panic(err)
	}

	fmt.Println("Computing scope...")
	spn.ComputeScope(S)
	var sums, prods int
	spn.BreadthFirst(S, func(Z spn.SPN) bool {
		ch := Z.Ch()
		n := len(ch)
		C, err := colorful.SoftPaletteEx(n, colorful.SoftPaletteSettings{
			CheckColor:  nil,
			Iterations:  50,
			ManySamples: true,
		})
		if err != nil {
			fmt.Println(err)
			fmt.Printf("Length: %d\n", len(C))
		}
		var Sc [][]int
		for _, c := range ch {
			Sc = append(Sc, c.Sc())
		}
		if t := Z.Type(); t == "sum" {
			//colorify(spath, Sc, C, T[0], sums)
			sums++
		} else if t == "product" {
			colorify(ppath, Sc, C, T[0], prods)
			prods++
		}
		return true
	})
	fmt.Println(prods)
}
