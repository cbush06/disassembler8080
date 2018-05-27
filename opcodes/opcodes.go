package opcodes

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// OpCode contains the actual numeric code, commonly accepted name for the code, and the length in bytes of the opcode and its arguments.
type OpCode struct {
	Code   uint8
	Name   string
	Length uint8
}

// InitOpcodes loads opcodes from a user-specified CSV file. Format of the CSV file should e: hex opcode (e.g. 0x15), common name of the
// opcode, and length in bytes of the opcode plus its arguments.
func InitOpcodes() map[uint8]OpCode {
	//var opCodeFilePath = getFilePath()
	return parseCsvToOpcodes("opcodes.csv")
}

func parseCsvToOpcodes(filePath string) map[uint8]OpCode {
	var file, _ = os.Open(filePath)
	var reader = csv.NewReader(bufio.NewReader(file))
	var opcodes = make(map[uint8]OpCode)

	for {
		var csvLineTokens, err = reader.Read()
		if err == io.EOF {
			break
		}

		if opcode, err := strconv.ParseUint(csvLineTokens[0], 0, 8); err == nil {
			if length, err := strconv.ParseUint(csvLineTokens[2], 10, 8); err == nil {
				opcodes[uint8(opcode)] = OpCode{
					Code:   uint8(opcode),
					Name:   csvLineTokens[1],
					Length: uint8(length),
				}
			} else {
				log.Printf("Error parsing opcode size: %s\n", err.Error())
			}
		} else {
			log.Printf("Error parsing opcode hex: %s\n", err.Error())
		}
	}

	return opcodes
}

func getFilePath() string {
	var consoleReader = bufio.NewReader(os.Stdin)
	var path string

	for {
		fmt.Print("OpCode CSV File: ")

		path, _ = consoleReader.ReadString('\n')
		path = strings.TrimSpace(path)

		// Confirm file extension
		fileExt := strings.ToLower(path[len(path)-3:])
		if fileExt != "csv" {
			fmt.Println("Invalid file type!")
			continue
		}

		// Confirm file exists
		if _, err := os.Lstat(path); err != nil {
			log.Println("File does not exist!")
			continue
		}

		break
	}

	return path
}
