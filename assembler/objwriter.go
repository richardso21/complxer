package assembler

import (
	"bufio"
	"os"
)

// extended bufio.Writer that allows for writing uint16 values
type objWriter struct {
	*bufio.Writer
}

func newObjWriter(f *os.File) *objWriter {
	return &objWriter{bufio.NewWriter(f)}
}

func (w *objWriter) writeUint16(u uint16) (int, error) {
	// separate uint16 into two bytes and write both to file
	return w.Write([]byte{byte(u >> 8), byte(u)})
}
