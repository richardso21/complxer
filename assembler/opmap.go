package assembler

type binFunc func(*[]string, *symTable, uint16) (uint16, error)

var opToBinMap = map[string]binFunc{
	"ADD":   addToBin(),
	"AND":   andToBin(),
	"BR":    brToBin(true, true, true), // same as BRnzp
	"BRN":   brToBin(true, false, false),
	"BRZ":   brToBin(false, true, false),
	"BRP":   brToBin(false, false, true),
	"BRNZ":  brToBin(true, true, false),
	"BRNP":  brToBin(true, false, true),
	"BRZP":  brToBin(false, true, true),
	"BRNZP": brToBin(true, true, true),
	"JMP":   jmpToBin(),
	"JSR":   jsrToBin(),
	"JSRR":  jsrrToBin(),
	"LD":    ldToBin(),
	"LDI":   ldiToBin(),
	"LDR":   ldrToBin(),
	"LEA":   leaToBin(),
	"NOT":   notToBin(),
	"ST":    stToBin(),
	"STI":   stiToBin(),
	"STR":   strToBin(),
	"TRAP":  trapToBin(),
}

type pseudoBinFunc func(*string, *uint16, *objWriter) error

var pseudoOpToBinMap = map[string]pseudoBinFunc{
	".FILL":    fillToBin(),
	".BLKW":    blkwToBin(),
	".STRINGZ": stringzToBin(),
}
