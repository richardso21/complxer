package lc3vm

import (
	"fmt"
	"log"

	"github.com/eiannone/keyboard"
)

var trapFuncs = [...]func(){
	LC3.tGetC,
	LC3.tOut,
	LC3.tPutS,
	LC3.tIn,
	func() {}, // 0x24 not used
	LC3.tHalt,
}

var trapFuncMap = make(map[uint16]func())

func init() {
	for i, fn := range trapFuncs {
		trapFuncMap[uint16(i)+0x20] = fn
	}
}

func (lc3 *LC3vm) trap() {
	// get trap vector, then execute respective function
	trapFuncMap[lc3.ir&0x00FF]()
}

func (lc3 *LC3vm) tGetC() {
	ch, controlKey, err := keyboard.GetSingleKey()
	if controlKey == keyboard.KeyCtrlC {
		lc3.halt = true
		log.Fatal("Keyboard interrupt")
	}
	if err != nil {
		panic(err)
	}
	if ch != 0 {
		lc3.reg[0] = uint16(ch)
	} else {
		// avoid null terminator
		lc3.reg[0] = 0
	}
}

func (lc3 *LC3vm) tOut() {
	val := lc3.reg[0]
	fmt.Printf("%c", val)
}

func (lc3 *LC3vm) tPutS() {
	addr := lc3.reg[0]
	for lc3.mem[addr] != 0 {
		fmt.Printf("%c", lc3.mem[addr])
		addr++
	}
}

func (lc3 *LC3vm) tIn() {
	fmt.Print("IN: ")
	ch, controlKey, err := keyboard.GetSingleKey()
	if controlKey == keyboard.KeyCtrlC {
		lc3.halt = true
		log.Fatal("Keyboard interrupt")
	}
	if err != nil {
		panic(err)
	}
	if ch != 0 {
		lc3.reg[0] = uint16(ch)
		fmt.Printf("\n%c\n", ch)

	} else {
		// avoid null terminator
		lc3.reg[0] = 0
	}
}

func (lc3 *LC3vm) tHalt() {
	// halt program
	lc3.halt = true
	// syscall.Exit(0)

}
