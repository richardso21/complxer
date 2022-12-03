package lc3vm

func (lc3 *LC3vm) Mem() *[1 << 16]uint16 {
	// pass by reference to prevent expensive value copying
	return &lc3.mem
}

func (lc3 *LC3vm) Reg() *[8]uint16 {
	return &lc3.reg
}

func (lc3 *LC3vm) Pc() uint16 {
	return lc3.pc
}

func (lc3 *LC3vm) Ir() uint16 {
	return lc3.ir
}

func (lc3 *LC3vm) IsHalt() bool {
	return lc3.halt
}

func (lc3 *LC3vm) incrStack() {
	lc3.stackSize++
}

func (lc3 *LC3vm) decrStack() {
	if lc3.stackSize == 0 {
		return
	}
	lc3.stackSize--
}
