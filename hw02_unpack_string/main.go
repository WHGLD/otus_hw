package main

import (
    "flag"
    "fmt"

    unpacking"./hw02unpackstring"
)

func main() {
    var input string

    flag.StringVar(&input, "input", "", "Unpack the string")
    flag.Parse()

    result, error := unpacking.Unpack(input)

    fmt.Println(result, error)
}
