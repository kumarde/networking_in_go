package main

import(
    "net"
    "fmt"
)

func main() {
    conn, _ := net.Dial("tcp", ":9988")
    fmt.Fprintf(conn, "debug12") 
    for{}
}
