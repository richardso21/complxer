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
		lc3.PC += signExt(IR&PCOFFSET9, 8)
	}
}

func (lc3 *LC3_st) _jmp() {
	IR := lc3.IR
	lc3.PC = lc3.REG[(IR&BaseR)>>6] // set PC to baseR unconditionally
}

func (lc3 *LC3_st) _jsr() {
	IR := lc3.IR
	lc3.REG[7] = lc3.PC                // set R7 to PC
	if IR&JSRR_TOGGLE == JSRR_TOGGLE { // check if JSR or JSRR
		lc3.PC = lc3.REG[(IR&BaseR)>>6]
	} else {
		lc3.PC += signExt(IR&PCOFFSET11, 10)
	}
}

func (lc3 *LC3_st) _add() {
	IR := lc3.IR
	add1 := lc3.REG[(IR&BaseR)>>6] // get baseR (SR1)
	var add2 uint16
	if IR&IMM5_TOGGLE == IMM5_TOGGLE { // get either imm5 or SR2 register
		add2 = signExt(IR&IMM5, 4)
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
	if IR&IMM5_TOGGLE == IMM5_TOGGLE { // get either imm5 or SR2 register
		and2 = signExt(IR&IMM5, 4)
	} else {
		and2 = lc3.REG[IR&SR2]
	}
	res := and1 & and2
	lc3.REG[(IR&DR)>>9] = res // set DR to and1 & and2
	lc3.updateCC(res)         // update NZP
}

func (lc3 *LC3_st) _not() {
	IR := lc3.IR
	res := ^lc3.REG[(IR&BaseR)>>6]
	lc3.REG[(IR&DR)>>9] = res
}

func (lc3 *LC3_st) _ld() {
	IR := lc3.IR
	offset := signExt(IR&PCOFFSET9, 8)
	location := lc3.PC + offset
	res := lc3.readMemory(location)
	lc3.REG[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3_st) _ldr() {
	IR := lc3.IR
	val := lc3.REG[(IR&BaseR)>>6]
	offset := signExt(IR&PCOFFSET6, 5)
	location := val + offset
	res := lc3.readMemory(location)
	lc3.REG[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3_st) _ldi() {
	IR := lc3.IR
	location := lc3.readMemory(lc3.PC + signExt(IR&PCOFFSET9, 8))
	res := lc3.readMemory(location)
	lc3.REG[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3_st) _st() {
	IR := lc3.IR
	val := lc3.REG[(IR&DR)>>9]
	location := lc3.PC + signExt(IR&PCOFFSET9, 8)
	lc3.writeMemory(location, val)
}

func (lc3 *LC3_st) _str() {
	IR := lc3.IR
	val := lc3.REG[(IR&DR)>>9]
	baseR := lc3.REG[(IR&BaseR)>>6]
	offset := signExt(IR&PCOFFSET6, 5)
	location := baseR + offset
	lc3.writeMemory(location, val)
}

func (lc3 *LC3_st) _sti() {
	IR := lc3.IR
	val := lc3.REG[(IR&DR)>>9]
	location := lc3.readMemory(lc3.PC + signExt(IR&PCOFFSET9, 8))
	lc3.writeMemory(location, val)
}

func (lc3 *LC3_st) _lea() {
	IR := lc3.IR
	val := lc3.PC + signExt(IR&PCOFFSET9, 8)
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

func signExt(val uint16, signedBit int) uint16 {
	var sigBit uint16 = 1 << signedBit
	sigMask := ^(sigBit - 1)
	if val&sigBit == sigBit {
		return val | sigMask
	}
	return val
}
