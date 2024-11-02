package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd *bufio.Reader
	mu sync.Mutex
}

func NewAof(path string) (*Aof , error){
	f,err:=os.OpenFile(path, os.O_CREATE|os.O_RDWR,0666)
	if err != nil{
		return nil, err
	}
	aof := &Aof{
		file: f,
		rd: bufio.NewReader(f),
	}

	//goroutine to execute all the command in the disk after recovery
	go func() {
        aof.mu.Lock()
        defer aof.mu.Unlock()
        
        for {
            value, err := aof.readCommand()
            if err != nil {
                if err == io.EOF {
                    break
                }
                fmt.Println("Error reading AOF file:", err)
                break
            }
            
            if len(value.array) > 0 {
                command := strings.ToUpper(value.array[0].bulk)
                args := value.array[1:]
                
                handler, ok := Handlers[command]
                if ok {
                    handler(args)  
                }
            }
        }
    }()

	//goroutine to sync AOF to disk every 1 sec
	go func() {
		for{
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Close() error{
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

func (aof *Aof) write(value Value) error{
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_,err := aof.file.Write(value.Marshal())

	if err !=nil{
		return err
	}

	return nil
}

func (aof *Aof) readCommand() (Value, error) {
    resp := NewResp(aof.rd)
    return resp.Read()
}