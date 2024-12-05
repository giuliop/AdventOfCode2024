package main

import (
	"bufio"
	"log"
	"os"
)

// PART 1

func readInput() [][]byte {
	file, err := os.Open("input/day4")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines [][]byte
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, []byte(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	return lines
}

type direction int

const (
	horizontalForward direction = iota
	horizontalBackward
	verticalForward
	verticalBackward
	diagonalUpRight
	diagonalUpLeft
	diagonalDownRight
	diagonalDownLeft
)

type position struct {
	row, col int
}

// getNextThreePositions returns the next three positions in the given direction
func getNextThreePositions(p position, dir direction) []position {
	switch dir {
	case horizontalForward:
		return []position{{p.row, p.col + 1}, {p.row, p.col + 2}, {p.row, p.col + 3}}
	case horizontalBackward:
		return []position{{p.row, p.col - 1}, {p.row, p.col - 2}, {p.row, p.col - 3}}
	case verticalForward:
		return []position{{p.row + 1, p.col}, {p.row + 2, p.col}, {p.row + 3, p.col}}
	case verticalBackward:
		return []position{{p.row - 1, p.col}, {p.row - 2, p.col}, {p.row - 3, p.col}}
	case diagonalUpRight:
		return []position{{p.row - 1, p.col + 1}, {p.row - 2, p.col + 2}, {p.row - 3, p.col + 3}}
	case diagonalUpLeft:
		return []position{{p.row - 1, p.col - 1}, {p.row - 2, p.col - 2}, {p.row - 3, p.col - 3}}
	case diagonalDownRight:
		return []position{{p.row + 1, p.col + 1}, {p.row + 2, p.col + 2}, {p.row + 3, p.col + 3}}
	case diagonalDownLeft:
		return []position{{p.row + 1, p.col - 1}, {p.row + 2, p.col - 2}, {p.row + 3, p.col - 3}}
	}
	return nil
}

// isXMAS returns true if the given position and direction form the word 'XMAS'
func isXMAS(p position, dir direction, lines [][]byte) bool {
	if lines[p.row][p.col] != 'X' {
		return false
	}
	positions := getNextThreePositions(p, dir)
	letters := []byte{'M', 'A', 'S'}
	for i, pos := range positions {
		if pos.row < 0 || pos.row >= len(lines) || pos.col < 0 || pos.col >= len(lines[pos.row]) {
			return false
		}
		if lines[pos.row][pos.col] != letters[i] {
			return false
		}
	}
	return true
}

func answer1() int {
	// input is a list of lines of text. Find all 'XMAS' sequences, which can be horizontal,
	// vertical or diagonal, also backwards, and return the number of times it appears.
	sum := 0
	directions := []direction{horizontalForward, horizontalBackward, verticalForward,
		verticalBackward, diagonalUpRight, diagonalUpLeft, diagonalDownRight, diagonalDownLeft}
	for row, chars := range readInput() {
		for col, char := range chars {
			if char == 'X' {
				for _, dir := range directions {
					if isXMAS(position{row, col}, dir, readInput()) {
						sum++
					}
				}
			}
		}
	}
	return sum
}

// -----------------------------------------------------------------------

// PART 2

// getDiagonalAdjacent returns the adjacent positions in the given diagonal direction
// /from the given position. note that diagonalUpRight is the same as diagonalDownLeft,
// and diagonalUpLeft is the same as diagonalDownRight, just in the opposite order
func getDiagonalAdjacent(p position, dir direction) []position {
	switch dir {
	case diagonalUpRight:
		return []position{{p.row - 1, p.col + 1}, {p.row + 1, p.col - 1}}
	case diagonalUpLeft:
		return []position{{p.row - 1, p.col - 1}, {p.row + 1, p.col + 1}}
	case diagonalDownRight:
		return []position{{p.row + 1, p.col + 1}, {p.row - 1, p.col - 1}}
	case diagonalDownLeft:
		return []position{{p.row + 1, p.col - 1}, {p.row - 1, p.col + 1}}
	default:
		log.Fatal("Invalid direction")
		return nil
	}
}

// countMAS counts how many times the given position forms 'MAS' diagonally,
// assuming it is the central position of the word (i.e. the 'A')
func countMAS(p position, lines [][]byte) int {
	sum := 0
	if lines[p.row][p.col] != 'A' {
		return sum
	}
	if p.row == 0 || p.row == len(lines)-1 || p.col == 0 || p.col == len(lines[p.row])-1 {
		return sum
	}
	directions := []direction{diagonalUpRight, diagonalUpLeft, diagonalDownRight, diagonalDownLeft}
	letters := []byte{'M', 'S'}
	for _, dir := range directions {
		positions := getDiagonalAdjacent(p, dir)
		if lines[positions[0].row][positions[0].col] == letters[0] &&
			lines[positions[1].row][positions[1].col] == letters[1] {
			sum++
		}
	}
	return sum
}

func answer2() int {
	// now we need to find the 'MAS' words that cross line in the below diagram.
	// MAS can be written forward or backward.
	// M.S
	// .A.
	// M.S
	sum := 0
	for row, chars := range readInput() {
		for col, char := range chars {
			if char == 'A' {
				if countMAS(position{row, col}, readInput()) == 2 {
					sum++
				}
			}
		}
	}
	return sum
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 2569,
	2: 1998,
}

var answerFuncs = map[int]func() int{
	1: answer1,
	2: answer2,
}

func printAndTest(question int) {
	answer := answerFuncs[question]()
	correctAnswer, ok := correctAnswers[question]
	if ok && answer != correctAnswer {
		log.Fatal("Wrong answer, expected ", correctAnswer, " got ", answer)
	}
	println(answer)
}

func main() {
	// if no argument, run all answers, otherwise only part 1 or 2
	if len(os.Args) == 1 || os.Args[1] == "1" {
		printAndTest(1)
	}
	if len(os.Args) == 1 || os.Args[1] == "2" {
		printAndTest(2)
	}
	if len(os.Args) > 1 && os.Args[1] != "1" && os.Args[1] != "2" {
		println("Give 1 or 2 as argument, or no argument at all")
	}
}
