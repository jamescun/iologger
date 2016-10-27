iologger
========

iologger implements io package compatible shims for logging and debugging data inbetween readers and writers.


[![GoDoc](https://godoc.org/github.com/jamescun/iologger?status.png)](https://godoc.org/github.com/jamescun/iologger)


Example
-------

```go
func handleConn(conn net.Conn) {
	defer conn.Close()

	rw := iologger.NewReadWriteLogger(conn, func(p []byte) {
		log.Printf("debug: read: [%x]\n", p)
	}, func(p []byte) {
		log.Printf("debug: write: [%x]\n", p)
	})

	data := make([]byte, 3)
	rw.Read(data)
	// debug: read: [aabbcc]

	rw.Write([]byte{0xDD, 0xEE, 0xFF})
	// debug: write: [ddeeff]
}
```
