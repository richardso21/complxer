package asmlc3

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/richardso21/complxer/lc3vm"
)

func loadOnLC3(lc3 *lc3vm.LC3vm, s *bufio.Scanner, symTable *symTable) error {
	// skip commented lines
	s.Scan()
	lineNumber := 1
	if getLine(s) == "" {
		s.Scan()
		lineNumber++
	}
	// get initial address
	addr, err := getOrigAddr(getLine(s))
	if err != nil {
		return err
	}
	// keep getting new line until .END
	for s.Scan() && getLine(s) != ".END" {
		lineNumber++
		line := getLine(s)
		// skip empty/comment lines
		if line == "" {
			continue
		}
		args := splitByDelim(line, ' ', ',')
		var binInstr uint16
		var err error
		// check if there is label
		if !(isKeyword(args[0])) {
			if len(args) == 1 {
				// label on its own line
				continue
			} else if !isKeyword(args[1]) {
				// two labels on one line, malformed instruction
				return errors.New("malformed instruction: " + line)
			} else {
				// label with instruction
				binInstr, err = argsToBinary(args[1:], symTable, addr, lineNumber)
				addr++
			}
		} else {
			// normal instruction line, increment address
			binInstr, err = argsToBinary(args, symTable, addr, lineNumber)
			addr++
		}

		if err != nil {
			return err
		}
		lc3.FillValue(addr, binInstr)
	}

	if !s.Scan() {
		return errors.New("missing or malformed .END directive")
	}
	return nil
}

func argsToBinary(args []string, st *symTable, addr uint16, ln int) (uint16, error) {
	if len(args) < 1 || len(args) > 4 {
		return 0, errors.New(fmt.Sprintln("malformed instruction: ", args))
	}
	if len(args) == 1 {
		switch args[0] {
		case "RTI":
			return 0x8000, nil // RTI instruction
		case "RET":
			return 0xC1C0, nil // RET instruction
		case "GETC":
			return 0xF020, nil // GETC trap
		case "OUT":
			return 0xF021, nil // OUT  trap
		case "PUTS":
			return 0xF022, nil // PUTS trap
		case "IN":
			return 0xF023, nil // IN   trap
		case "HALT":
			return 0xF025, nil // HALT trap
		default:
			if isKeyword(args[0]) {
				return 0, assemblerErr(ln, "not enough arguments for "+args[0])
			} else {
				return 0, assemblerErr(ln, "invalid instruction: "+args[0])
			}
		}
	} else if len(args) == 3 {
		switch args[0] {
		case "NOT":
			if !isReg(args[1]) {
				return 0, assemblerErr(ln, "invalid NOT arguments: "+args[1])
			} else if !isReg(args[2]) {
				return 0, assemblerErr(ln, "invalid NOT arguments: "+args[2])
			} else {
				return 0x903F | (getReg(args[1]) << 9) | (getReg(args[2]) << 6), nil
			}
		case "LEA":
			if !isReg(args[1]) {
				return 0, assemblerErr(ln, "invalid LEA arguments: "+args[1])
			}
		}
	}

	return 0, nil
}

func assemblerErr(ln int, e string) error {
	return errors.New(fmt.Sprintln("line ", ln, ": ", e))
}
