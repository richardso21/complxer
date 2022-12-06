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
	// testOBJ()
	testASM()
	LC3.Run()
	fmt.Println(LC3.Reg())
	for _, v := range LC3.Mem()[0x3000:0x300F] {
		fmt.Printf("%04X ", v)
	}
	fmt.Println("\n==== Program finished ====")
}

func testOBJ() {
	f, err := os.Open("./example/fibloop.obj")
	if err != nil {
		log.Fatal(err)
	}
	LC3.LoadObjFile(f)
	// LC3.Run()
}

func testASM() {
	f, err := os.Open("./example/fibloop.asm")
	if err != nil {
		log.Fatal(err)
	}
	st, err := asmlc3.LoadASMFile(LC3, f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(st)
}
