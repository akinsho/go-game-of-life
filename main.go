// Provides a Game of Life
package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten"
)

//World represents the game world
type World struct {
	area [][]bool
	rnd  *rand.Rand
}

// NewWorld Creates a new world
func NewWorld(width, height int) *World {
	world := World{
		area: makeArea(width, height),
		rnd:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return &world
}

// RandomSeed inits world with a random state
func (w *World) RandomSeed(limit int) {
	height := len(w.area)
	width := len(w.area[0])
	for i := 0; i < limit; i++ {
		x := w.rnd.Intn(width)
		y := w.rnd.Intn(height)
		w.area[y][x] = true
	}
}

//Progress advances the game state by checking the current position
// condition of the cells against the rules
func (w *World) Progress() {
	height := len(w.area)
	width := len(w.area[0])
	next := makeArea(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pop := neighbourCount(w.area, x, y)
			switch {
			case pop < 2:
				//Rule 1. Any live cell with fewerthan 2 neibouers dies, as if caused by under population
				next[y][x] = false
			case (pop == 2 || pop == 3) && w.area[y][x]:
				//Rule 2. Any live cell with two or three live neighbours lives on to the next generation
				next[y][x] = true
			case pop > 3:
				// Rule 3. Any live cell with more than three live neighbours dies, as if bu over population
				next[y][x] = false
			case pop == 3:
				//Rule 4. any dead cell with exactly three live neighbours becomes a live cell, as if by reproduction
				next[y][x] = true
			}
		}
	}
	w.area = next
}

//DrawImage Draws the pixels which make up the world
func (w *World) DrawImage(pix []uint8) {
	height := len(w.area)
	width := len(w.area[0])
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pos := 4*y*width + 4*x
			if w.area[y][x] {
				pix[pos] = 0xff
				pix[pos+1] = 0xff
				pix[pos+2] = 0xff
				pix[pos+3] = 0xff
			} else {
				pix[pos] = 0
				pix[pos+1] = 0
				pix[pos+2] = 0
				pix[pos+3] = 0
			}
		}
	}
}

func neighbourCount(a [][]bool, x, y int) int {
	height := len(a)
	width := len(a[0])
	lowX := 0
	if x > 0 {
		lowX = x - 1
	}
	lowY := 0
	if y > 0 {
		lowY = y - 1
	}
	highX := width - 1
	if x < width-1 {
		highX = x + 1
	}
	highY := height - 1
	if y < height-1 {
		highY = y + 1
	}
	near := 0
	for pY := lowY; pY <= highY; pY++ {
		for pX := lowX; pX <= highX; pX++ {
			if !(pX == x && pY == y) && a[pY][pX] {
				near++
			}
		}
	}
	return near
}

func makeArea(width, height int) [][]bool {
	area := make([][]bool, height)
	for i := 0; i < height; i++ {
		area[i] = make([]bool, width)
	}
	return area
}

const (
	screenWidth  = 320
	screenHeight = 240
)

var (
	world  = NewWorld(screenWidth, screenHeight)
	pixels = make([]uint8, screenWidth*screenHeight*4)
)

func update(screen *ebiten.Image) error {
	world.Progress()
	if ebiten.IsRunningSlowly() {
		return nil
	}
	world.DrawImage(pixels)
	screen.ReplacePixels(pixels)
	return nil
}

func main() {
	world.RandomSeed(int((screenWidth * screenHeight) / 10))
	if err := ebiten.Run(update, screenWidth, screenHeight, 2.0, "Game of Life"); err != nil {
		log.Fatal(err)
	}
}
