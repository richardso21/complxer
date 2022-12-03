package lc3vm

import (
	"bufio"
	"os"
)

// modified from bufio library to scan 16 bits instead of 8 (makeshift 2-byte scanner)
func scan16Bits(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	return 2, data[0:2], nil
}

func read16Bits(sf *bufio.Scanner) uint16 {
	return uint16(sf.Bytes()[0])<<8 | uint16(sf.Bytes()[1])
}

func (lc3 *LC3vm) LoadObjFile(file *os.File) {
	// read file into memory
	sf := bufio.NewScanner(file)
	sf.Split(scan16Bits)
	sf.Scan()                  // read first line (header)
	currAddr := read16Bits(sf) // set current addr to .ORIG
	lc3.pc = currAddr          // set PC to beginning program
	// now to read the rest of the program
	for ; sf.Scan(); currAddr++ {
		lc3.FillValue(currAddr, read16Bits(sf))
	}
}

func (lc3 *LC3vm) FillValue(addr uint16, value uint16) {
	lc3.mem[addr] = value
}
