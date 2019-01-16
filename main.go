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
	if nargs < 2 || nargs > 4 {
		fmt.Printf("Usage: %s [model filename] | [filename]\n", os.Args[0])
		os.Exit(1)
	}

	var S []spn.SPN
	R, L, T, J, Sc, _ := ImagesToData("data/caltech_simple", -1, 39, 100, 100, 255)
	//R, L, T, J, Sc, _ := ImagesToData("data/english_small", -1, 5, 100, 100, 255)
	fmt.Printf("|R|: %d, |L|: %d, |T|: %d, |J|: %d, |Sc|: %d\n", len(R), len(L), len(T), len(J), len(Sc))
	if nargs == 2 {
		fmt.Println("Loading SPNs...")
		S = Load(os.Args[1])
	} else if nargs == 3 {
		fmt.Println("Extracting images...")
		fmt.Println("Learning...")
		if os.Args[1] == "gens" {
			S = Gens(R, L, Sc)
		} else if os.Args[1] == "dennis" {
			S = Dennis(R, L, Sc)
		} else {
			log.Fatalf("No architecture named %s.", os.Args[1])
		}
		fmt.Println("Saving model...")
		Save(S, os.Args[2])
	}
	if nargs == 4 && os.Args[3] == "cmpl" {
		if S == nil {
			S = Load(os.Args[1])
		}
		fmt.Println("Performing completion...")
		CompleteData(S, T, J)
	} else {
		fmt.Println("Coloring scope...")
		R := data.Split(T, ClassVar.Categories, J)
		Q := conc.NewSingleQueue(-1)
		for i := range S {
			Q.Run(func(id int) {
				ColorScope(S[id], os.Args[1], id, R[id])
			}, i)
		}
		Q.Wait()
	}
	fmt.Println("Done.")
}
