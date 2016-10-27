package iologger

import (
	"bytes"
	"io"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleNewReadWriteLogger() {
	var dataSource io.ReadWriter

	// the below will emit seperate log entries for reads and writes
	// with the hex of the data
	rw := NewReadWriteLogger(dataSource, func(p []byte) {
		log.Printf("debug: read: [%x]\n", p)
	}, func(p []byte) {
		log.Printf("debug: write: [%x]\n", p)
	})

	// use dataSource instead of conn directly
	rw.Write([]byte("foo"))
}

func TestReadLoggerNew(t *testing.T) {
	r := NewReadLogger(bytes.NewReader(nil), func(p []byte) {})
	if rl, ok := r.(*readLogger); ok {
		assert.NotNil(t, rl.r)
		assert.NotNil(t, rl.fn)
	} else {
		t.Fatal("reader is not readLogger")
	}
}

func TestReadLoggerRead(t *testing.T) {
	var invoked bool

	b := bytes.NewReader([]byte("foo"))
	r := &readLogger{
		r: b,
		fn: func(p []byte) {
			invoked = true
			assert.Equal(t, []byte("foo"), p)
		},
	}

	d := make([]byte, 6)

	n, err := r.Read(d)
	if assert.NoError(t, err) {
		assert.Equal(t, 3, n)
		assert.Equal(t, []byte("foo\x00\x00\x00"), d)
	}

	assert.True(t, invoked, "LoggerFunc not invoked")
}

func TestWriteLoggerNew(t *testing.T) {
	w := NewWriteLogger(bytes.NewBuffer(nil), func(p []byte) {})
	if wl, ok := w.(*writeLogger); ok {
		assert.NotNil(t, wl.w)
		assert.NotNil(t, wl.fn)
	} else {
		t.Fatal("writer is not writeLogger")
	}
}

func TestWriteLoggerWrite(t *testing.T) {
	var invoked bool

	b := bytes.NewBuffer(nil)
	w := &writeLogger{
		w: b,
		fn: func(p []byte) {
			invoked = true
			assert.Equal(t, []byte("foo"), p)
		},
	}

	n, err := w.Write([]byte("foo"))
	if assert.NoError(t, err) {
		assert.Equal(t, 3, n)
	}

	assert.True(t, invoked, "LoggerFunc not invoked")
}

func TestNewReadWriteLogger(t *testing.T) {
	rw := NewReadWriteLogger(bytes.NewBuffer(nil), nil, nil)
	if rwl, ok := rw.(*readWriteLogger); ok {
		assert.NotNil(t, rwl.Reader)
		assert.NotNil(t, rwl.Writer)
	} else {
		t.Fatal("readWriter is not readWriteLogger")
	}

	rw = NewReadWriteLogger(bytes.NewBuffer(nil), func(p []byte) {}, nil)
	if rwl, ok := rw.(*readWriteLogger); ok {
		if assert.NotNil(t, rwl.Reader) {
			if rl, ok := rwl.Reader.(*readLogger); ok {
				assert.NotNil(t, rl.fn)
			} else {
				t.Fatal("reader is not readLogger")
			}
		}
	} else {
		t.Fatal("readWriter is not readWriteLogger")
	}

	rw = NewReadWriteLogger(bytes.NewBuffer(nil), nil, func(p []byte) {})
	if rwl, ok := rw.(*readWriteLogger); ok {
		if assert.NotNil(t, rwl.Writer) {
			if wl, ok := rwl.Writer.(*writeLogger); ok {
				assert.NotNil(t, wl.fn)
			} else {
				t.Fatal("writer is not writeLogger")
			}
		}
	} else {
		t.Fatal("readWriter is not readWriteLogger")
	}
}
