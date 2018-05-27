package disassembly

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/cbush06/disassembler8080/opcodes"
)

// LittleEndian : Should we assume little-endian?
var LittleEndian = true

// Disassemble accepts a map of OpCode structures and disassembles a user-specified binary file using that map. All output is written to stdin.
func Disassemble(opcodesMap map[uint8]opcodes.OpCode) {
	//var filePath = getFilePath()
	var filePath = "INVADERS.H"
	var readerPos uint64
	var unkOpcodes = false
	var unkOpcodeStart uint64
	var fileHandle, _ = os.Open(filePath)
	var fileReader = bufio.NewReader(fileHandle)
	var tabWriter = tabwriter.NewWriter(os.Stdout, 12, 4, 1, ' ', 0)

	for {
		var opcode, err = fileReader.ReadByte()
		var opcodeStruct opcodes.OpCode

		if err != nil {
			break
		}

		readerPos++

		// Check if next code is recognized
		opcodeStruct, ok := opcodesMap[opcode]
		if !ok {
			// If it's unrecognized and it's the first unrecognized byte
			// in the current series of unknown codes, record the position
			if !unkOpcodes {
				unkOpcodeStart = readerPos
				unkOpcodes = true
			}

			continue
		}

		// If this is the first recognized opcode after a series of unknown codes,
		// report the start and end byte positions of the unrecognized block of bytes
		if unkOpcodes {
			fmt.Printf("\nUnrecognized bytes: %d through %d\n\n", unkOpcodeStart, readerPos)
			unkOpcodeStart = 0
			unkOpcodes = false
		}

		// Output the current OpCode name
		fmt.Fprint(tabWriter, "\n", opcodeStruct.Name)

		// If this opcode has arguments, print them
		if opcodeStruct.Length > 1 {
			var args = make([]byte, opcodeStruct.Length-1)
			if _, err := io.ReadAtLeast(fileReader, args, int(opcodeStruct.Length-1)); err != nil {
				log.Println(err)
			}

			// Progress the variable tracking our reader position
			readerPos += uint64(opcodeStruct.Length - 1)

			fmt.Fprintf(tabWriter, "\t")

			// If LittleEndian, print the args in reverse, otherwise, print them in order
			if LittleEndian {
				for i := len(args) - 1; i > -1; i-- {
					fmt.Fprintf(tabWriter, "\t0x%02x", args[i])
				}
			} else {
				for _, e := range args {
					fmt.Fprintf(tabWriter, "\t0x%02x", e)
				}
			}
		}

		tabWriter.Flush()
	}
}

func getFilePath() string {
	var consoleReader = bufio.NewReader(os.Stdin)
	var path string

	for {
		fmt.Print("Binary File: ")

		path, _ = consoleReader.ReadString('\n')
		path = strings.TrimSpace(path)

		// Confirm file exists
		if _, err := os.Lstat(path); err != nil {
			log.Println("File does not exist!")
			continue
		}

		break
	}

	return path
}
