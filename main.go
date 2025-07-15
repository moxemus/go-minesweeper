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
	//var input string

	// for {
	// 	_, err := fmt.Scanln(&input)
	// 	if err != nil {
	// 		fmt.Println("Error:", err)
	// 		continue
	// 	} else {
	//         fmt.Println("Answer: ", input)
	//     }

	// }

	// fmt.Println("Game over.")

	clearScreen()
	initMap(&cells)
	drawMap(cells)
}

func initMap(vals *[height][width]int) {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			vals[i][j] = rand.Intn(10)
		}
	}
}

func drawMap(vals [height][width]int) {
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
			//fmt.Printf("%3d ", vals[i][j])
			fmt.Printf("%3c ", '*')
		}

		fmt.Println()
	}
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
