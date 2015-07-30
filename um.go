package main

import (
    "io/ioutil"
    "fmt"
    "bytes"
    "encoding/binary"
)

func read_platters(path string) []uint32 {
    b, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    } 
    buf := bytes.NewReader(b)
    platters := make([]uint32, len(b) / 4)
    err = binary.Read(buf, binary.LittleEndian, &platters)
    if err != nil {
	    panic(err)
    }
    return platters
}

func main() {
    platters := read_platters("out.um")
    fmt.Printf("%d", len(platters))
}
