package assembler

import (
	"os"
	"path/filepath"
)

// writes an object file from a given assembly file
func AsmToObj(asmFile *os.File) (string, error) {
	// create a line-by-line scanner for assembly file
	asmScanner := newAsmScanner(asmFile)
	asmFN := asmFile.Name()

	// create object file writer on "out" directory
	if err := os.Mkdir("out", 0777); err != nil {
		return "", err
	}
	objFN := "out/" + filepath.Base(asmFN) + ".obj" // object file name
	objFile, err := os.Create(objFN)                // create file representation on disk (empties existing file)
	if err != nil {
		// something went wrong with creating file
		return "", err
	}
	objWriter := newObjWriter(objFile) // create writer for object file

	// perform first pass (symbol table, check .ORIG & .END)
	table, origAddr, err := getSymTable(asmScanner)
	if err != nil {
		return "", err
	}

	// write origin address as first 16-bit word of object file
	objWriter.writeUint16(origAddr)

	// perform second pass (assembly)
	asmFile.Seek(0, 0)                  // reset file pointer
	asmScanner = newAsmScanner(asmFile) // recreate scanner for file

	err = asmToBin(asmScanner, objWriter, table, origAddr)
	if err != nil {
		return "", err
	}
	// save changes to writer
	err = objWriter.Flush()
	if err != nil {
		return "", err
	}
	objFile.Close()
	// return back reference to object file
	return objFN, nil
}
