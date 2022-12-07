package asmlc3

import (
	"bufio"
	"strings"

	"github.com/richardso21/complxer/lc3vm"
)

// keep track of current line number for error messaging
var currentLine int = 1

func loadOnLC3(lc3 *lc3vm.LC3vm, s *bufio.Scanner, st *symTable) error {
	// skip commented lines
	s.Scan()
	currentLine = 1
	for getLine(s) == "" {
		s.Scan()
		currentLine++
	}
	// get initial address
	addr, err := getOrigAddr(getLine(s))
	if err != nil {
		return err
	}
	// keep getting new line until .END
	for s.Scan() && getLine(s) != ".END" {
		currentLine++
		line := getLine(s)
		// skip empty/comment lines
		if line == "" {
			continue
		}
		tokens := splitByDelim(line, ',')
		var binInstr uint16
		var err error
		if isPseudoOp(tokens[0]) { // check if pseudo op
			err := pseudoOpToBin(lc3, tokens, st, &addr)
			if err != nil {
				return err
			}
			continue
		} else if !(isKeyword(tokens[0])) { // check if label
			if len(tokens) == 1 {
				// label on its own line
				continue
			} else if !isKeyword(tokens[1]) {
				// two labels on one line, malformed instruction
				return assemblerErr("unknown instruction or two labels on one line" +
					tokens[1])
			} else if isPseudoOp(tokens[1]) {
				// pseudo op with label
				err := pseudoOpToBin(lc3, tokens[1:], st, &addr)
				if err != nil {
					return err
				}
				continue
			} else {
				// label with instruction
				binInstr, err = tokensToBin(tokens[1:], st, addr)
			}
		} else {
			// normal instruction line, increment address
			binInstr, err = tokensToBin(tokens, st, addr)
		}
		if err != nil {
			return err
		}
		lc3.FillValue(addr, binInstr)
		addr++
		// incr addr after each fill loop
	}

	if getLine(s) != ".END" {
		return assemblerErr("missing or malformed .END directive: " + getLine(s))
	}
	return nil
}

func tokensToBin(tokens []string, st *symTable, addr uint16) (uint16, error) {
	instr := tokens[0]
	args := tokens[1:]

	switch len(args) {
	case 0:
		// check on instructions with no args
		if val, ok := noArgOpMap[instr]; ok {
			return val, nil
		}
		return errTokenNotFound(instr)

	case 1:
		// check on instructions with one arg
		if fn, ok := oneArgOpMap[instr]; ok {
			// use corresponding function to get binary instruction
			val, err := fn(&args, st, addr)
			if err != nil {
				return 0, err
			}
			return val, nil
		}
		return errTokenNotFound(instr)

	case 2:
		// two args
		if fn, ok := twoArgOpMap[instr]; ok {
			val, err := fn(&args, st, addr)
			if err != nil {
				return 0, err
			}
			return val, nil
		}
		return errTokenNotFound(instr)

	case 3:
		// three args
		if fn, ok := threeArgOpMap[instr]; ok {
			val, err := fn(&args, st, addr)
			if err != nil {
				return 0, err
			}
			return val, nil
		}
		return errTokenNotFound(instr)

	default:
		return 0, assemblerErr("invalid number of arguments")
	}
}

func pseudoOpToBin(lc3 *lc3vm.LC3vm, tokens []string, st *symTable, addr *uint16) error {
	if len(tokens) != 2 && tokens[0] != ".STRINGZ" {
		return assemblerErr("invalid number of arguments for " + tokens[0])
	}
	switch tokens[0] {
	case ".FILL":
		// check if it is a label that we are filling
		if val, ok := (*st)[tokens[1]]; ok {
			offset, err := getOffset(tokens[1], st, *addr, 16)
			if err != nil {
				return err
			}
			lc3.FillValue(offset, val)
			*addr++
			return nil
		}
		val, err := strToUint16(tokens[1])
		if err != nil {
			return err
		}
		lc3.FillValue(*addr, val)
		*addr++
		return nil

	case ".BLKW":
		val, err := strToUint16(tokens[1])
		if err != nil {
			return err
		}
		var i uint16
		for i = 0; i < val; i++ {
			*addr++
		}
		return nil

	case ".STRINGZ":
		str := strings.Join(tokens[1:], " ")
		if str[0] != '"' || str[len(str)-1] != '"' {
			return assemblerErr("invalid string format")
		}
		str = str[1 : len(str)-1] // remove quotes
		for i := 0; i < len(str); i++ {
			lc3.FillValue(*addr, uint16(str[i]))
			*addr++
		}
		return nil

	default:
		return assemblerErr("invalid pseudo op: " + tokens[0])
	}
}

func errTokenNotFound(token string) (uint16, error) {
	if isKeyword(token) {
		return 0, assemblerErr("wrong number of arguments for " + token)
	} else {
		return 0, assemblerErr("invalid instruction: " + token)
	}
}
