package main

import (
	"fmt"
	"log"
	"os"

	"github.com/richardso21/complxer/lc3vm"
)

// get initialized LC3 vm
var LC3VM = lc3vm.LC3

func main() {
	f, err := os.Open("example/lower.obj")
	if err != nil {
		log.Fatal(err)
	}
	LC3VM.LoadObjFile(f)
	LC3VM.Run()
	fmt.Println("==== Program finished ====")
}
