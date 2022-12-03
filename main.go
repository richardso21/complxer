package main

import (
	"fmt"
	"log"
	"os"

	"github.com/richardso21/complxer/asmlc3"
	"github.com/richardso21/complxer/lc3vm"
)

// get initialized LC3 vm
var LC3 = lc3vm.LC3

func main() {
	// run2048()
	testASM()
	fmt.Println("\n==== Program finished ====")
}

func run2048() {
	f, err := os.Open("example/2048.obj")
	if err != nil {
		log.Fatal(err)
	}
	LC3.LoadObjFile(f)
	LC3.Run()
}

func testASM() {
	f, err := os.Open("./example/fibloop.asm")
	if err != nil {
		log.Fatal(err)
	}
	st, err := asmlc3.LoadASMFile(LC3, f)
	fmt.Println(st)
	if err != nil {
		log.Fatal(err)
	}
}
