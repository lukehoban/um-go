package main

import (
    "io/ioutil"
    "fmt"
    "bytes"
    "os"
    "flag"
    "encoding/binary"
)

func run(program []uint32) {
    reg := [8]uint32{0,0,0,0,0,0,0,0}
    platters := [][]uint32{program}
    freePlatters := []uint32{} 
    var pc uint32 = 0
    
    for ;; {
        instruction := platters[0][pc]
        op := (instruction >> 28) & 15
        a := ((instruction >> 6) & 7)
        b := ((instruction >> 3) & 7) 
        c := ((instruction >> 0) & 7)        
        switch op {
            case 0: if reg[c] != 0 { reg[a] = reg[b] }
            case 1: reg[a] = platters[reg[b]][reg[c]]
            case 2: platters[reg[a]][reg[b]] = reg[c]
            case 3: reg[a] = reg[b] + reg[c]
            case 4: reg[a] = reg[b] * reg[c]
            case 5: reg[a] = reg[b] / reg[c]
            case 6: reg[a] = ^(reg[b] & reg[c])
            case 7: return
            case 8: { 
                newPlatter := make([]uint32, reg[c])
                if len(freePlatters) > 0 {
                    platters[freePlatters[0]] = newPlatter
                    reg[b] = freePlatters[0]
                    freePlatters = freePlatters[1:]
                } else {
                    platters = append(platters, newPlatter); 
                    reg[b] = uint32(len(platters) - 1)
                } 
            }
            case 9: { platters[reg[c]] = nil; freePlatters = append(freePlatters, reg[c]) }
            case 10: os.Stdout.Write([]byte{byte(reg[c])})
            case 11: { b := []byte{0}; _, err := os.Stdin.Read(b); check(err); reg[c] = uint32(b[0]) }
            case 12: {
                if reg[b] != 0 { 
                    platters[0] = make([]uint32, len(platters[reg[b]]))
                    copy(platters[0], platters[reg[b]])
                }
                pc = reg[c]
                continue
            }
            case 13: reg[(instruction >> 25) & 7] = instruction & 0x01FFFFFF
            default: panic(fmt.Errorf("Failed on %d", op)) 
        }
        pc++
    }
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
    program := flag.String("program", "sandmark.umz", "The program to run on the Universal Machine.")
    flag.Parse()
    platters := read_platters(*program)
    run(platters)
}
