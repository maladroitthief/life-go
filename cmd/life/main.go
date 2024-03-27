package main

import (
	"context"
	"fmt"
	"time"

	"github.com/maladroitthief/life-go/internal/terminal"
)

func main() {
	ctx := context.Background()

	termWriter := terminal.NewWriter(ctx)
	termWriter.Start()
	defer termWriter.Stop()

	for i := range 10 {
		termWriter.Write([]byte(fmt.Sprintf("%v\n", i)))
		time.Sleep(1 * time.Second)
	}
}
