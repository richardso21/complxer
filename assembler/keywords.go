package assembler

import (
	"golang.org/x/exp/slices"
)

var opNames = [...]string{
	"BR",
	"BRN",
	"BRZ",
	"BRP",
	"BRNZ",
	"BRNP",
	"BRZP",
	"BRNZP",
	"ADD",
	"LD",
	"ST",
	"JSR",
	"JSRR",
	"AND",
	"LDR",
	"STR",
	"RTI",
	"NOT",
	"LDI",
	"STI",
	"JMP",
	"RET",
	"LEA",
	"TRAP",
}

func isOp(s string) bool {
	return slices.Contains(opNames[:], s)
}

var trapNames = [...]string{
	"GETC",
	"OUT",
	"PUTS",
	"IN",
	// "PUTSP", nobody uses PUTSP :P
	"HALT",
}

func isTrap(s string) bool {
	return slices.Contains(trapNames[:], s)
}

var pseudoOpNames = [...]string{
	".ORIG",
	".FILL",
	".BLKW",
	".STRINGZ",
	".END",
}

func isPseudoOp(s string) bool {
	return slices.Contains(pseudoOpNames[:], s)
}

func isKeyword(s string) bool {
	return isOp(s) || isPseudoOp(s) || isTrap(s)
}

var regNames = [...]string{
	"R0",
	"R1",
	"R2",
	"R3",
	"R4",
	"R5",
	"R6",
	"R7",
}

func isReg(s string) bool {
	return slices.Contains(regNames[:], s)
}

func getReg(s string) uint16 {
	return uint16(slices.Index(regNames[:], s))
}
