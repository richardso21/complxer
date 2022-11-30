package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/richardso21/complxer/lc3vm"
)

// get initialized LC3 vm
var LC3VM = lc3vm.LC3

func main() {
	f, err := os.Open("example/fibloop.txt")
	if err != nil {
		log.Fatal(err)
	}
	sf := bufio.NewScanner(f)
	for i := 0; sf.Scan(); i++ {
		binLine, err := strconv.ParseUint(sf.Text(), 2, 16)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%04X\n", binLine)
		LC3VM.MEMORY[0x3000+uint16(i)] = uint16(binLine)
	}
	LC3VM.Run()
	fmt.Printf("%04X %04X\n", LC3VM.PC, LC3VM.REG[0])
}
