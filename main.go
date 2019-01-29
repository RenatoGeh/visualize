package main

import (
	"fmt"
	"github.com/RenatoGeh/gospn/conc"
	"github.com/RenatoGeh/gospn/data"
	"github.com/RenatoGeh/gospn/spn"
	"log"
	"os"
)

func main() {
	nargs := len(os.Args)
	if nargs < 4 || nargs > 5 {
		fmt.Printf("Usage: %s [--train | --load] [model filename]\n", os.Args[0])
		os.Exit(1)
	}

	var S []spn.SPN
	//R, L, T, J, Sc, _ := ImagesToData("data/caltech_simple", -1, 0, 100, 100, 7)
	R, L, T, J, Sc, _ := ImagesToData("data/olivetti_simple", -1, 0, 46, 56, 7)
	//R, L, T, J, Sc, _ := ImagesToData("data/english_small", -1, 5, 100, 100, 255)
	fmt.Printf("|R|: %d, |L|: %d, |T|: %d, |J|: %d, |Sc|: %d\n", len(R), len(L), len(T), len(J), len(Sc))
	LearnType = os.Args[2]
	filename := os.Args[3]
	fmt.Printf("Learning type: %s\n", LearnType)
	fmt.Printf("Filename: %s\n", filename)
	if os.Args[1] == "--load" {
		fmt.Println("Loading SPNs...")
		S = Load(filename)
	} else if os.Args[1] == "--train" {
		fmt.Println("Extracting images...")
		fmt.Println("Learning...")
		if LearnType == "gens" {
			S = Gens(R, L, Sc)
		} else if LearnType == "dennis" {
			S = Dennis(R, L, Sc)
		} else {
			log.Fatalf("No architecture named %s.", LearnType)
		}
		fmt.Println("Saving model...")
		Save(S, filename)
	}
	if nargs == 5 && os.Args[4] == "cmpl" {
		if S == nil {
			S = Load(filename)
		}
		fmt.Println("Performing completion...")
		CompleteData(S, T, J)
	} else {
		fmt.Println("Coloring scope...")
		R := data.Split(R, ClassVar.Categories, L)
		Q := conc.NewSingleQueue(-1)
		for i := range S {
			Q.Run(func(id int) {
				ColorScope(S[id], filename, id, R[id])
			}, i)
		}
		Q.Wait()
	}
	fmt.Println("Done.")
}
