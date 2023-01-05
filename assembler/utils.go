package assembler

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var __asmScanner *asmScanner

// global variables to track line for better error messages
// TODO: extend Scanner struct to include these variables/methods
// var __currentLine string
// var __currentTokens []string
// var __lineNumber int

// // use nextLine() and getCurrLine() in place for interacting with scanner object
// func nextLine() bool {
// 	// get next line of assembly file scanner
// 	if !__asmScanner.Scan() {
// 		return false // no more lines or error
// 	}
// 	getCurrTokens()
// 	__lineNumber++ // update global line number
// 	return true
// }

// func currLine() string {
// 	line := __asmScanner.Text()
// 	__currentLine = line // update global current line
// 	return line
// }

// func getCurrTokens() []string {
// 	__currentTokens = getTokens(currLine()) // alias to get tokens from current line
// 	return __currentTokens
// }

func getTokens(line string) []string {
	line = strings.Split(line, ";")[0] // strip comments
	if line == "" {
		return []string{} // empty line has no tokens
	}
	return splitByComma(strings.ToUpper(line)) // split by commas and whitespace, uppercased
}

func splitByComma(line string) []string {
	// split by commas and whitespace (common operation for LC3 token parsing)
	return splitByDelim(line, ',')
}

func splitByDelim(line string, delims ...rune) []string {
	// split given line w/ delimiters and whitespace
	return strings.FieldsFunc(line, getSplitFunc(delims...))
}

func getSplitFunc(delims ...rune) func(rune) bool {
	// return a function that splits by delimiter parameters and whitespace
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
		return "", asmLineErr("string not found (missing double quotes?)")
	}
	if after[len(after)-1] != '"' {
		return "", asmLineErr("string not terminated")
	}
	return after[:len(after)-1], nil
}

// convert (hex 0x, binary 0b, decimal #) string to uint16
func strToUint16(strNum string) (uint16, error) {
	var res uint64
	var val int64
	var err error
	if strNum[0] == 'X' {
		res, err = strconv.ParseUint(strNum[1:], 16, 16)
	} else if strNum[0] == 'B' {
		res, err = strconv.ParseUint(strNum[1:], 2, 16)
	} else if strNum[0] == '#' {
		// if not negative
		if strNum[1] != '-' {
			res, err = strconv.ParseUint(strNum[1:], 10, 16)
		} else {
			val, err = strconv.ParseInt(strNum[1:], 10, 16)
			res = uint64(val)
		}
	} else {
		if strNum[0] != '-' {
			res, err = strconv.ParseUint(strNum, 10, 16)
		} else {
			val, err = strconv.ParseInt(strNum, 10, 16)
			res = uint64(val)
		}
	}
	if err != nil {
		return 0, asmLineErr(err.Error())
	} else if int16(res) > 32767 || int16(res) < -32768 {
		return 0, asmLineErr("value is out of bounds: " + strNum)
	}
	return uint16(res), nil
}

func asmLineErr(errMessage string) error {
	return fmt.Errorf("assembler ERROR @ line #%d: \n==>\t%s\n%s",
		__asmScanner.currentLineNum, __asmScanner.currentLine, errMessage)
}

func asmGlobalErr(errMessage string) error {
	return fmt.Errorf("assembler ERROR: %s", errMessage)
}
