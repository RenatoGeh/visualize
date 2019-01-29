package main

import (
	"container/heap"
	"fmt"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/lucasb-eyer/go-colorful"
	"image"
	"image/png"
	"os"
)

const ScopeThreshold = 200

var LearnType string

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

type heapElement struct {
	p float64
	M spn.VarSet
	i int
}
type maxHeap []*heapElement

func (h maxHeap) Len() int           { return len(h) }
func (h maxHeap) Less(i, j int) bool { return h[i].p > h[j].p }
func (h maxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].i = i
	h[j].i = j
}
func (h *maxHeap) Push(e interface{}) {
	i := e.(*heapElement)
	i.i = len(*h)
	*h = append(*h, i)
}
func (h *maxHeap) Pop() interface{} {
	o := *h
	n := len(o)
	e := o[n-1]
	e.i = -1
	*h = o[0 : n-1]
	return e
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func completify(path string, S spn.SPN, idx int) {
	fmt.Printf("%d: Completify started\n", idx)
	ch := S.Ch()
	var H maxHeap
	for i, c := range ch {
		p, M := spn.MAP(c, make(spn.VarSet))
		heap.Push(&H, &heapElement{
			p: p,
			M: M,
			i: i,
		})
	}
	var n int
	if LearnType == "gens" {
		n = min(H.Len(), 3)
	} else {
		n = 1
	}
	for i := 0; i < n && H.Len() > 0; i++ {
		fmt.Printf("%d: Completing child %d\n", idx, i)
		e := heap.Pop(&H).(*heapElement)
		I := image.NewRGBA(image.Rect(0, 0, Width, Height))
		for p, v := range e.M {
			x, y := p%Width, p/Width
			// Order: Blue, Green, Red
			u := float64(v) / float64(Max)
			//if LearnType == "gens" {
			//I.Set(x, y, colorful.Color{
			//float64(((i & 2) >> 1)) * u,
			//float64((i & 1)) * u,
			//float64((((i + 1) >> 1) ^ 1)) * u,
			//})
			//} else {
			I.Set(x, y, colorful.Color{u, u, u})
			//}
		}
		var filename string
		if LearnType == "gens" {
			filename = fmt.Sprintf("%s%d_%d.png", path, idx, i)
		} else {
			filename = fmt.Sprintf("%s%d.png", path, idx)
		}
		f, err := os.Create(filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%d: Drawing child %d...\n", idx, i)
		if err = png.Encode(f, I); err != nil {
			fmt.Println(err)
		}
		f.Close()
	}
	fmt.Printf("%d: Completify ended\n", idx)
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
	spn.BreadthFirst(S, func(Z spn.SPN) int {
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
		var m int
		var Sc [][]int
		for _, c := range ch {
			csc := c.Sc()
			Sc = append(Sc, csc)
			m += len(csc)
		}
		if m < ScopeThreshold {
			return -1
		}
		if t := Z.Type(); t == "sum" {
			if LearnType == "gens" {
				completify(spath, Z, sums)
				sums++
			}
		} else if t == "product" {
			colorify(ppath, Sc, C, MeanImage(T), prods)
			if LearnType == "dennis" {
				for i, c := range ch {
					if c.Type() == "sum" && len(Sc[i]) > ScopeThreshold {
						completify(fmt.Sprintf("%s%d_", spath, prods), c, sums)
						sums++
					}
				}
			}
			prods++
		}
		return 0
	})
	fmt.Println(prods)
}
