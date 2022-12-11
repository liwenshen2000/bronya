package pkg

import (
	"io"
	"sync"
)

func Relay(inbound, outbound io.ReadWriter) {
	// ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup

	forward := func(dst io.Writer, src io.Reader) {
		_, e := io.Copy(dst, src)

		if e != nil {
			if conn, ok := src.(interface{ CloseRead() error }); ok {
				conn.CloseRead()
			}
		}
		if conn, ok := dst.(interface{ CloseWrite() error }); ok {
			conn.CloseWrite()
		}

		wg.Done()
	}

	wg.Add(2)
	go forward(inbound, outbound)
	go forward(outbound, inbound)
	wg.Wait()

	if c, ok := inbound.(io.Closer); ok {
		defer c.Close()
	}
	if c, ok := outbound.(io.Closer); ok {
		c.Close()
	}
}
