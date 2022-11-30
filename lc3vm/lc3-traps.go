package lc3vm

import (
	"fmt"
)

const (
	_TRAP_GETC = 0x20
	_TRAP_OUT  = 0x21
	_TRAP_PUTS = 0x22
	_TRAP_IN   = 0x23
	_TRAP_HALT = 0x25
)

var TRAP_FUNCMAP = map[uint16]func(){
	_TRAP_GETC: LC3._trap_getc,
	_TRAP_OUT:  LC3._trap_out,
	_TRAP_PUTS: LC3._trap_puts,
	_TRAP_IN:   LC3._trap_in,
	_TRAP_HALT: LC3._trap_halt,
}

func (lc3 *LC3_st) _rti() {
	// not implemented
}

func (lc3 *LC3_st) _trap() {
	// get trap vector, then execute respective function
	TRAP_FUNCMAP[lc3.IR&0x00FF]()
}

func (lc3 *LC3_st) _trap_getc() {
	var input string
	fmt.Scanf("%c", &input)
	lc3.REG[0] = uint16(input[0])
}

func (lc3 *LC3_st) _trap_out() {
	addr := lc3.REG[0]
	fmt.Printf("%c", lc3.MEMORY[addr])
}

func (lc3 *LC3_st) _trap_puts() {
	addr := lc3.REG[0]
	for lc3.MEMORY[addr] != 0 {
		fmt.Printf("%c", lc3.MEMORY[addr])
		addr++
	}
}

func (lc3 *LC3_st) _trap_in() {
	var input string
	fmt.Print("IN: ")
	fmt.Scanf("%c", &input)
	lc3.REG[0] = uint16(input[0])
	fmt.Printf("%c", input[0])
}

func (lc3 *LC3_st) _trap_halt() {
	// halt program
	lc3.HALT = true
	// syscall.Exit(0)

}
