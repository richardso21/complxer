package lc3vm

import (
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
)

func (lc3 *LC3_st) readMemory(location uint16) uint16 {
	if location == KBSRADDR { // listen to keyboard once reading KBSR
		ch, controlKey, err := keyboard.GetSingleKey()
		if controlKey == keyboard.KeyCtrlC {
			lc3.HALT = true
			log.Fatal("Keyboard interrupt")
		}
		if err != nil {
			panic(err)
		}
		if ch != 0 {
			lc3.MEMORY[KBDRADDR] = uint16(ch)
			lc3.MEMORY[KBSRADDR] = 1 << 15
		} else {
			// avoid null terminator
			lc3.MEMORY[KBSRADDR] = 0
		}
	}
	return lc3.MEMORY[location]
}

func (lc3 *LC3_st) writeMemory(location uint16, value uint16) {
	if location == DDRADDR { // print to screen when writing DSR
		// print char
		fmt.Printf("%c", value)
	} else if location < 0x3000 || location >= 0xFE00 {
		lc3.HALT = true
		log.Fatalf("Accessing memory out of bounds, %04X", location)
	} else {
		lc3.MEMORY[location] = value
	}
}
