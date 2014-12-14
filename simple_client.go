package main

import(
    "net"
    "fmt"
)

func main() {
    conn, _ := net.Dial("tcp", ":8080")
    fmt.Fprintf(conn, "newmessage\n") 
}
