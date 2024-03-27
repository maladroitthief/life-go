package terminal

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type (
	Writer struct {
		ctx    context.Context
		cancel context.CancelFunc

		Writer io.Writer
		buffer bytes.Buffer
		mutex  sync.Mutex
		once   sync.Once
		lines  int
		active bool
		kill   bool
		ticker *time.Ticker
	}
)

const (
	screenRefreshDuration = 20
)

func NewWriter(ctx context.Context) *Writer {
	writer := &Writer{
		Writer: os.Stdout,
		ticker: time.NewTicker(screenRefreshDuration),
	}
	writer.ctx, writer.cancel = context.WithCancel(ctx)

	return writer
}

func (w *Writer) Start() {
	w.kill = false
	w.ticker.Reset(screenRefreshDuration)

	go w.Listen()
}

func (w *Writer) Stop() {
	w.once.Do(
		func() {
			w.ticker.Stop()
			w.Flush()
			w.kill = true
		},
	)
}

func (w *Writer) Listen() {
	if w.active {
		return
	}

	w.mutex.Lock()
	w.active = true
	w.mutex.Unlock()

ListenerLoop:
	for !w.kill {
		select {
		case <-w.ctx.Done():
			break ListenerLoop
		case <-w.ticker.C:
			w.Wait()
		}
	}

	w.active = false
	return
}

func (w *Writer) Wait() {
	w.Flush()
}

func (w *Writer) Write(input []byte) (length int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	return w.buffer.Write(input)
}

func (w *Writer) Flush() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.buffer.Len() == 0 {
		return nil
	}

	w.clear()

	lines := 0
	for _, byte := range w.buffer.Bytes() {
		if byte == '\n' {
			lines++
		}
	}

	w.lines = lines
	_, err := w.Writer.Write(w.buffer.Bytes())
	w.buffer.Reset()

	return err
}

func (w *Writer) clear() {
	for i := 0; i < w.lines; i++ {
		fmt.Fprint(w.Writer, "\033[0A")
		fmt.Fprint(w.Writer, "\033[2K\r")
	}
}
