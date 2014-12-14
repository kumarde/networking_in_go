package main

import(
    "net"
    "fmt"
)

func main() {
    ln, _ := net.Listen("tcp", ":8080")
    
    for{
        conn, _ := ln.Accept()
        var cmd []byte
        fmt.Fscan(conn, &cmd)
        fmt.Println("Message: ", string(cmd))
    }
}
