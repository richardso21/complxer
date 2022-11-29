package lc3vm

// structure representing LC3_t
type LC3_t struct {
	MEMORY [1 << 16]uint16
	REG    [8]uint16
	PC     uint16
	IR     uint16
	NZP    uint16
	_HALT  bool
}

// memory-mapped I/O
// TODO: implement ST/LD with memory-mapped I/O
const (
	KBSRADDR = 0xFE00 // keyboard status
	KBDRADDR = 0xFE02 // keyboard data
	DSRADDR  = 0xFE04 // display status
	DDRADDR  = 0xFE06 // display data
)

func createLC3() LC3_t {
	var lc3 LC3_t
	lc3.Reset(true)
	return lc3
}

func (lc3 *LC3_t) Reset(resetPC bool) {
	lc3._HALT = false
	lc3.MEMORY = [1 << 16]uint16{}
	lc3.REG = [8]uint16{}
	if resetPC {
		lc3.PC = 0x3000
	}
}

func (lc3 *LC3_t) RunLine() {
	// if halted, do nothing
	if lc3._HALT {
		return
	} else if lc3.PC > 0xFE00 || lc3.PC < 0x3000 { // off-limits from usable memory
		lc3._HALT = true
		return
	}
	// FETCH - get instruction from memory
	lc3.fetch()
	// EXECUTE - decode/run instruction
	lc3.execute(lc3.IR)
}

func (lc3 *LC3_t) fetch() {
	lc3.IR = lc3.MEMORY[lc3.PC] // set IR
	lc3.PC++                    // increment PC (PC*)
}

func (lc3 *LC3_t) execute(IR uint16) {
	op := IR >> 12 // get op-code
	OP_FUNCMAP[op]()
}

var LC3VM = createLC3()
