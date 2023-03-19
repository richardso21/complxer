package assembler

func fillToBin() pseudoBinFunc {
	return func(line *string, table *symTable, addr *uint16, writer *objWriter) error {
		tokens := getTokens(*line)
		// strip label
		if !isKeyword(tokens[0]) {
			tokens = tokens[1:]
		}
		args := tokens[1:]
		if len(args) != 1 {
			return asmLineErr("invalid number of arguments for .FILL")
		}
		if val, ok := (*table)[args[0]]; ok {
			// if symbol is label, write where label points to
			writer.writeUint16(val)
		} else {
			// check if its a valid number and write
			val, err := strToUint16(args[0])
			if err != nil {
				return err
			}
			writer.writeUint16(val)
		}
		// increment address
		*addr++
		return nil
	}
}

func blkwToBin() pseudoBinFunc {
	return func(line *string, table *symTable, addr *uint16, writer *objWriter) error {
		tokens := getTokens(*line)
		if !isKeyword(tokens[0]) {
			tokens = tokens[1:]
		}
		args := tokens[1:]
		if len(args) != 1 {
			return asmLineErr("invalid number of arguments for .BLKW")
		}
		val, err := strToUint16(args[0])
		if err != nil {
			return err
		}
		// skip `val` amount of steps
		for i := uint16(0); i < val; i++ {
			*addr++
		}
		return nil
	}
}

func stringzToBin() pseudoBinFunc {
	return func(line *string, table *symTable, addr *uint16, writer *objWriter) error {
		stringz, err := getStringzStr(*line)
		if err != nil {
			return err
		}
		for i := 0; i < len(stringz); i++ {
			// escape forward slash
			if stringz[i] == '\\' {
				i++
				if stringz[i] == 'n' {
					writer.writeUint16('\n')
				} else if stringz[i] == 't' {
					writer.writeUint16('\t')
				} else if stringz[i] == '\\' {
					writer.writeUint16('\\')
				} else {
					return asmLineErr("invalid escape sequence")
				}
				// thanks github copilot :D
			} else {
				writer.writeUint16(uint16(stringz[i]))
			}
			*addr++
		}
		// null terminator
		writer.writeUint16(0)
		*addr++
		return nil
	}
}
