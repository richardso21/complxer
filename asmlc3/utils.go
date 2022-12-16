package asmlc3

import (
	"bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func getOrigAddr(line string) (uint16, error) {
	args := splitByDelim(line)
	if len(args) != 2 || args[0] != ".ORIG" {
		return 0, errors.New("invalid .ORIG first line, must be in format: .ORIG x3000")
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
		return strings.ContainsRune(string(delims), r) || unicode.IsSpace(r)
	}
}

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
		return 0, assemblerErr(err.Error())
	} else if int16(res) > 32767 || int16(res) < -32768 {
		return 0, assemblerErr("value is out of bounds: " + strNum)
	}
	return uint16(res), nil
}

func assemblerErr(e string) error {
	// currentLine from
	// return errors.New(fmt.Sprintln("line ", currentLine, ": ", e))
	return fmt.Errorf("[line %d: %s]\nERROR: %s", currentLine, getLine(s), e)
}
