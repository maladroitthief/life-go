package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/maladroitthief/life-go/internal/life"
	"github.com/maladroitthief/life-go/internal/terminal"
)

var (
	fps      int
	duration int
	height   int
	width    int
	filename string
)

func main() {
	ctx := context.Background()

	termWriter := terminal.NewWriter(ctx)
	termWriter.Start()
	defer termWriter.Stop()

	flag.IntVar(&fps, "fps", 20, "frames per second")
	flag.IntVar(&duration, "duration", -1, "game of life duration")
	flag.IntVar(&height, "height", 64, "height of the game")
	flag.IntVar(&width, "width", 64, "width of the game")
	flag.StringVar(&filename, "file", "", "input file")

	flag.Parse()

	var game *life.Game
	if filename != "" {
		game = life.LoadRound(filename)
	} else {
		game = life.GenerateRound(height, width)
	}

	for i := 0; i != duration; i++ {
		game.Update()
		time.Sleep(time.Second / time.Duration(fps))
		fmt.Fprintf(termWriter, "%v\n", game.Render())
	}
}
