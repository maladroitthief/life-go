package life

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

type (
	Game struct {
		cells  [][]int
		height int
		width  int
	}
)

func newGame(height, width int) *Game {
	cells := make([][]int, width)
	for rows := range cells {
		cells[rows] = make([]int, height)
	}

	game := &Game{
		cells:  cells,
		height: height,
		width:  width,
	}

	return game
}

func LoadRound(filename string) *Game {
	switch filepath.Ext(filename) {
	case ".txt":
		return LoadTxt(filename)
	default:
		return LoadTxt(filename)
	}
}

func GenerateRound(height, width int) *Game {
	game := newGame(height, width)

	for i := 0; i < (height * width / 4); i++ {
		game.processCell(rand.Intn(width), rand.Intn(height), 1)
	}

	return game
}

func LoadTxt(filename string) *Game {
	info, err := os.Stat(filename)
	if err != nil {
		fmt.Printf("%v file does not exist", filename)
		os.Exit(1)
	}

	if info.IsDir() {
		fmt.Printf("%v file is a directory", filename)
		os.Exit(2)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("%v file could not be read", filename)
		os.Exit(3)
	}

	split_data := bytes.Split(data, []byte("\n"))
	height := len(split_data)
	width := 0

	for i := range split_data {
		if len(split_data[i]) > width {
			width = len(split_data[i])
		}
	}

	game := newGame(height, width)
	for y := range split_data {
		for x := range len(split_data[y]) {
			char := split_data[y][x]
			switch {
			case char == '1':
				game.processCell(x, y, 1)
			default:
				game.processCell(x, y, 0)
			}
		}
	}

	return game
}

func (g *Game) Cells() [][]int {
	return g.cells
}

func (g *Game) processCell(x, y, status int) {
	x += g.width
	x %= g.width

	y += g.height
	y %= g.height

	g.cells[x][y] = status
}

func (g *Game) cellStatus(x, y int) int {
	x += g.width
	x %= g.width

	y += g.height
	y %= g.height

	return g.cells[x][y]
}

func (g *Game) LivingNeighbors(x, y int) int {
	count := 0
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}

			if g.cellStatus(x+i, y+j) > 0 {
				count++
			}
		}
	}

	return count
}

func (g *Game) NextStatus(x, y int) int {
	livingNeighbors := g.LivingNeighbors(x, y)
	isAlive := g.cellStatus(x, y) > 0

	if livingNeighbors == 3 {
		return 1
	}

	if livingNeighbors == 2 && isAlive {
		return 1
	}

	return 0
}

func (g *Game) Update() {
	nextRound := newGame(g.height, g.width)
	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
			nextRound.processCell(x, y, g.NextStatus(x, y))
		}
	}

	g.cells = nextRound.cells
}

func (g *Game) Render() string {
	var buffer bytes.Buffer

	for y := 0; y < g.height; y++ {
		for x := 0; x < g.width; x++ {
			if g.cellStatus(x, y) > 0 {
				buffer.WriteString("*")
			} else {
				buffer.WriteByte(byte(' '))
			}
		}
		buffer.WriteByte('\n')
	}

	return buffer.String()
}
