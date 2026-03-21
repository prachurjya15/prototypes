package services

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type IdGenReq struct{}

type IdGenResp struct {
	Id int
}

type Service struct {
	Name      string
	rpcClient *rpc.Client
}

func NewService(name string, connStr string) *Service {
	rc, err := rpc.Dial("tcp", connStr)
	if err != nil {
		log.Fatalf("Error in dailing the rpc conn. Error: %s\n", err)
	}
	s := Service{
		Name:      name,
		rpcClient: rc,
	}
	return &s
}

func (s *Service) doWork() {
	// Mimic some DB entry or something
	// Call Id_Generator and get Id and insert into db
	idGenResp := &IdGenResp{}
	t := time.Now()
	err := s.rpcClient.Call("IdGenerator.GetNextId", &IdGenReq{}, idGenResp)
	if err != nil {
		log.Printf("Error in calling rpc function. Error : %s\n", err)
	}
	tt := time.Since(t).Milliseconds()
	log.Printf("[%s] Time taken for Id Generation is: [%d] \n", s.Name, tt)
	stat := fmt.Sprintf("INSERT INTO DBX (id) VALUES(%d)", idGenResp.Id)
	log.Printf("%s called DB with insert query: [%s] \n", s.Name, stat)
	// MIMIC DB CALL
	time.Sleep(1 * time.Millisecond)
}

func (s *Service) Work() {
	//At various frequencies call doWork
	for i := range 10 {
		if i%3 == 0 {
			s.doWork()
			time.Sleep(1000 * time.Millisecond)
		}
	}
}
