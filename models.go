package main

import (
	"fmt"
	"github.com/RenatoGeh/gospn/conc"
	"github.com/RenatoGeh/gospn/data"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/learn/dennis"
	"github.com/RenatoGeh/gospn/learn/gens"
	"github.com/RenatoGeh/gospn/spn"
)

func Gens(D spn.Dataset, L []int, Sc map[int]*learn.Variable) []spn.SPN {
	Q := conc.NewSingleQueue(-1)
	n := ClassVar.Categories
	T := data.Split(D, n, L)
	S := make([]spn.SPN, n)
	for i := 0; i < n; i++ {
		Q.Run(func(id int) {
			fmt.Printf("  Learning category %d...\n", id)
			nSc := make(map[int]*learn.Variable)
			for k, v := range Sc {
				if k != ClassVar.Varid {
					nSc[k] = v
				}
			}
			S[id] = gens.Learn(nSc, T[id], 2, 0.01, 4.0, 4.0)
			fmt.Printf("  Finished learning category %d.\n", id)
		}, i)
	}
	Q.Wait()
	return S
}

func Dennis(D spn.Dataset, L []int, Sc map[int]*learn.Variable) []spn.SPN {
	Q := conc.NewSingleQueue(-1)
	n := ClassVar.Categories
	T := data.Split(D, n, L)
	S := make([]spn.SPN, n)
	for i := 0; i < n; i++ {
		Q.Run(func(id int) {
			fmt.Printf("  Learning category %d...\n", id)
			nSc := make(map[int]*learn.Variable)
			for k, v := range Sc {
				if k != ClassVar.Varid {
					nSc[k] = v
				}
			}
			S[id] = dennis.Structure(T[id], nSc, 1, 4, 4, 0.95)
			fmt.Printf("  Finished learning category %d.\n", id)
		}, i)
	}
	Q.Wait()
	return S
}
