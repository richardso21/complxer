package lc3vm

import (
	"bufio"
	"os"
)

// structure representing LC3_st
type LC3_st struct {
	MEMORY [1 << 16]uint16
	REG    [8]uint16
	PC     uint16
	IR     uint16
	NZP    uint16
	HALT   bool
}

// memory-mapped I/O
const (
	KBSRADDR = 0xFE00 // keyboard status
	KBDRADDR = 0xFE02 // keyboard data
	DSRADDR  = 0xFE04 // display status
	DDRADDR  = 0xFE06 // display data
)

func createLC3() *LC3_st {
	var lc3 LC3_st
	lc3.Reset(true)
	return &lc3
}

func (lc3 *LC3_st) Reset(resetPC bool) {
	lc3.HALT = false
	lc3.MEMORY = [1 << 16]uint16{}
	// assuming DSR is always on ready :)
	lc3.MEMORY[DSRADDR] = 1 << 15
	lc3.REG = [8]uint16{}
	if resetPC {
		lc3.PC = 0x3000
	}
}

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

func (lc3 *LC3_st) LoadObjFile(filename *os.File) {
	// read file into memory
	sf := bufio.NewScanner(filename)
	sf.Split(scan16Bits)
	sf.Scan()                  // read first line (header)
	currAddr := read16Bits(sf) // set current addr to .ORIG
	lc3.PC = currAddr          // set PC to beginning program
	// now to read the rest of the program
	for ; sf.Scan(); currAddr++ {
		lc3.MEMORY[currAddr] = read16Bits(sf)
	}
}

func (lc3 *LC3_st) Run() {
	// run in a loop until HALT or PC goes off-limits
	for !lc3.HALT {
		if lc3.PC < 0x3000 || lc3.PC >= 0xFE00 { // off-limits from usable memory
			lc3.HALT = true
			continue
		}
		// FETCH - get instruction from memory
		lc3.fetch()
		// EXECUTE - decode/run instruction
		lc3.execute()
	}
}

func (lc3 *LC3_st) fetch() {
	lc3.IR = lc3.MEMORY[lc3.PC] // set IR
	lc3.PC++                    // increment PC (PC*)
}

func (lc3 *LC3_st) execute() {
	op := lc3.IR >> 12 // get op-code
	OP_FUNCMAP[op]()
}

var LC3 = createLC3()
