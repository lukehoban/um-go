package main

import (
    "io/ioutil"
    "fmt"
    "bytes"
    "encoding/binary"
)

func run(program []uint32) {
    reg := [8]uint32{0,0,0,0,0,0,0,0}
    platters := [][]uint32{program}
    var pc uint32 = 0
    
    for ;; {
        instruction := platters[0][pc]
        op := (instruction >> 28) & 15
        a := ((instruction >> 6) & 7)
        b := ((instruction >> 3) & 7) 
        c := ((instruction >> 0) & 7)
        // fmt.Printf("PC: %d\n", pc)
        // fmt.Printf("Instruction: %032b\n", instruction)
        // fmt.Printf("Operation: %d\n", op)
        switch op {
            case 0: if reg[c] != 0 { reg[a] = reg[b] }
            case 1: reg[a] = platters[reg[b]][reg[c]]
            case 2: platters[reg[a]][reg[b]] = reg[c]
            case 3: reg[a] = reg[b] + reg[c]
            case 4: reg[a] = reg[b] * reg[c]
            case 5: reg[a] = reg[b] / reg[c]
            case 6: reg[a] = ^(reg[b] & reg[c])
            case 7: return
            case 8: { platters = append(platters, make([]uint32, reg[c])); reg[b] = uint32(len(platters) - 1) }
            case 9: { platters[reg[c]] = nil }
            case 12: {
                if reg[b] != 0 { copy(platters[0], platters[reg[b]]) };
                pc = reg[c];
                continue
            }
            case 13: reg[(instruction >> 25) & 7] = instruction & 0x01FFFFFF
            default: panic(fmt.Errorf("Failed on %d", op)) 
        }
        pc++
    }
    
    fmt.Printf("%A\n", reg)
}

func read_platters(path string) []uint32 {
    b, err := ioutil.ReadFile(path)
    check(err)
    platters := make([]uint32, len(b) / 4)
    err = binary.Read(bytes.NewReader(b), binary.BigEndian, &platters)
    check(err)
    return platters
}

func check(err error) {
    if err != nil {
        panic(err)
    }    
}

func main() {
    platters := read_platters("out.um")
    fmt.Printf("Platters: %d\n", len(platters))
    run(platters)
}
