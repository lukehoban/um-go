package main

import (
    "io/ioutil"
    "fmt"
    "bytes"
    "errors"
    "encoding/binary"
)

func run(program []uint32) {
    reg := [8]uint32{0,0,0,0,0,0,0,0}
    platters := [][]uint32{program}
    pc := 0
    
    for ;; {
        instruction := platters[0][pc]
        op := (instruction >> 28) & 15
        a := ((instruction >> 6) & 7)
        b := ((instruction >> 3) & 7) 
        c := ((instruction >> 0) & 7)
        fmt.Printf("PC: %d\n", pc)
        fmt.Printf("Instruction: %0b\n", instruction)
        fmt.Printf("Operation: %d\n", op)
        switch op {
            case 0: if reg[c] != 0 { reg[a] = reg[b] }
            default: panic(errors.New("Failed")) 
        }
        fmt.Printf("%d\n", instruction)
    }
    
    fmt.Printf("%A\n", reg)
}

func read_platters(path string) []uint32 {
    b, err := ioutil.ReadFile(path)
    check(err)
    platters := make([]uint32, len(b) / 4)
    err = binary.Read(bytes.NewReader(b), binary.LittleEndian, &platters)
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
