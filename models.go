package main

import (
	"github.com/RenatoGeh/gospn/data"
	"github.com/RenatoGeh/gospn/learn"
	"github.com/RenatoGeh/gospn/learn/gens"
	"github.com/RenatoGeh/gospn/spn"
)

func Gens(D spn.Dataset, L []int, Sc map[int]*learn.Variable) spn.SPN {
	T := data.MergeLabel(D, L, ClassVar)
	return gens.LearnConcurrent(Sc, T, 2, 0.01, 4.0, 4.0, -1)
}
