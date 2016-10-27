// iologger implements io package compatible shims for logging and debugging
// data inbetween readers and writers.
package iologger

import (
	"io"
)

// LoggerFunc is a function called AFTER every Read/Write operation with a
// byte slice of the data representing what was sent by the underlying
// reader/writer.
type LoggerFunc func([]byte)

type readLogger struct {
	r  io.Reader
	fn LoggerFunc
}

// NewReadLogger() returns a new io.Reader compatible reader configured
// to invoke the passed LoggerFunc AFTER every read.
func NewReadLogger(r io.Reader, fn LoggerFunc) io.Reader {
	return &readLogger{
		r:  r,
		fn: fn,
	}
}

func (r *readLogger) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)

	r.fn(p[:n])

	return n, err
}

type writeLogger struct {
	w  io.Writer
	fn LoggerFunc
}

// NewWriteLogger() returns a new io.Writer compatible writer configured
// to invoke the passed LoggerFunc AFTER every write.
func NewWriteLogger(w io.Writer, fn LoggerFunc) io.Writer {
	return &writeLogger{
		w:  w,
		fn: fn,
	}
}

func (w *writeLogger) Write(p []byte) (int, error) {
	n, err := w.w.Write(p)

	w.fn(p[:n])

	return n, err
}

type readWriteLogger struct {
	io.Reader
	io.Writer
}

// NewReadWriteLogger() returns a new io.ReadWriter configured to invoke
// the passed LoggerFuncs AFTER every read/write. To only configure a single
// LoggerFunc, set the other to nil.
func NewReadWriteLogger(rw io.ReadWriter, r LoggerFunc, w LoggerFunc) io.ReadWriter {
	n := &readWriteLogger{
		Reader: rw,
		Writer: rw,
	}

	if r != nil {
		n.Reader = NewReadLogger(rw, r)
	}
	if w != nil {
		n.Writer = NewWriteLogger(rw, w)
	}

	return n
}
