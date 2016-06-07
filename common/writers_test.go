package common

import (
	"bytes"
	"fmt"
	"testing"
)

type testSliceWriter struct {
	b [][]byte
}

func (tsw *testSliceWriter) Write(p []byte) (n int, err error) {
	l := make([]byte, len(p))
	copy(l, p)
	tsw.b = append(tsw.b, l)
	return len(p), nil
}

func TestLineWriter(t *testing.T) {
	tsw := &testSliceWriter{}
	lw := NewLineWriter(tsw)

	lineCount := 7
	lw.Write([]byte("0 line\n1 line\n2 line\n\n4 line"))
	lw.Write([]byte("+more\n5 line\n"))
	lw.Write([]byte("6 line"))

	lw.Flush()

	if len(tsw.b) != lineCount {
		t.Errorf("Expected %v individual rows; got %v", lineCount, len(tsw.b))
	}

	for x := 0; x < len(tsw.b); x++ {
		l := fmt.Sprintf("%v line\n", x)
		if x == 3 {
			if len(tsw.b[x]) != 1 {
				t.Errorf("Expected slice with only newline; got %v", tsw.b[x])
			}
			continue
		} else if x == 4 {
			l = "4 line+more\n"
		}
		if !bytes.Equal(tsw.b[x], []byte(l)) {
			t.Errorf("Expected slice %s equal to %s", []byte(l), tsw.b[x])
		}
	}
}

func testLastWritesWriter(writes, size int, t *testing.T) {
	lnw, _ := NewLastWritesWriter(size)
	for x := 0; x < writes; x++ {
		lnw.Write([]byte(fmt.Sprintf("%v line", x)))
	}

	r := lnw.Fetch()
	if len(r) != min(size, writes) {
		t.Errorf("Expected %v lines; got %v", min(size, writes), len(r))
	}

	expectedStart := writes - min(size, writes)
	for x := 0; x < len(r); x++ {
		if !bytes.Contains(r[x], []byte(fmt.Sprintf("%v line", expectedStart+x))) {
			t.Errorf("Expected %vth line; got '%s'", expectedStart+x, r[x])
		}
	}
}

func TestLastWritesWriter(t *testing.T) {
	testLastWritesWriter(0, 1, t)
	testLastWritesWriter(1, 4, t)
	testLastWritesWriter(4, 4, t)
	testLastWritesWriter(10, 4, t)
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}