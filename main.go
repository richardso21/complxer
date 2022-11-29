package main

import (
	"fmt"

	"github.com/richardso21/complxer/lc3vm"
)

// get initialized LC3 vm
var LC3VM = lc3vm.LC3

func main() {
	for i := 0; i < 10; i++ {
		LC3VM.RunLine()
	}

	fmt.Printf("%X\n", LC3VM.PC)
}
