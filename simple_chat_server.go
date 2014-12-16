package main

import(
    "fmt"
    "net"
    "container/list"
    "bytes"
)

type Client struct{
    Name string
    Incoming chan string
    Outgoing chan string
    Conn net.Conn
    Quit chan bool
    ClientList *list.List
}

func (c *Client) Read(buffer []byte) bool{
    bytesRead, err := c.Conn.Read(buffer)
    if err != nil{ //Network error problems
        c.Close()
        Log(err)
        return false
    }
    Log("Read, ", bytesRead, " bytes")
    return true
}

func (c *Client) Close(){
    c.Quit <- true
    c.Conn.Close()
    c.RemoveMe()
}

func (c *Client) Equal(other *Client) bool{
    if bytes.Equal([]byte(c.Name), []byte(other.Name)){
        if c.Conn == other.Conn{
            return true 
        } 
    }
    return false
}

func (c *Client) RemoveMe(){ //Returns void?
    for entry := c.ClientList.Front(); entry != nil; entry = entry.Next() {
        client := entry.Value.(Client)
        if c.Equal(&client){
            Log("Remove me: ", c.Name)
            c.ClientList.Remove(entry)
        }
    }
}

func Log(v ...interface{}){
    fmt.Println(v...)
}

func IOHandler(Incoming <- chan string, clientList *list.List){
    for{
        Log("IOHandler: Waiting for input")
        input := <-Incoming
        Log("IOHandling: ", input);
        for e := clientList.Front(); e != nil; e = e.Next(){
            client := e.Value.(Client)
            client.Incoming <-input
        }
    }
}

func ClientReader(client *Client){
    buffer := make([]byte, 2048)

    for client.Read(buffer){
        if bytes.Equal(buffer, []byte("/quit")){
            client.Close()
            break
        }
        Log("ClientReader received ", client.Name, "> ", string(buffer))
        send := client.Name + "> " + string(buffer)
        client.Outgoing <- send 
        for i := 0; i < 2048; i++ {
            buffer[i] = 0x00 
        }
    }

    client.Outgoing <- client.Name + " has left chat"
    Log("ClientReader stopped for ", client.Name)
}

func ClientSender(client *Client){
    for{
        select{
            case buffer := <-client.Incoming: 
                Log("Case buffer is sending: ", string(buffer), " to ", client.Name)
                count := 0
                for i := 0; i < len(buffer); i++{
                    if buffer[i] == 0x00{
                        break 
                    }
                    count++
                }
                Log("Send size: ", count)
                client.Conn.Write([]byte(buffer)[0: count])
            case <-client.Quit:
                Log("Quitting client: ", client.Name)
                client.Conn.Close()
                break
        }
    }
}

func ClientHandler(conn net.Conn, ch chan string, clientList *list.List){
    buffer := make([]byte, 1024)
    bytesRead, _ := conn.Read(buffer)
    name := string(buffer[0: bytesRead])
    newClient := &Client{name, make(chan string), ch, conn, make(chan bool), clientList}

    go ClientSender(newClient)
    go ClientReader(newClient)
    clientList.PushBack(*newClient)
    ch <-string(name + " has just joined the chat")
}

func main(){
    Log("hello, server!")

    clientList := list.New()
    in := make(chan string)
    go IOHandler(in, clientList)

    port := ":9988"
    tcpAddr, _ := net.ResolveTCPAddr("tcp", port)
    
    netListen, _ := net.Listen(tcpAddr.Network(), tcpAddr.String())
   
    defer netListen.Close()

    for{
        Log("Waiting for Client")
        connection, err := netListen.Accept()
        if err != nil{
            Log("Client error: ", err) 
        } else{
            go ClientHandler(connection, in, clientList) 
        }
    }
}
