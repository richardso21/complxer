package assembler

import (
	"bufio"
	"os"
)

// wrapper for bufio.Scanner w/ additional functionality (line tracking, reset, etc.)
type asmScanner struct {
	*bufio.Scanner
	currentLine    string
	currentTokens  []string
	currentLineNum int
}

func newAsmScanner(f *os.File) *asmScanner {
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines) // default to scan by line
	return &asmScanner{s, "", nil, 0}
}

// return status of line read, number of tokens in read line
func (s *asmScanner) getNextLine() (bool, int) {
	if !s.Scanner.Scan() {
		// EOF or error (check s.Scanner.Err)
		return false, 0
	}
	// update variables with current line
	s.currentLine = s.Scanner.Text()
	s.currentTokens = getTokens(s.currentLine)
	s.currentLineNum++
	return true, len(s.currentTokens) // successful line read
}

// // reset scanner mainly for 2nd pass of file
// func (s *asmScanner) resetScanner() {
// 	// set variables to default
// 	s.currentLine = ""
// 	s.currentTokens = nil
// 	s.currentLineNum = 0
// 	s.file.Seek(0, 0)                    // seek file obj back to beginning
// 	s.scanner = bufio.NewScanner(s.file) // create new scanner object
// }
