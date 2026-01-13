package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
)

const height = 9
const width = 9
const bombs_count = 5

type Point struct {
	X int
	Y int
}

type Cell struct {
	Value   int
	Visible bool
}

type GameMap struct {
	height    int
	width     int
	bombCount int
	grid      [][]Cell
}

func (g *GameMap) Init(h, w, bombs int) {
	g.height = h
	g.width = w
	g.bombCount = bombs
	g.grid = make([][]Cell, g.height)

	// Init map
	for i := 0; i < g.height; i++ {
		g.grid[i] = make([]Cell, g.width)

		for j := 0; j < g.width; j++ {
			g.grid[i][j].Value = 0
			g.grid[i][j].Visible = false
		}
	}

	// Add bombs
	for i := 0; i < g.bombCount; i++ {
		var x, y int

		// Find an empty place for bomb
		for {
			y = rand.Intn(g.height)
			x = rand.Intn(g.width)

			if g.grid[x][y].Value != -1 {
				break
			}
		}

		g.grid[x][y].Value = -1

		// Update numbers near bomb
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				ny, nx := y+dy, x+dx

				if ny >= 0 && ny < g.height && nx >= 0 && nx < g.width && g.grid[nx][ny].Value != -1 {
					g.grid[nx][ny].Value++
				}
			}
		}
	}
}

func (g *GameMap) restart() {
	g.Init(g.height, g.width, g.bombCount)
}

func (g *GameMap) openZeros(start Point) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			x := start.X + dx
			y := start.Y + dy

			if x < 0 || y < 0 || x >= width || y >= g.height || g.grid[x][y].Visible || g.grid[x][y].Value == -1 {
				continue
			}

			g.grid[x][y].Visible = true

			if g.grid[x][y].Value == 0 {
				g.openZeros(Point{X: x, Y: y})
			}
		}
	}
}

func (g *GameMap) checkWin() bool {
	for i := 0; i < g.height; i++ {
		for j := 0; j < g.width; j++ {
			if g.grid[j][i].Value != -1 && !g.grid[j][i].Visible {
				return false
			}
		}
	}

	return true
}

func (g *GameMap) openCell(point Point) bool {
	g.grid[point.X][point.Y].Visible = true

	return g.grid[point.X][point.Y].Value == -1
}

// Drawer logic
type Drawer interface {
	drawMap(cells [][]Cell, height, width int)
	write(s string)
	clearScreen()
}

type TerminalDrawer struct {
	drawer Drawer
}

func (s *TerminalDrawer) write(text string) {
	fmt.Println(text)
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

func (s *TerminalDrawer) clearScreen() {
	fmt.Println("\033[H\033[2J")
}

// User input logic
type UserMessage struct {
	Point
	Command string
}

type UserInput interface {
	getInput() (UserMessage, error)
}

type TerminalInput struct {
	input UserInput
}

func (t *TerminalInput) getInput() (UserMessage, error) {
	var input string
	_, err := fmt.Scanln(&input)
	userMessage := UserMessage{}

	if len(input) == 3 && input[0] >= '1' && input[0] <= '9' && input[2] >= '1' && input[2] <= '9' && err == nil {
		userMessage.X = int(input[0]-'0') - 1
		userMessage.Y = int(input[2]-'0') - 1
	} else if len(input) == 1 && err == nil {
		userMessage.Command = input
	} else {
		return userMessage, errors.New("Invalid input. Please select cell coordinates in format X,Y")
	}

	return userMessage, nil
}

func main() {
	gameMap := &GameMap{}
	gameMap.Init(height, width, bombs_count)

	var drawHandler Drawer
	var inputHandler UserInput

	drawHandler = &TerminalDrawer{}
	inputHandler = &TerminalInput{}

	for {
		drawHandler.drawMap(gameMap.grid, gameMap.height, gameMap.width)

		userMessage, err := inputHandler.getInput()

		drawHandler.clearScreen()

		if err != nil {
			drawHandler.write(err.Error())
			continue
		}

		// Handle user input
		switch userMessage.Command {
		case "q":
			return
		case "r":
			gameMap.restart()
		case "":
			// Handle lose
			if gameMap.openCell(userMessage.Point) {
				drawHandler.write("Game over. Game restarted.")
				gameMap.restart()
				break
			}

			drawHandler.write("Your choise is: " + strconv.Itoa(userMessage.X+1) + " " + strconv.Itoa(userMessage.Y+1))
			gameMap.openZeros(userMessage.Point)

			// Handle victory
			if gameMap.checkWin() {
				drawHandler.write("Victory! Game restarted.")
				gameMap.restart()
			}
		default:
			drawHandler.write("Invalid command")
		}
	}
}
