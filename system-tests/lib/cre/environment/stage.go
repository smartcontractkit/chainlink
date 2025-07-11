package environment

import (
	"fmt"
	"sync"
	"time"
)

type StageGen struct {
	mu        sync.Mutex
	no        int
	total     int
	startTime time.Time
	label     string
}

func NewStageGen(total int, label string) *StageGen {
	return &StageGen{
		no:        1,
		total:     total,
		startTime: time.Now(),
		label:     label,
	}
}

func (s *StageGen) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("[%s %d/%d]", s.label, s.no, s.total)
}

func (s *StageGen) Next() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.no++
	s.startTime = time.Now()
}

func (s *StageGen) Wrap(msg string, args ...any) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return fmt.Sprintf("[%s %d/%d] %s\n\n", s.label, s.no, s.total, fmt.Sprintf(msg, args...))
}

func (s *StageGen) WrapAndNext(msg string, args ...any) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	formatted := fmt.Sprintf("\n[%s %d/%d] %s\n\n", s.label, s.no, s.total, fmt.Sprintf(msg, args...))
	s.no++
	s.startTime = time.Now()
	return formatted
}

func (s *StageGen) Elapsed() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return time.Since(s.startTime)
}

func (s *StageGen) ElapsedAndNext() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	duration := time.Since(s.startTime)
	s.no++
	s.startTime = time.Now()
	return duration
}
