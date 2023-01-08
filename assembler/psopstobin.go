package assembler

func fillToBin() pseudoBinFunc {
	return func(line *string, addr *uint16, writer *objWriter) error {
		tokens := getTokens(*line)
		// strip label
		if !isKeyword(tokens[0]) {
			tokens = tokens[1:]
		}
		args := tokens[1:]
		if len(args) != 1 {
			return asmLineErr("invalid number of arguments for .FILL")
		}
		// only accepts number literals (no labels)
		val, err := strToUint16(args[0])
		if err != nil {
			return err
		}
		*addr++
		writer.writeUint16(val)
		return nil
	}
}

func blkwToBin() pseudoBinFunc {
	return func(line *string, addr *uint16, writer *objWriter) error {
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
	return func(line *string, addr *uint16, writer *objWriter) error {
		stringz, err := getStringzStr(*line)
		if err != nil {
			return err
		}
		for i := 0; i < len(stringz); i++ {
			writer.writeUint16(uint16(stringz[i]))
			*addr++
		}
		return nil
	}
}
