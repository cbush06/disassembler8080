package main

import (
	"github.com/cbush06/disassembler8080/disassembly"
	"github.com/cbush06/disassembler8080/opcodes"
)

func main() {
	var opcodes = opcodes.InitOpcodes()
	disassembly.Disassemble(opcodes)
}
