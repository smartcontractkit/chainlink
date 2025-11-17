package stagegen

import (
	"strings"
	"testing"
	"time"
)

func TestNewStageGen(t *testing.T) {
	s := NewStageGen(5, "Test")
	if s.no != 1 || s.total != 5 || s.label != "Test" {
		t.Errorf("unexpected StageGen values: %+v", s)
	}
}

func TestString(t *testing.T) {
	s := NewStageGen(3, "Stage")
	str := s.String()
	if !strings.Contains(str, "[Stage 1/3]") {
		t.Errorf("unexpected String output: %s", str)
	}
}

func TestNext(t *testing.T) {
	s := NewStageGen(2, "Next")
	s.Next()
	if s.no != 2 {
		t.Errorf("Next() did not increment stage: got %d, want 2", s.no)
	}
}

func TestWrap(t *testing.T) {
	s := NewStageGen(4, "Wrap")
	msg := s.Wrap("hello %s", "world")
	if !strings.Contains(msg, "[Wrap 1/4] hello world") {
		t.Errorf("Wrap() output incorrect: %s", msg)
	}
}

func TestWrapAndNext(t *testing.T) {
	s := NewStageGen(2, "WrapNext")
	msg := s.WrapAndNext("msg %d", 42)
	if !strings.Contains(msg, "[WrapNext 1/2] msg 42") {
		t.Errorf("WrapAndNext() output incorrect: %s", msg)
	}
	if s.no != 2 {
		t.Errorf("WrapAndNext() did not increment stage: got %d, want 2", s.no)
	}
}

func TestElapsedAndNext(t *testing.T) {
	s := NewStageGen(2, "Elapsed")
	time.Sleep(10 * time.Millisecond)
	d := s.ElapsedAndNext()
	if d < 10*time.Millisecond {
		t.Errorf("ElapsedAndNext() duration too short: %v", d)
	}
	if s.no != 2 {
		t.Errorf("ElapsedAndNext() did not increment stage: got %d, want 2", s.no)
	}
}

func TestElapsed(t *testing.T) {
	s := NewStageGen(2, "Elapsed")
	time.Sleep(5 * time.Millisecond)
	d := s.Elapsed()
	if d < 5*time.Millisecond {
		t.Errorf("Elapsed() duration too short: %v", d)
	}
}
