package main

import (
	"fmt"
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
	R, L, T, J, Sc, _ := ImagesToData("data/caltech_simple", 35, 5, 200, 200, 255)
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
	fmt.Println("Performing completion...")
	CompleteData(S, T, J)
	fmt.Println("Done.")
}
