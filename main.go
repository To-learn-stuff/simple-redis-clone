package main

import (
	"fmt"
	"net"
)

func main(){

	l, err := net.Listen("tcp",":6379")
	if err != nil {
		fmt.Println(err)
		return 
	}

	conn, err := l.Accept()
	if err != nil{
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for{
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			fmt.Println(err)
			return
		}
		_=value

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})


	}
}