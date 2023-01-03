package asmlc3

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/richardso21/complxer/lc3vm"
)

func LoadASMFile(lc3 *lc3vm.LC3vm, file *os.File) (symTable, error) {
	sf := bufio.NewScanner(file)
	sf.Split(bufio.ScanLines)
	// check for empty/corrupt file first before doing anything
	if sf.Err() != nil {
		return nil, errors.New(sf.Err().Error())
	}
	table, err := getSymTable(sf) // generate symbol table
	if err != nil {
		return nil, err
	}

	file.Seek(0, io.SeekStart) // reset file pointer for second pass
	sf = bufio.NewScanner(file)
	err = assemble(lc3, sf, &table)

	return table, err
}
