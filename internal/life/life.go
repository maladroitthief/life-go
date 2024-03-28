package life

import (
	"bytes"
	"math/rand"
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
	for columns := range cells {
		cells[columns] = make([]int, height)
	}

	game := &Game{
		cells:  cells,
		height: height,
		width:  width,
	}

	return game
}

func GenerateRound(height, width int) *Game {
	game := newGame(height, width)

	for i := 0; i < (height * width / 4); i++ {
		game.processCell(rand.Intn(width), rand.Intn(height), 1)
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

	for x := 0; x < g.width; x++ {
		for y := 0; y < g.height; y++ {
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
