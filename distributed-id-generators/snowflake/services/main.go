package services

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Service struct {
	Name   string
	Id     int16
	SeqNum int
	mu     sync.Mutex
	lastMs int64
}

var customEpoch int64 = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

func GetServiceId() int16 {
	// Ideally read from zookeeper/etcd
	return int16(rand.Intn(1024))
}

func (s *Service) GetSnowFlakeId() int {
	// 64bit Id -->
	nowTime := time.Now().UnixMilli()
	epochTime := (nowTime - customEpoch) & 0x1FFFFFFFFFF
	machineId := s.Id & 0x3FF // 0011 1111 1111
	s.mu.Lock()
	if s.lastMs == nowTime { // If I am asking for a new Id within same ms, then only update the seq number, else reset it back to 0
		s.SeqNum = s.SeqNum + 1
	} else {
		s.SeqNum = 0
		s.lastMs = nowTime
	}
	seqNum := s.SeqNum & 0xFFF // 1111 1111 1111
	s.mu.Unlock()

	return int(epochTime<<22 | int64(machineId)<<12 | int64(seqNum))
}

func NewService(name string) *Service {
	return &Service{
		Name: name,
		Id:   GetServiceId(),
	}
}

func (s *Service) Work() {
	var wg sync.WaitGroup
	for range 10 {
		wg.Go(func() {
			sfId := s.GetSnowFlakeId()
			stmt := fmt.Sprintf("[%s] INSERT INTO DBX id VALUES(%d)", s.Name, sfId)
			log.Println(stmt)
		})
	}
	wg.Wait()

}
