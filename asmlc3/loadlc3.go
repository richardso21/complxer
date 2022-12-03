package asmlc3

import (
	"bufio"

	"github.com/richardso21/complxer/lc3vm"
)

func loadOnLC3(lc3 *lc3vm.LC3vm, s *bufio.Scanner, addr uint16, symTable *symTable) error {
	for s.Scan() {
		// skip commented lines
		// for s.Scan() {
		// 	s.Scan()
		// }

	}
	return nil
}

func readASMLine(lc3 *lc3vm.LC3vm, s string, st *symTable, addr *uint16) {
	// (*st)["PC"] = *addr
	args := splitByDelim(s)
	if len(args) == 0 {
		return // empty line
	}
}
