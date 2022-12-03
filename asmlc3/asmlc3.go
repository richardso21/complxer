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
	// check for empty file first before doing anything
	if !sf.Scan() {
		return nil, errors.New("empty file")
	}
	startAddr, table, err := getSymTable(sf) // generate symbol table
	if err != nil {
		return nil, err
	}

	file.Seek(0, io.SeekStart) // reset file pointer for second pass
	loadOnLC3(lc3, sf, startAddr, &table)

	return table, nil
}
