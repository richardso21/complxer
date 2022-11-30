package lc3vm

import (
	"fmt"
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
// TODO: implement ST/LD with memory-mapped I/O
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
	lc3.REG = [8]uint16{}
	if resetPC {
		lc3.PC = 0x3000
	}
}

func (lc3 *LC3_st) Run() {
	// run in a loop until HALT or PC goes off-limits
	for !lc3.HALT {
		if lc3.PC < 0x3000 || lc3.PC >= 0xFE00 { // off-limits from usable memory
			lc3.HALT = true
			continue
		}
		// syncIO - update memory with display/keyboard, and vice versa
		lc3.syncIO()
		// FETCH - get instruction from memory
		lc3.fetch()
		// EXECUTE - decode/run instruction
		lc3.execute()
	}
}

func (lc3 *LC3_st) syncIO() {
	// check if keyboard is ready (STDIN has data)
	stat, _ := os.Stdin.Stat()
	if ((stat.Mode() & os.ModeCharDevice) == 0) && lc3.MEMORY[KBSRADDR] == 0 {
		// set keyboard status
		lc3.MEMORY[KBSRADDR] = 1 << 15
		// read char into keyboard data
		fmt.Scanf("%c", &lc3.MEMORY[KBDRADDR])
	}

	// check if display is ready (we put something into DDRADDR)
	if lc3.MEMORY[DSRADDR] != 0 {
		// print char
		fmt.Printf("%c", lc3.MEMORY[DDRADDR])
		lc3.MEMORY[DSRADDR] = 0 // clear display status
		lc3.MEMORY[DDRADDR] = 0 // clear display data
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
