package asmlc3

import (
	"strings"

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
	"JSRR",
	"AND",
	"LDR",
	"STR",
	"RTI",
	"NOT",
	"LDI",
	"STI",
	"JMP",
	"LEA",
	"TRAP",
}

func isOp(s string) bool {
	return slices.Contains(opNames[:], strings.ToUpper(s))
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
	return slices.Contains(trapNames[:], strings.ToUpper(s))
}

var pseudoOpNames = [...]string{
	".ORIG",
	".FILL",
	".BLKW",
	".STRINGZ",
	".END",
}

func isPseudoOp(s string) bool {
	return slices.Contains(pseudoOpNames[:], strings.ToUpper(s))
}

func isKeyword(s string) bool {
	return isOp(s) || isPseudoOp(s) || isTrap(s)
}
