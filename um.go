package main

import (
    "io/ioutil"
    "fmt"
    "bytes"
    "encoding/binary"
)

func check(err error) {
    if err != nil {
        panic(err)
    }    
}

func read_platters(path string) []uint32 {
    b, err := ioutil.ReadFile(path)
    check(err)
    platters := make([]uint32, len(b) / 4)
    err = binary.Read(bytes.NewReader(b), binary.LittleEndian, &platters)
    check(err)
    return platters
}

func main() {
    platters := read_platters("out.um")
    fmt.Printf("%d", len(platters))
}
