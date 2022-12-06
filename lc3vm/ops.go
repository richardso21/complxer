package lc3vm

var opFuncs = [...]func(){ // array of opFuncs in order of op-code value
	LC3.br,
	LC3.add,
	LC3.ld,
	LC3.st,
	LC3.jsr,
	LC3.and,
	LC3.ldr,
	LC3.str,
	func() {}, // RTI not implemented
	LC3.not,
	LC3.ldi,
	LC3.sti,
	LC3.jmp,
	func() {}, //reserved
	LC3.lea,
	LC3.trap,
}

var opFuncMap = make(map[uint16]func())

func init() {
	for i, fn := range opFuncs {
		opFuncMap[uint16(i)] = fn
	}
}

const (
	IMM5_TOGGLE = 1 << 5
	JSR_TOGGLE  = 1 << 11
	IMM5        = 0b0000_0000_0001_1111
	DR          = 0b0000_1110_0000_0000
	BaseR       = 0b0000_0001_1100_0000
	SR2         = 0b0000_0000_0000_0111
	PCOFFSET9   = 0b0000_0001_1111_1111
	PCOFFSET11  = 0b0000_0111_1111_1111
	PCOFFSET6   = 0b0000_0000_0011_1111
)

func (lc3 *LC3vm) br() {
	IR := lc3.ir
	if lc3.nzp&((IR&DR)>>9) != 0 { // check if current NZP matches BR nzp
		lc3.pc += signExt(IR&PCOFFSET9, 9)
	}
}

func (lc3 *LC3vm) jmp() {
	IR := lc3.ir
	reg := IR & BaseR >> 6
	if reg == 7 {
		lc3.decrStack() // if R7 (RET), decrement stack
	}
	lc3.pc = lc3.reg[reg] // set PC to baseR unconditionally
}

func (lc3 *LC3vm) jsr() {
	IR := lc3.ir
	lc3.reg[7] = lc3.pc              // set R7 to PC
	if IR&JSR_TOGGLE == JSR_TOGGLE { // check if JSR or JSRR
		lc3.pc += signExt(IR&PCOFFSET11, 11)
	} else {
		lc3.pc = lc3.reg[(IR&BaseR)>>6]
	}
	lc3.incrStack()
}

func (lc3 *LC3vm) add() {
	IR := lc3.ir
	add1 := lc3.reg[(IR&BaseR)>>6] // get baseR (SR1)
	var add2 uint16
	if IR&IMM5_TOGGLE == IMM5_TOGGLE { // get either imm5 or SR2 register
		add2 = signExt(IR&IMM5, 5)
	} else {
		add2 = lc3.reg[IR&SR2]
	}
	res := add1 + add2
	lc3.reg[(IR&DR)>>9] = res // set DR to add1 + add2
	lc3.updateCC(res)         // update NZP
}

func (lc3 *LC3vm) and() {
	IR := lc3.ir
	and1 := lc3.reg[(IR&BaseR)>>6] // get baseR (SR1)
	var and2 uint16
	if IR&IMM5_TOGGLE == IMM5_TOGGLE { // get either imm5 or SR2 register
		and2 = signExt(IR&IMM5, 5)
	} else {
		and2 = lc3.reg[IR&SR2]
	}
	res := and1 & and2
	lc3.reg[(IR&DR)>>9] = res // set DR to and1 & and2
	lc3.updateCC(res)         // update NZP
}

func (lc3 *LC3vm) not() {
	IR := lc3.ir
	res := ^lc3.reg[(IR&BaseR)>>6]
	lc3.reg[(IR&DR)>>9] = res
}

func (lc3 *LC3vm) ld() {
	IR := lc3.ir
	offset := signExt(IR&PCOFFSET9, 9)
	location := lc3.pc + offset
	res := lc3.readMemory(location)
	lc3.reg[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3vm) ldr() {
	IR := lc3.ir
	val := lc3.reg[(IR&BaseR)>>6]
	offset := signExt(IR&PCOFFSET6, 6)
	location := val + offset
	res := lc3.readMemory(location)
	lc3.reg[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3vm) ldi() {
	IR := lc3.ir
	location := lc3.readMemory(lc3.pc + signExt(IR&PCOFFSET9, 9))
	res := lc3.readMemory(location)
	lc3.reg[(IR&DR)>>9] = res
	lc3.updateCC(res)
}

func (lc3 *LC3vm) st() {
	IR := lc3.ir
	val := lc3.reg[(IR&DR)>>9]
	location := lc3.pc + signExt(IR&PCOFFSET9, 9)
	lc3.writeMemory(location, val)
}

func (lc3 *LC3vm) str() {
	IR := lc3.ir
	val := lc3.reg[(IR&DR)>>9]
	baseR := lc3.reg[(IR&BaseR)>>6]
	offset := signExt(IR&PCOFFSET6, 6)
	location := baseR + offset
	lc3.writeMemory(location, val)
}

func (lc3 *LC3vm) sti() {
	IR := lc3.ir
	val := lc3.reg[(IR&DR)>>9]
	location := lc3.readMemory(lc3.pc + signExt(IR&PCOFFSET9, 9))
	lc3.writeMemory(location, val)
}

func (lc3 *LC3vm) lea() {
	IR := lc3.ir
	val := lc3.pc + signExt(IR&PCOFFSET9, 9)
	lc3.reg[(IR&DR)>>9] = val
}

func (lc3 *LC3vm) updateCC(res uint16) {
	signedRes := int16(res) // convert to signed int
	if signedRes < 0 {
		lc3.nzp = 1 << 2
	} else if signedRes > 0 {
		lc3.nzp = 1
	} else {
		lc3.nzp = 1 << 1
	}
}

func signExt(val uint16, signedBitLocation int) uint16 {
	signedBitLocation--
	var sigBit uint16 = 1 << signedBitLocation
	sigMask := ^(sigBit - 1)
	if val&sigBit == sigBit {
		return val | sigMask
	}
	return val
}

// thank god for github copilot
