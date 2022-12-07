package asmlc3

import "golang.org/x/exp/slices"

func brToBin(n bool, z bool, p bool) tokenFunc {
	var nzp uint16 = 0
	if n {
		nzp = nzp | 0b100
	}
	if z {
		nzp = nzp | 0b010
	}
	if p {
		nzp = nzp | 0b001
	}
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		offset, err := getOffset((*args)[0], st, addr, 9)
		if err != nil {
			return 0, err
		}
		return nzp<<9 | offset, nil
	}
}

func jmpToBin() tokenFunc {
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		// jmp argument must be register
		if !isReg((*args)[0]) {
			return 0, assemblerErr("invalid JMP argument: " + (*args)[0])
		}
		reg := getReg((*args)[0])
		return 0xC000 | (reg << 6), nil
	}
}

func jsrToBin() tokenFunc {
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		offset, err := getOffset((*args)[0], st, addr, 11)
		if err != nil {
			return 0, err
		}
		return 0x4800 | offset, nil
	}
}

func jsrrToBin() tokenFunc {
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		if !isReg((*args)[0]) {
			return 0, assemblerErr("invalid JSRR argument: " + (*args)[0])
		}
		reg := getReg((*args)[0])
		return 0x4000 | (reg << 6), nil
	}
}

func trapToBin() tokenFunc {
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		val, err := strToUint16((*args)[0])
		if err != nil {
			return 0, assemblerErr("invalid TRAP argument: " + (*args)[0])
		}
		// check if trap vector is actually implemented
		trapVals := []uint16{0x20, 0x21, 0x22, 0x23, 0x25}
		if !slices.Contains(trapVals, val) {
			return 0, assemblerErr("invalid TRAP vector: " + (*args)[0])
		}
		return 0xF000 | val, nil
	}
}

func ldToBin() tokenFunc {
	return dROffset9Func(0x2000)
}

func ldiToBin() tokenFunc {
	return dROffset9Func(0xA000)
}

func stToBin() tokenFunc {
	return dROffset9Func(0x3000)
}

func stiToBin() tokenFunc {
	return dROffset9Func(0xB000)
}

func leaToBin() tokenFunc {
	return dROffset9Func(0xE000)
}

func dROffset9Func(opCode uint16) tokenFunc {
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		offset, err := getOffset((*args)[1], st, addr, 9)
		if err != nil {
			return 0, err
		}
		if !isReg((*args)[0]) {
			return 0, assemblerErr("invalid LD argument: " + (*args)[0])
		}
		destReg := getReg((*args)[0])
		return opCode | (destReg << 9) | offset, nil
	}
}

func notToBin() tokenFunc {
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		if !isReg((*args)[0]) {
			return 0, assemblerErr("invalid NOT argument: " + (*args)[0])
		} else if !isReg((*args)[1]) {
			return 0, assemblerErr("invalid NOT argument: " + (*args)[1])
		}
		destReg := getReg((*args)[0])
		srcReg := getReg((*args)[1])
		return 0x903F | (destReg << 9) | (srcReg << 6), nil
	}
}

func addToBin() tokenFunc {
	return aluToBin(0x1000)
}

func andToBin() tokenFunc {
	return aluToBin(0x5000)
}

func aluToBin(opCode uint16) tokenFunc {
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		if !isReg((*args)[0]) {
			return 0, assemblerErr("invalid ALU argument: " + (*args)[0])
		} else if !isReg((*args)[1]) {
			return 0, assemblerErr("invalid ALU argument: " + (*args)[1])
		}
		destReg := getReg((*args)[0])
		baseReg := getReg((*args)[1])
		if isReg((*args)[2]) {
			sr2 := getReg((*args)[2])
			return opCode | (destReg << 9) | (baseReg << 6) | sr2, nil
		} else {
			imm5, err := getImm5((*args)[2])
			if err != nil {
				return 0, err
			}
			return opCode | (destReg << 9) | (baseReg)<<6 | (0b1 << 5) | imm5, nil
		}
	}
}

func getImm5(strNum string) (uint16, error) {
	imm5, err := strToUint16(strNum)
	if err != nil {
		return 0, err
	} else if int16(imm5) > 15 || int16(imm5) < -16 {
		return 0, assemblerErr("IMM5 out of bounds: " + strNum)
	}
	return imm5 & 0x001F, nil // mask out all but the last 5 bits
}

func ldrToBin() tokenFunc {
	return dBaseROffset6Func(0x6000)
}

func strToBin() tokenFunc {
	return dBaseROffset6Func(0x7000)
}

func dBaseROffset6Func(opCode uint16) tokenFunc {
	return func(args *[]string, st *symTable, addr uint16) (uint16, error) {
		offset, err := getOffset((*args)[2], st, addr, 6)
		if err != nil {
			return 0, err
		}
		if !isReg((*args)[0]) {
			return 0, assemblerErr("invalid LD argument: " + (*args)[0])
		} else if !isReg((*args)[1]) {
			return 0, assemblerErr("invalid LD argument: " + (*args)[1])
		}
		destReg := getReg((*args)[0])
		baseReg := getReg((*args)[1])
		return opCode | (destReg << 9) | (baseReg << 6) | offset, nil
	}
}

func getOffset(arg string, st *symTable, addr uint16, offsetSize uint16) (uint16, error) {
	// addr + 1 because of PC*
	addr++
	// check if offset is a label
	if val, ok := (*st)[arg]; ok {
		offset, err := getOffsetDist(addr, val, offsetSize)
		if err != nil {
			return 0, err
		}
		return offset, nil
	} else {
		// offset is a number
		val, err := strToUint16(arg)
		if err != nil {
			return 0, assemblerErr("invalid value: " + arg)
		}
		return val, nil
	}
}

func getOffsetDist(addrFrom uint16, addrTo uint16, offsetSize uint16) (uint16, error) {
	// if offsetsize 16, then do not mask or check
	if offsetSize == 16 {
		return uint16(addrTo - addrFrom), nil
	}
	offset := int16(addrTo - addrFrom)
	if offset > (1<<offsetSize)-1 || offset < -(1<<offsetSize) {
		return 0, assemblerErr("offset out of range")
	}
	sigMask := uint16(1<<(offsetSize) - 1)
	return uint16(offset) & sigMask, nil
}
