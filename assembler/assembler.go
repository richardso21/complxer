package assembler

import "os"

// writes an object file from a given assembly file
func AsmToObj(asmFile *os.File) error {
	// create a line-by-line scanner for assembly file
	asmScanner := newAsmScanner(asmFile)

	// create object file writer
	objFN := asmFile.Name() + ".obj" // object file name
	objFile, err := os.Create(objFN) // create file representation on disk (empties existing file)
	if err != nil {
		// something went wrong with creating file
		return err
	}
	objWriter := newObjWriter(objFile) // create writer for object file

	// perform first pass (symbol table, check .ORIG & .END)
	table, origAddr, err := getSymTable(asmScanner)
	if err != nil {
		return err
	}

	// write origin address as first 16-bit word of object file
	objWriter.writeUint16(origAddr)

	// perform second pass (assembly)
	asmFile.Seek(0, 0)                  // reset file pointer
	asmScanner = newAsmScanner(asmFile) // recreate scanner for file

	err = asmToBin(asmScanner, objWriter, table, origAddr)
	return err
}
