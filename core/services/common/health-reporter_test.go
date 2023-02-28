package common

import (
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHealthReporter(t *testing.T) {
	hr := NewHealthReporter(time.Second)
	err1 := fmt.Errorf("error A")
	hr.Add(NewHealthError("test-system-1", err1))
	want := map[string]error{
		"test-system-1": errors.Join(err1),
	}
	got := hr.Report()
	t.Logf("got %+v want %+v", got, want)

	assertReportsEqual(t, want, got)
	assert.True(t, errors.Is(got["test-system-1"], err1))
	assertErrorsIs(t, got["test-system-1"], err1)

	hr.Add(NewHealthError("test-system-1", os.ErrExist))
	want = map[string]error{
		"test-system-1": errors.Join(err1, os.ErrExist),
	}
	got = hr.Report()
	t.Logf("got %+v want %+v", got, want)

	assertReportsEqual(t, want, got)
	assertErrorsIs(t, got["test-system-1"], err1, os.ErrExist)

	err3 := io.EOF

	hr.Add(NewHealthError("test-system-1", err3))
	want = map[string]error{
		"test-system-1": errors.Join(err1, os.ErrExist, err3),
	}
	got = hr.Report()
	t.Logf("got %+v want %+v", got, want)

	assertReportsEqual(t, want, got)
	assertErrorsIs(t, got["test-system-1"], err1, os.ErrExist)
}

func assertReportsEqual(t *testing.T, want, got map[string]error) {

	assert.Equal(t, len(want), len(got))
	for name, wantErr := range want {
		gotErr, exists := got[name]
		assert.True(t, exists)
		assert.Equal(t, wantErr, gotErr)
	}
}

func assertErrorsIs(t *testing.T, got error, wantAs ...error) {
	for i, wanted := range wantAs {
		t.Logf("got, want i %+v, %+v, %d", got, wanted, i)
		assert.ErrorIsf(t, got, wanted, "got error is not as wanted for %d (%s)", i, wanted)
	}
}
