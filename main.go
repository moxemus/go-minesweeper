package main

import (
	"fmt"
	"math/rand"
	"strconv"
)

const height = 9
const width = 9
const bombs_count = 5

type Cell struct {
	Value   int
	Visible bool
}

type GameMap struct {
	Height    int
	Width     int
	BombCount int
	Grid      [][]Cell
}

func (g *GameMap) Init(h, w, bombs int) {
	g.Height = h
	g.Width = w
	g.BombCount = bombs
	g.Grid = make([][]Cell, g.Height)

	// Init map
	for i := 0; i < g.Height; i++ {
		g.Grid[i] = make([]Cell, g.Width)

		for j := 0; j < g.Width; j++ {
			g.Grid[i][j].Value = 0
			g.Grid[i][j].Visible = false
		}
	}

	// Add bombs
	for i := 0; i < g.BombCount; i++ {
		y := rand.Intn(g.Height)
		x := rand.Intn(g.Width)

		g.Grid[x][y].Value = -1

		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				ny, nx := y+dy, x+dx

				if ny < 0 || ny >= g.Height || nx < 0 || nx >= g.Width || g.Grid[nx][ny].Value == -1 {
					continue
				}

				g.Grid[nx][ny].Value++
			}
		}
	}
}

func (g *GameMap) openZeros(x, y int) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			nx := x + dx
			ny := y + dy

			if nx < 0 || ny < 0 || nx >= width || ny >= g.Height {
				continue
			}

			if g.Grid[nx][ny].Visible {
				continue
			}

			g.Grid[nx][ny].Visible = true

			if g.Grid[nx][ny].Value == 0 {
				g.openZeros(nx, ny)
			}
		}
	}
}

func (g *GameMap) checkWin() bool {
	for i := 0; i < g.Height; i++ {
		for j := 0; j < g.Width; j++ {
			if g.Grid[j][i].Value != -1 && !g.Grid[j][i].Visible {
				return false
			}
		}
	}

	return true
}

func (g *GameMap) pressCell(x, y int) bool {
	g.Grid[x][y].Visible = true

	return g.Grid[x][y].Value == -1
}

// Screen Handler
type Drawer interface {
	drawMap()
}

type TerminalDrawer struct {
	drawer Drawer
}

func Init(d Drawer) *TerminalDrawer {
	return &TerminalDrawer{
		drawer: d,
	}
}

func (s *TerminalDrawer) write(text string) {
	fmt.Print(text)
}

func (s *TerminalDrawer) drawMap(cells [][]Cell, height, width int) {
	// Column coordinates
	fmt.Print("    ")
	for i := 0; i < width; i++ {
		fmt.Printf("%3d ", i+1)
	}
	fmt.Println()

	// Hortizontal border
	fmt.Print("   ")
	for i := 0; i < width; i++ {
		fmt.Print("----")
	}
	fmt.Println()

	for i := 0; i < height; i++ {
		// Vertical border
		fmt.Printf("%2d |", i+1)

		for j := 0; j < width; j++ {
			if !cells[j][i].Visible {
				fmt.Printf("%3c ", '#')
			} else {
				fmt.Printf("%3d ", cells[j][i].Value)
			}
		}

		fmt.Println()
	}
}

func (s *TerminalDrawer) clear() {
	fmt.Println("\033[H\033[2J")
}

func main() {
	gameMap := &GameMap{}
	gameMap.Init(height, width, bombs_count)

	drawHandler := &TerminalDrawer{}

	var input string

	for {
		drawHandler.drawMap(gameMap.Grid, gameMap.Height, gameMap.Width)

		_, err := fmt.Scanln(&input)
		drawHandler.clear()

		// Handle user input
		if len(input) == 3 && input[0] >= '1' && input[0] <= '9' && input[2] >= '1' && input[2] <= '9' && err == nil {
			x := int(input[0]-'0') - 1
			y := int(input[2]-'0') - 1

			// Handle cell click and check is user was exploded
			if gameMap.pressCell(x, y) {
				drawHandler.write("Game over. Game restarted.")
				gameMap.Init(height, width, bombs_count)
			} else {
				drawHandler.write("Your choise is: " + strconv.Itoa(x+1) + strconv.Itoa(y+1) + " result - 0")
				gameMap.openZeros(x, y)
			}

			if gameMap.checkWin() {
				drawHandler.write("Victory! Game restarted.")
				gameMap.Init(height, width, bombs_count)
			}

		} else if input == "q" {
			break
		} else if input == "r" {
			gameMap.Init(height, width, bombs_count)
		} else {
			drawHandler.write("Invalid input. Please select cell coordinates in format X,Y")
		}
	}
}
