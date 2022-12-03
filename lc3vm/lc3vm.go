package lc3vm

import (
	"log"
)

// structure representing LC3vm
type LC3vm struct {
	mem       [1 << 16]uint16
	reg       [8]uint16
	pc        uint16
	ir        uint16
	nzp       uint16
	halt      bool
	stackSize uint16
}

// memory-mapped I/O
const (
	KBSRADDR = 0xFE00 // keyboard status
	KBDRADDR = 0xFE02 // keyboard data
	DSRADDR  = 0xFE04 // display status
	DDRADDR  = 0xFE06 // display data
)

func createLC3() *LC3vm {
	var lc3 LC3vm
	lc3.Reset(true, true)
	return &lc3
}

func (lc3 *LC3vm) Stop() {
	lc3.halt = true
}

func (lc3 *LC3vm) Reset(resetPC bool, resetMem bool) {
	lc3.halt = false
	if resetMem {
		lc3.mem = [1 << 16]uint16{}
	}
	// assuming DSR is always on ready :)
	lc3.mem[DSRADDR] = 1 << 15
	lc3.reg = [8]uint16{}
	if resetPC {
		lc3.pc = 0x3000
	}
}

func (lc3 *LC3vm) Run() {
	// run in a loop until HALT or PC goes off-limits
	for !lc3.halt {
		lc3.Step()
	}
}

func (lc3 *LC3vm) Next() {
	currStackSize := lc3.stackSize
	lc3.Step()
	// next over JSR/JSRR statements
	for lc3.stackSize > currStackSize {
		lc3.Step()
	}
}

func (lc3 *LC3vm) Step() {
	if lc3.pc < 0x3000 || lc3.pc >= 0xFE00 { // off-limits from usable memory
		log.Fatalf("PC (%04X) out of bounds!", lc3.pc)
	}
	// FETCH - get instruction from memory
	lc3.fetch()
	// EXECUTE - decode/run instruction
	lc3.execute()
}

func (lc3 *LC3vm) fetch() {
	lc3.ir = lc3.mem[lc3.pc] // set IR
	lc3.pc++                 // increment PC (PC*)
}

func (lc3 *LC3vm) execute() {
	op := lc3.ir >> 12 // get op-code
	opFuncMap[op]()
}

var LC3 = createLC3()
