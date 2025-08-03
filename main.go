package main

import (
	"fmt"
	"math/rand"
)

const height = 9
const width = 9
const bombs_count = 5

func main() {
	var cells [height][width]int
	var visible [height][width]bool
	var input string

	initMap(&cells, &visible)

	for {
		drawMap(cells, visible)

		_, err := fmt.Scanln(&input)
		clearScreen()

		if len(input) == 3 && input[0] >= '1' && input[0] <= '9' && input[2] >= '1' && input[2] <= '9' && err == nil {
			// Coordinates for current game
			x := int(input[0] - '0')
			y := int(input[2] - '0')
			var result = cells[x-1][y-1]

			if result == -1 {
				fmt.Println("Game over. Game restarted.")
				initMap(&cells, &visible)
			} else {
				fmt.Println("Your choise is: ", x, y, " result - ", result)
				visible[x-1][y-1] = true
			}

		} else if input == "q" {
			break
		} else if input == "r" {
			initMap(&cells, &visible)
			continue
		} else {
			fmt.Println("Invalid input. Please select cell coordinates in format X,Y")
		}
	}
}

func initMap(vals *[height][width]int, visible *[height][width]bool) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			vals[i][j] = 0
		}
	}

	for i := 0; i < bombs_count; i++ {
		y := rand.Intn(height)
		x := rand.Intn(width)

		vals[x][y] = -1
		updateBombCells(vals, x, y)
	}

	// Clear visible
	for i := range visible {
		for j := range visible[i] {
			visible[i][j] = false
		}
	}
}

func updateBombCells(vals *[height][width]int, x, y int) {
	for dy := -1; dy <= 1; dy++ {
		for dx := -1; dx <= 1; dx++ {
			ny, nx := y+dy, x+dx

			if ny < 0 || ny >= height || nx < 0 || nx >= width || vals[nx][ny] == -1 {
				continue
			}

			vals[nx][ny]++
		}
	}
}

func drawMap(vals [height][width]int, visible [height][width]bool) {
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
			if !visible[j][i] {
				fmt.Printf("%3c ", '#')
			} else {
				if vals[j][i] == -1 {
					fmt.Printf("%3c ", 'x')
				} else {
					fmt.Printf("%3d ", vals[j][i])
				}
			}
		}

		fmt.Println()
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
