package main

import (
	"fmt"
	"net"
	"strings"
)
func handleConnection(conn net.Conn, aof *Aof) {
    defer conn.Close()

    for {
        resp := NewResp(conn)
        value, err := resp.Read()
        if err != nil {
            fmt.Println(err)
            return
        }
        
        if value.typ != "array" {
            fmt.Println("Invalid request")
            continue
        }

        if len(value.array) == 0 {
            fmt.Println("Expected length is greater than 0")
            continue
        }

        command := strings.ToUpper(value.array[0].bulk)
        args := value.array[1:]

        writer := NewWriter(conn)
        handler, ok := Handlers[command]
        if !ok {
            fmt.Println("Invalid command: ", command)
            writer.Write(Value{typ: "string", str: ""})
            continue
        }
        
        if command == "SET" || command == "HSET" {
            aof.write(value)  
        }

        result := handler(args)
        writer.Write(result)
    }
}

func main() {

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := NewAof("db.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	for {
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	go handleConnection(conn, aof)

	}
}
