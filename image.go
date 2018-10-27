package main

import (
	"fmt"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/sys"
	"gocv.io/x/gocv"
	"image"
	"io/ioutil"
	"os"
)

var (
	Width  int
	Height int
	Max    int

	ClassVar *learn.Variable
)

// ImageToInstance takes an image file, resizes to (w, h), quantizes if max != 255 and returns a
// VarSet representation of the image.
func ImageToInstance(file string, w, h, max int) map[int]int {
	img := gocv.IMRead(file, gocv.IMReadGrayScale)
	gocv.Resize(img, &img, image.Point{w, h}, 0, 0, gocv.InterpolationLinear)
	if max != 255 {
		img.MultiplyFloat(float32(max) / 255.0)
	}
	ptr := img.DataPtrUint8()
	I := make(map[int]int)
	for i, p := range ptr {
		I[i] = int(p)
	}
	return I
}

// ImagesToData takes a directory dir and a number of samples n and returns training and test
// datasets containing n and m elements of each class in dir respectively. Arguments w, h and max
// are width, height and max pixel value (resolution).
func ImagesToData(dir string, n, m, w, h, max int) ([]map[int]int, []int, []map[int]int, []int, map[int]*learn.Variable, int) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	var cats []os.FileInfo
	for _, d := range files {
		if d.IsDir() {
			cats = append(cats, d)
		}
	}
	k := len(cats)
	var D, E []map[int]int
	var L, J []int
	for i, c := range cats {
		cpath := dir + "/" + c.Name()
		files, err = ioutil.ReadDir(cpath)
		if err != nil {
			panic(err)
		}
		for j, f := range files {
			fpath := cpath + "/" + f.Name()
			if j < n {
				I := ImageToInstance(fpath, w, h, max)
				D = append(D, I)
				L = append(L, i)
			} else if j < n+m {
				I := ImageToInstance(fpath, w, h, max)
				E = append(E, I)
				J = append(J, i)
			} else {
				break
			}
		}
	}
	N := w * h
	Sc := make(map[int]*learn.Variable)
	for i := 0; i < N; i++ {
		Sc[i] = &learn.Variable{Varid: i, Categories: max + 1, Name: fmt.Sprintf("Pixel %d", i)}
	}
	Sc[N] = &learn.Variable{Varid: N, Categories: k, Name: fmt.Sprintf("Class")}
	sys.Width, sys.Height, sys.Max = w, h, max+1
	Width, Height, Max = w, h, max
	ClassVar = Sc[N]
	return D, L, E, J, Sc, k
}

func colorize(M gocv.Mat, ch int) gocv.Mat {
	C := make([]gocv.Mat, 3)
	C[ch] = M
	I := gocv.NewMat()
	gocv.Merge(C, &I)
	return I
}

// Saves an spn.VarSet to filename.
func SaveInstance(I map[int]int, filename string) {
	c := 255.0 / float32(Max)
	M := gocv.NewMatWithSize(Width, Height, gocv.MatTypeCV8U)
	for k, v := range I {
		x := k % Width
		y := k / Width
		p := float32(v) * c
		M.SetUCharAt(y, x, uint8(p))
	}
	gocv.IMWrite(filename, M)
}
