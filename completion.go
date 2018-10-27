package main

import (
	"fmt"
	"github.com/RenatoGeh/gospn/conc"
	"github.com/RenatoGeh/gospn/data"
	"github.com/RenatoGeh/gospn/io"
	"github.com/RenatoGeh/gospn/spn"
	"github.com/RenatoGeh/gospn/sys"
	"math"
)

func subtract(U spn.VarSet, V spn.VarSet) {
	for k, _ := range V {
		delete(U, k)
	}
}

// CompleteHalf performs completion on the top, bottom, left and right halves of an image I, saving
// them to a file. Arguments l and i are label of I and its index in the dataset.
func CompleteHalf(S spn.SPN, I spn.VarSet, l, i int) {
	st := spn.NewStorer()
	tk := st.NewTicket()
	H := make(map[io.CmplType]spn.VarSet)
	SaveInstance(I, fmt.Sprintf("full_%d_%d.png", i, l))
	H[io.Left], H[io.Right] = io.SplitHalf(I, io.Left, Width, Height)
	H[io.Top], H[io.Bottom] = io.SplitHalf(I, io.Top, Width, Height)
	n := Width * Height
	for t, h := range H {
		//_, _, M := spn.StoreMAP(S, h, tk, st)
		for p := 0; p < n; p++ {
			pm, vmax := math.Inf(-1), -1
			for v := 0; v < sys.Max; v++ {
				h[p] = v
				spn.StoreInference(S, h, tk, st)
				pi, _ := st.Single(tk, S)
				if pi > pm {
					pm, vmax = pi, v
				}
				st.Reset(tk)
			}
			h[p] = vmax
		}
		s := fmt.Sprintf("cmpl_%s_%d_%d.png", t, i, l)
		SaveInstance(h, s)
	}
}

func CompleteData(S []spn.SPN, D spn.Dataset, L []int) {
	Q := conc.NewSingleQueue(-1)
	G, H := data.Divide(D, L, Q.Allowed())
	k := len(G)
	for i := range G {
		Q.Run(func(id int) {
			g, h := G[id], H[id]
			for j := range g {
				l := h[j]
				CompleteHalf(S[l], g[j], l, k*id+j)
			}
		}, i)
	}
	Q.Wait()
}
