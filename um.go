package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

func run(program []uint32) {

	reg := [8]uint32{0, 0, 0, 0, 0, 0, 0, 0}
	platters := [][]uint32{program}
	freePlatters := []uint32{}
	var pc uint32

	for {
		instruction := platters[0][pc]
		op := (instruction >> 28) & 15
		a := ((instruction >> 6) & 7)
		b := ((instruction >> 3) & 7)
		c := ((instruction >> 0) & 7)
		switch op {
		case 0:
			if reg[c] != 0 {
				reg[a] = reg[b]
			}
		case 1:
			reg[a] = platters[reg[b]][reg[c]]
		case 2:
			platters[reg[a]][reg[b]] = reg[c]
		case 3:
			reg[a] = reg[b] + reg[c]
		case 4:
			reg[a] = reg[b] * reg[c]
		case 5:
			reg[a] = reg[b] / reg[c]
		case 6:
			reg[a] = ^(reg[b] & reg[c])
		case 7:	
			return
		case 8:
			{
				newPlatter := make([]uint32, reg[c])
				if len(freePlatters) > 0 {
					platters[freePlatters[0]] = newPlatter
					reg[b] = freePlatters[0]
					freePlatters = freePlatters[1:]
				} else {
					platters = append(platters, newPlatter)
					reg[b] = uint32(len(platters) - 1)
				}
			}
		case 9:
			{
				platters[reg[c]] = nil
				freePlatters = append(freePlatters, reg[c])
			}
		case 10:
			os.Stdout.Write([]byte{byte(reg[c])})
		case 11:
			{
				b := []byte{0}
				_, err := os.Stdin.Read(b)
				check(err)
				reg[c] = uint32(b[0])
			}
		case 12:
			{
				if reg[b] != 0 {
					platters[0] = make([]uint32, len(platters[reg[b]]))
					copy(platters[0], platters[reg[b]])
				}
				pc = reg[c]
				continue
			}
		case 13:
			reg[(instruction>>25)&7] = instruction & 0x01FFFFFF
		default:
			panic(fmt.Errorf("Failed on %d", op))
		}
		pc++
	}
}

func readPlatters(path string) []uint32 {
	b, err := ioutil.ReadFile(path)
	check(err)
	platters := make([]uint32, len(b)/4)
	err = binary.Read(bytes.NewReader(b), binary.BigEndian, &platters)
	check(err)
	return platters
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func decompile(program []uint32) {
	for i := range program {
		instruction := program[i]
		op := (instruction >> 28) & 15
		a := ((instruction >> 6) & 7)
		b := ((instruction >> 3) & 7)
		c := ((instruction >> 0) & 7)
		var text string
		switch op {
		case 0:
			text = fmt.Sprintf("IF REG[%d] != 0 { REG[%d] = REG[%d] }", c, a, b)
		case 1:
			text = fmt.Sprintf("REG[%d] = PLATTERS[REG[%d]][REG[%d]]", a, b, c)
		case 2:
			text = fmt.Sprintf("PLATTERS[REG[%d]][REG[%d]] = REG[%d]", a, b, c)
		case 3:
			text = fmt.Sprintf("REG[%d] = REG[%d] + REG[%d]", a, b, c)
		case 4:
			text = fmt.Sprintf("REG[%d] = REG[%d] * REG[%d]", a, b, c)
		case 5:
			text = fmt.Sprintf("REG[%d] = REG[%d] / REG[%d]", a, b, c)
		case 6:
			text = fmt.Sprintf("REG[%d] = ^(REG[%d] & REG[%d])", a, b, c)
		case 7:
			text = fmt.Sprintf("HALT")
		case 8:
			text = fmt.Sprintf("REG[%d] = MALLOC(REG[%d])", b, c)
		case 9:
			text = fmt.Sprintf("ABND %d", c)
		case 10:
			text = fmt.Sprintf("OTPT %d", c)
		case 11:
			text = fmt.Sprintf("INPT %d", c)
		case 12:
			text = fmt.Sprintf("PLATTERS[0] = PLATTERS[REG[%d]]; GOTO REG[%d]", b, c)
		case 13:
			text = fmt.Sprintf("REG[%d] = %x", (instruction>>25)&7, instruction&0x01FFFFFFc)
		}
		fmt.Printf("[%08x] %08x: %s\n", i, instruction, text)
	}
}

func main() {
	program := flag.String("program", "sandmark.umz", "The program to run on the Universal Machine.")
	decomp := flag.Bool("decompile", false, "Decompile instead of execute.")
	flag.Parse()

	platters := readPlatters(*program)

	// If -decompile, decompile the program instead of running
	if *decomp {
		decompile(platters)
		return
	}

	// Else run
	run(platters)
}
