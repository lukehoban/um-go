package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func run(program []uint32) {

	// State
	reg := [8]uint32{0, 0, 0, 0, 0, 0, 0, 0}
	platters := [][]uint32{program}
	freePlatters := []uint32{}
	var pc uint32

	// Analytics
	iteration := int64(0)
	startTime := time.Now()
	defer func() {
		duration := time.Now().Sub(startTime)
		fmt.Printf("\n\n** UM Execution Statistics **\n")
		fmt.Printf("Ops/s: %d\n", (iteration*1e9)/duration.Nanoseconds())
		totalArrayLength := 0
		for _, arr := range platters {
			totalArrayLength += len(arr)
		}
		fmt.Printf("Platters: %v, Free:  %v, TotalBytes: %v\n", len(platters), len(freePlatters), totalArrayLength*4)
	}()

	// Spin cycle
	for {
		iteration++
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
				if err != nil {
					if err == io.EOF {
						reg[c] = uint32(0xffffffff)
					} else {
						panic("Failed to read from stdin.")
					}
				}
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

func readPlatters(path string) ([]uint32, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	platters := make([]uint32, len(b)/4)
	err = binary.Read(bytes.NewReader(b), binary.BigEndian, &platters)
	if err != nil {
		return nil, err
	}
	return platters, nil
}

func main() {
	program := flag.String("program", "sandmark.umz", "The program to run on the Universal Machine.")
	flag.Parse()

	platters, err := readPlatters(*program)
	if err != nil {
		panic("Could not read program")
	}

	run(platters)
}
