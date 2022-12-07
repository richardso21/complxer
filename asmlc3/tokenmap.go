package asmlc3

var noArgOpMap = map[string]uint16{
	"GETC": 0xF020,
	"OUT":  0xF021,
	"PUTS": 0xF022,
	"IN":   0xF023,
	"HALT": 0xF025,
	"RTI":  0x8000,
	"RET":  0xC1C0,
}

type tokenFunc func(*[]string, *symTable, uint16) (uint16, error)

var oneArgOpMap = map[string]tokenFunc{
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
	"TRAP":  trapToBin(),
}

var twoArgOpMap = map[string]tokenFunc{
	"LD":  ldToBin(),
	"LDI": ldiToBin(),
	"ST":  stToBin(),
	"STI": stiToBin(),
	"LEA": leaToBin(),
	"NOT": notToBin(),
}

var threeArgOpMap = map[string]tokenFunc{
	"ADD": addToBin(),
	"AND": andToBin(),
	"LDR": ldrToBin(),
	"STR": strToBin(),
}
