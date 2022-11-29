package lc3vm

var opFuncs = [...]func(){ // array of opFuncs in order of op-code value
	LC3._br,
	LC3._add,
	LC3._ld,
	LC3._st,
	LC3._jsr,
	LC3._and,
	LC3._ldr,
	LC3._str,
	LC3._rti,
	LC3._not,
	LC3._ldi,
	LC3._sti,
	LC3._jmp,
	LC3._rsvd,
	LC3._lea,
	LC3._trap,
}

var OP_FUNCMAP = make(map[uint16]func())

func init() {
	for i, fn := range opFuncs {
		OP_FUNCMAP[uint16(i)] = fn
	}
}

const (
	IMM5_TOGGLE = 1 << 5
	JSRR_TOGGLE = 1 << 11
	IMM5        = 0b0000_0000_0001_1111
	DR          = 0b0000_1110_0000_0000
	BaseR       = 0b0000_0001_1100_0000
	SR2         = 0b0000_0000_0000_0111
	PCOFFSET9   = 0b0000_0001_1111_1111
	PCOFFSET11  = 0b0000_0111_1111_1111
	PCOFFSET6   = 0b0000_0000_0011_1111
)

func (lc3 *LC3_st) _br() {
	IR := lc3.IR
	if lc3.NZP&((IR&DR)>>9) != 0 { // check if current NZP matches BR nzp
		lc3.PC += IR & PCOFFSET9
	}
}

func (lc3 *LC3_st) _jmp() {
	IR := lc3.IR
	lc3.PC = lc3.REG[(IR&BaseR)>>6] // set PC to baseR unconditionally
}

func (lc3 *LC3_st) _jsr() {
	IR := lc3.IR
	lc3.REG[7] = lc3.PC      // set R7 to PC
	if IR&JSRR_TOGGLE == 1 { // check if JSR or JSRR
		lc3.PC = lc3.REG[(IR&BaseR)>>6]
	} else {
		lc3.PC += IR & PCOFFSET11
	}
}

func (lc3 *LC3_st) _add() {
	IR := lc3.IR
	add1 := lc3.REG[(IR&BaseR)>>6] // get baseR (SR1)
	var add2 uint16
	if IR&IMM5_TOGGLE == 1 { // get either imm5 or SR2 register
		add2 = IR & IMM5
	} else {
		add2 = lc3.REG[IR&SR2]
	}
	res := add1 + add2
	lc3.REG[(IR&DR)>>9] = res // set DR to add1 + add2
	lc3.updateCC(res)         // update NZP
}

func (lc3 *LC3_st) _and() {
	IR := lc3.IR
	and1 := lc3.REG[(IR&BaseR)>>6] // get baseR (SR1)
	var and2 uint16
	if IR&IMM5_TOGGLE == 1 { // get either imm5 or SR2 register
		and2 = IR & IMM5
	} else {
		and2 = lc3.REG[IR&SR2]
	}
	res := and1 & and2
	lc3.REG[(IR&DR)>>9] = res // set DR to and1 & and2
	lc3.updateCC(res)         // update NZP
}

func (lc3 *LC3_st) _not() {
	IR := lc3.IR
	res := ^lc3.REG[IR&BaseR]
	lc3.REG[(IR&DR)>>9] = res
}

func (lc3 *LC3_st) _ld() {
	IR := lc3.IR
	offset := IR & PCOFFSET9
	location := lc3.PC + offset
	if location == KBDRADDR {
		// toggle keyboard status if reading from keyboard
		lc3.MEMORY[KBSRADDR] = 0
	}
	res := lc3.MEMORY[location]
	lc3.REG[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3_st) _ldr() {
	IR := lc3.IR
	val := lc3.REG[(IR&BaseR)>>6]
	offset := IR & PCOFFSET6
	location := val + offset
	if location == KBDRADDR {
		lc3.MEMORY[KBSRADDR] = 0
	}
	res := lc3.MEMORY[location]
	lc3.REG[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3_st) _ldi() {
	IR := lc3.IR
	location := lc3.MEMORY[lc3.PC+(IR&PCOFFSET9)]
	if location == KBDRADDR {
		lc3.MEMORY[KBSRADDR] = 0
	}
	res := lc3.MEMORY[location]
	lc3.REG[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3_st) _st() {
	IR := lc3.IR
	val := lc3.REG[(IR&DR)>>9]
	location := lc3.PC + (IR & PCOFFSET9)
	if (location < 0x3000) || (location >= 0xFE00) {
		return // prevent overwriting unpriveleged memory
	} else if location == DSRADDR { // toggle display status if writing to display
		lc3.MEMORY[DSRADDR] = 1 << 15
	}
	lc3.MEMORY[location] = val
}

func (lc3 *LC3_st) _str() {
	IR := lc3.IR
	val := lc3.REG[(IR&DR)>>9]
	baseR := lc3.REG[(IR&BaseR)>>6]
	offset := IR & PCOFFSET6
	location := baseR + offset
	if (location < 0x3000) || (location >= 0xFE00) {
		return
	} else if location == DSRADDR {
		lc3.MEMORY[DSRADDR] = 1 << 15
	}
	lc3.MEMORY[location] = val
}

func (lc3 *LC3_st) _sti() {
	IR := lc3.IR
	val := lc3.REG[(IR&DR)>>9]
	itmdAddr := lc3.MEMORY[lc3.PC+(IR&PCOFFSET9)]
	if (itmdAddr < 0x3000) || (itmdAddr >= 0xFE00) {
		return
	} else if itmdAddr == DSRADDR {
		lc3.MEMORY[DSRADDR] = 1 << 15
	}
	lc3.MEMORY[itmdAddr] = val
}

func (lc3 *LC3_st) _lea() {
	IR := lc3.IR
	val := lc3.PC + (IR & PCOFFSET9)
	lc3.REG[(IR&DR)>>9] = val
}

func (lc3 *LC3_st) _rsvd() {
	// do nothing, reserved
}

func (lc3 *LC3_st) updateCC(res uint16) {
	signedRes := int16(res) // convert to signed int
	if signedRes < 0 {
		lc3.NZP = 1 << 2
	} else if signedRes > 0 {
		lc3.NZP = 1
	} else {
		lc3.NZP = 1 << 1
	}
}

// thank god for github copilot
