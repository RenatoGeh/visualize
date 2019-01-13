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
	if nargs < 2 || nargs > 3 {
		fmt.Printf("Usage: %s [model filename] | [filename]\n", os.Args[0])
		os.Exit(1)
	}

	var S []spn.SPN
	//R, L, T, J, Sc, _ := ImagesToData("data/caltech_pruned", -1, 5, 100, 100, 255)
	R, L, T, J, Sc, _ := ImagesToData("data/english_small", -1, 5, 100, 100, 255)
	if nargs == 2 {
		fmt.Println("Loading SPNs...")
		var err error
		S = Load(os.Args[1])
		if err != nil {
			panic(err)
		}
	} else {
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
