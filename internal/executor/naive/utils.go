// exec implementation in Go is awful! It's full of race conditions (on case we specified own buffer)
// So I created wrapper for bytes.Buffer which is mt-safe

package naive

import (
	"bytes"
	"sync"
)

type TSafeBuf struct {
	buf bytes.Buffer
	mut *sync.Mutex
}

func (tb *TSafeBuf) Write(b []byte) (int, error) {
	tb.mut.Lock()
	defer tb.mut.Unlock()

	return tb.buf.Write(b)
}

func (tb *TSafeBuf) Read(p []byte) (n int, err error) {
	tb.mut.Lock()
	defer tb.mut.Unlock()

	return tb.buf.Read(p)
}

func (tb *TSafeBuf) Reset() {
	tb.mut.Lock()
	defer tb.mut.Unlock()

	tb.buf.Reset()
}

func GetTSafeBuf() *TSafeBuf {
	return &TSafeBuf{
		mut: new(sync.Mutex),
	}
}
