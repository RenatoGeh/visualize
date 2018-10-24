package main

import (
	"fmt"
)

func main() {
	fmt.Println("Extracting images...")
	R, L, T, J, Sc, _ := ImagesToData("data/caltech_simple", 25, 8, 200, 200, 7)
	fmt.Println("Learning...")
	S := Gens(R, L, Sc)
	fmt.Println("Performing completion...")
	CompleteData(S, T, J)
	fmt.Println("Done.")
}
