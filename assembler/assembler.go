package assembler

import (
	"bufio"
	"os"
)

func AsmToObj(asmFile *os.File) error {
	// create a line-by-line scanner for assembly file
	asmScanner := newAsmScanner(asmFile)

	// create object file writer
	objFN := asmFile.Name() + ".obj" // object file name
	objFile, err := os.Create(objFN) // create file representation on disk (empties existing file)
	if err != nil {
		// something went wrong with creating file
		// TODO
		panic(err)
	}
	objWriter := bufio.NewWriter(objFile) // create writer for object file

	// perform first pass (symbol table)
	table, err := getSymTable(asmScanner)
	if err != nil {
		return err
	}

	// perform second pass (assembly)
	asmFile.Seek(0, 0)                  // reset file pointer
	asmScanner = newAsmScanner(asmFile) // recreate scanner for file

}
