package assembler

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
)

// variables to track line for better error messages
var __asmScanner *bufio.Scanner
var __currentLine string
var __lineNumber int

// use nextLine() and getLine() in place for interacting with scanner object
func nextLine() bool {
	if !__asmScanner.Scan() {
		return false
	}
	__lineNumber++
	return true
}

func getLine() string {
	line := __asmScanner.Text()
	__currentLine = line // update current line
	return line
}

func getTokens(line string) []string {
	line = strings.Split(line, ";")[0] // strip comments
	return splitByComma(line)          // split by comma
}

func splitByComma(line string) []string {
	return splitByDelim(line, ',')
}

func splitByDelim(line string, delims ...rune) []string {
	return strings.FieldsFunc(line, getSplitFunc(delims...))
}

func getSplitFunc(delims ...rune) func(rune) bool {
	if len(delims) == 0 {
		return unicode.IsSpace
	}
	return func(r rune) bool {
		return strings.ContainsRune(string(delims), r) || unicode.IsSpace(r)
	}
}

func getStringzStr(line string) (string, error) {
	_, after, ok := strings.Cut(line, "\"") // get argument in double quotes
	// error handling (no string found or string not terminated)
	if !ok {
		return "", assemblerErr("string not found (missing double quotes?)")
	}
	if after[len(after)-1] != '"' {
		return "", assemblerErr("string not terminated")
	}
	return after[:len(after)-1], nil
}

func assemblerErr(errMessage string) error {
	// TODO
	return fmt.Errorf("assembler ERROR @ line #%d: \n==>\t%s\n%s",
		__lineNumber, __currentLine, errMessage)
}
