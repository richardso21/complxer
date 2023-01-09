package main

import (
	"fmt"
	"log"
	"os"

	"github.com/richardso21/complxer/asmlc3"
	"github.com/richardso21/complxer/assembler"
	"github.com/richardso21/complxer/lc3vm"
)

// get initialized LC3 vm
var LC3 = lc3vm.LC3

func main() {
	// testOBJ()
	// testASM()
	testASMnew()
	LC3.Run()

	fmt.Println("\n==== Memory Slice ====")
	memSlice := LC3.Mem()[0x3000:LC3.Pc()]
	for i := range memSlice {
		fmt.Printf("%04X ", memSlice[i])
	}

	fmt.Println("\n==== Registers ====")
	for i, val := range LC3.Reg() {
		fmt.Printf("R%d: 0x%04X ", i, val)
	}

	fmt.Println("\n==== Program finished ====")
}

func testOBJ() {
	f, err := os.Open("./example/fibloop.obj")
	if err != nil {
		log.Fatal(err)
	}
	LC3.LoadObjFile(f)
}

func testASM() {
	f, err := os.Open("./example/fibloop.asm")
	if err != nil {
		log.Fatal(err)
	}
	_, err = asmlc3.LoadASMFile(LC3, f)
	if err != nil {
		log.Fatal(err)
	}
}

func testASMnew() {
	f, err := os.Open("./example/lower.asm")
	if err != nil {
		log.Fatal(err)
	}
	objFN, err := assembler.AsmToObj(f)
	if err != nil {
		log.Fatal(err)
	}
	objFile, err := os.Open(objFN)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("assembled!")
	LC3.LoadObjFile(objFile)
	fmt.Println("loaded!")
}
