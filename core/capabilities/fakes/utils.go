package fakes

import (
	"fmt"
)

type SimpleStats struct {
	Stats map[string]int64
}

func NewSimpleStats() *SimpleStats {
	return &SimpleStats{Stats: make(map[string]int64)}
}

func (s *SimpleStats) UpdateMaxStat(key string, value int64) {
	if value >= s.Stats[key] {
		s.Stats[key] = value
	}
}

func (s *SimpleStats) PrintToStdout(title string) {
	fmt.Println(title)
	for k, v := range s.Stats {
		fmt.Printf("  %s: %d\n", k, v)
	}
}
