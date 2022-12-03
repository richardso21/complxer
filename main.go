package main

import (
	"fmt"
	"log"
	"os"

	"github.com/richardso21/complxer/lc3vm"
)

// get initialized LC3 vm
var LC3 = lc3vm.LC3

func main() {
	f, err := os.Open("example/2048.obj")
	if err != nil {
		log.Fatal(err)
	}
	LC3.LoadObjFile(f)
	LC3.Run()
	fmt.Println("\n==== Program finished ====")
}
