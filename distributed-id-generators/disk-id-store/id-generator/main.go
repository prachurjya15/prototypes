// This is a separate service that runs and other services will use this
package main

import (
	"io"
	"log"
	"net"
	"net/rpc"
	"os"
	"strconv"
	"sync"
)

var Buffer int = 100 // Every Buffer increase of val, we will write to file

type IdGenerator struct {
	currId int
	mu     sync.Mutex
	file   *os.File
}

type IdGenReq struct{}

type IdGenResp struct {
	Id int
}

func (gen *IdGenerator) GetNextId(req *IdGenReq, resp *IdGenResp) error {

	gen.mu.Lock()
	gen.currId = gen.currId + 1
	resp.Id = gen.currId
	// DONT WRITE TILL BUFFER
	if gen.currId%Buffer == 0 {
		// We are doing this rather than WriteFile since we already have the file opened in the main call
		gen.file.Seek(0, 0)
		gen.file.Truncate(0)
		gen.file.WriteString(strconv.Itoa(gen.currId))
		gen.file.Sync()
	}
	gen.mu.Unlock()
	return nil

}

func main() {
	file, err := os.OpenFile("id-gen.txt", os.O_CREATE|os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		log.Fatalf("Error in opening the src of truth")
	}

	b := make([]byte, 64)
	numBytes, err := file.Read(b)
	if err != nil && err != io.EOF {
		log.Fatalf("Error in reading the src of truth. Error:[%s]", err)
	}
	var currId int

	if numBytes == 0 {
		currId = 0
	} else {
		s := string(b[:numBytes])
		currId, err = strconv.Atoi(s)
		currId = currId + Buffer + 1 //Safest Value
		if err != nil {
			log.Fatalf("Can't convert string: [%s] to int", s)
		}

	}

	gen := &IdGenerator{currId: currId, file: file}
	rpc.Register(gen)

	l, err := net.Listen("tcp", ":8085")
	if err != nil {
		log.Fatalf("Error in listening to tcp addr: %s", "8085")
	}
	log.Println("Id generator service is running in port :8085")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error", err)
			continue
		}
		go rpc.ServeConn(conn)
	}

}
