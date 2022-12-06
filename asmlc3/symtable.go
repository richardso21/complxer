package asmlc3

import "bufio"

type symTable map[string]uint16

func getSymTable(s *bufio.Scanner) (symTable, error) {
	sTable := make(map[string]uint16)
	s.Scan() // get first line
	// skip comments until .ORIG
	for getLine(s) == "" {
		s.Scan()
	}
	// try getting starting address
	addr, err := getOrigAddr(getLine(s))
	if err != nil {
		return nil, err
	}

	for s.Scan() {
		line := getLine(s)
		// skip empty lines and comments
		if line == "" {
			continue
		}
		args := splitByDelim(line, ' ', ',')
		if !(isKeyword(args[0])) {
			// add to symbol table
			sTable[args[0]] = addr
			// check if label is with instruction
			if len(args) > 1 {
				addr++
			}
			// otherwise, label refers to instruction afterwards
		} else {
			// normal instruction line without label, increment address
			addr++
		}
	}
	return sTable, nil
}
