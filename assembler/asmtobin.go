package assembler

func asmToBin(scanner *asmScanner, writer *objWriter, table symTable, addr uint16) error {
	// track the scanner object for easier error handling
	__asmScanner = scanner

	// constantly loop line read
	for ok, numTokens := scanner.getNextLine(); ok; ok, numTokens = scanner.getNextLine() {
		if numTokens == 0 {
			continue // skip empty lines
		}
		tokens := scanner.currentTokens

		// check if label as first token
		if !isKeyword(tokens[0]) {
			if len(tokens) == 1 {
				continue // if label is only token, read next line
			}
			tokens = tokens[1:] // deal with tokens without leading label
		}

		op := tokens[0]
		args := tokens[1:]

		// test if op is any of the (pseudo)ops
		switch op {
		case ".ORIG":
			// orig address already written in obj file
			continue
		case ".END":
			// do not evaluate anything after .END
			// (valid .END is already checked in getSymTable)
			goto FINISH
		case ".FILL", ".BLKW", ".STRINGZ":
			// use func in map to take care of these pseudo ops
			if err := pseudoOpToBinMap[op](&scanner.currentLine, &table, &addr, writer); err != nil {
				return err
			}
			// no need to write or increment addr, done in func call
			// writer.writeUint16(bin)
		default:
			// check if op is valid in map
			opToBin, ok := opToBinMap[op]
			if !ok {
				// op doesn't match with any op or pseudo op
				return asmLineErr("invalid operator: " + op)
			}
			// change operator to binary representation
			bin, err := opToBin(&args, &table, addr)
			if err != nil {
				return err
			}
			// write binary to object file, and increment addr
			writer.writeUint16(bin)
			addr++
		}
	}
FINISH:
	return nil
}
