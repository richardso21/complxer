package asmlc3

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"unicode"
)

func getOrigAddr(line string) (uint16, error) {
	args := splitByDelim(line)
	if len(args) != 2 || args[0] != ".ORIG" {
		return 0, errors.New("invalid .ORIG line, must be in format: .ORIG x3000")
	}
	addr, err := strToUint16(args[1])
	if err != nil {
		return 0, err
	}
	return addr, nil
}

func getLine(s *bufio.Scanner) string {
	// remove comments (any text after ';')
	res := strings.Split(s.Text(), ";")[0]
	res = strings.TrimSpace(res) // trim whitespaces
	res = strings.ToUpper(res)   // convert to uppercase to match with keywords
	return res
}

func splitByDelim(line string, delims ...rune) []string {
	return strings.FieldsFunc(line, getSplitFunc(delims...))
}

func getSplitFunc(delims ...rune) func(rune) bool {
	if len(delims) == 0 {
		return unicode.IsSpace
	}
	return func(r rune) bool {
		return strings.ContainsRune(string(delims), r)
	}
}

func strToUint16(str string) (uint16, error) {
	if str[0] == 'x' {
		res, err := strconv.ParseUint(str[1:], 16, 16)
		return uint16(res), err
	} else if str[0] == 'b' {
		res, err := strconv.ParseUint(str[1:], 2, 16)
		return uint16(res), err
	} else if str[0] == '#' {
		res, err := strconv.ParseUint(str[1:], 10, 16)
		return uint16(res), err
	} else {
		res, err := strconv.ParseUint(str, 10, 16)
		return uint16(res), err
	}
}
