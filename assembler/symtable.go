package assembler

type symTable map[string]uint16

func getSymTable(scanner *asmScanner) (symTable, uint16, error) {
	__asmScanner = scanner
	table := make(symTable)

	ok, numTokens := scanner.getNextLine()
	// constantly loop line read until first real line is found
	for ok && numTokens == 0 {
		ok, numTokens = scanner.getNextLine()
		if !ok {
			// EOF reached w/o .ORIG or file/scanner error
			return nil, 0, asmGlobalErr("no .ORIG found")
		}
	}
	// extract origin address
	origAddr, err := getOrigAddr(scanner.currentTokens)
	addr := origAddr
	if err != nil {
		return nil, 0, err
	}

	// keep looping until EOF
	for ok, numTokens = scanner.getNextLine(); ok; ok, numTokens = scanner.getNextLine() {
		if numTokens == 0 {
			continue // skip empty lines
		}
		tokens := scanner.currentTokens
		if numTokens == 1 && tokens[0] == ".END" {
			// successful once .END is found
			return table, origAddr, nil
		}
		if !isKeyword(tokens[0]) {
			// add to symbol table
			table[tokens[0]] = addr
			// check if label is with instruction
			if len(tokens) > 1 {
				switch tokens[1] {
				case ".BLKW", ".STRINGZ":
					// increment address arbitrarily by performing pseudo op with dummy writer
					if err := pseudoOpToBinMap[tokens[1]](
						&scanner.currentLine, &table, &addr, newObjWriter(nil)); err != nil {
						return nil, 0, err
					}
				default:
					addr++
				}
			}
			// else label refers to instruction afterwards, so don't increment address
		} else {
			// normal instruction w/o label (or label at preceding line), increment address
			addr++
		}
	}

	// .END not found or is invalid
	return nil, 0, asmGlobalErr(".END not found or is invalid")
}

func getOrigAddr(tokens []string) (uint16, error) {
	if tokens[0] != ".ORIG" {
		return 0, asmGlobalErr("first line must define .ORIG address")
	}
	if len(tokens) != 2 {
		return 0, asmLineErr("incorrect number of arguments for .ORIG line (expected 1)")
	}
	return strToUint16(tokens[1])
}
