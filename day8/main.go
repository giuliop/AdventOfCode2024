package main

import (
	"io"
	"log"
	"os"
)

// PART 1
// We have a grid with antennas marked by symbols other than '.'. Each pair of
// same symbol antennas form an antinode in line with the two antennas at the
// point where one antenna is twice as far from the antinode as the other.
// Count the number of antinodes in the grid.

// Pos is a position in the grid
type Pos struct {
	X, Y int
}

type World struct {
	Antennas   map[byte][]Pos
	MaxX, MaxY int
}

func readInput() World {
	file, err := os.Open("input/day8")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	world := World{Antennas: map[byte][]Pos{}}

	y, x := 0, 0
	for {
		b := make([]byte, 1)
		_, err := file.Read(b)

		if err == io.EOF {
			world.MaxY = y - 1
			break
		}
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}

		switch b[0] {
		case '\n':
			world.MaxX = x - 1
			y++
			x = 0
			continue
		case '.':
			// empty space
		default:
			world.Antennas[b[0]] = append(world.Antennas[b[0]], Pos{x, y})
		}
		x++
	}
	return world
}

// getAntinodes returns the antinodes formed by the two antennas
func (w *World) getAntinodes(antenna1 Pos, antenna2 Pos) []Pos {
	antinodes := []Pos{}
	distX := antenna2.X - antenna1.X
	distY := antenna2.Y - antenna1.Y
	candidateAntinodes := []Pos{
		{antenna1.X - distX, antenna1.Y - distY},
		{antenna2.X + distX, antenna2.Y + distY},
	}
	for _, pos := range candidateAntinodes {
		if pos.X >= 0 && pos.X <= w.MaxX && pos.Y >= 0 && pos.Y <= w.MaxY {
			antinodes = append(antinodes, pos)
		}
	}
	return antinodes
}

func answer1() int {
	w := readInput()
	antinodes := map[Pos]bool{}
	for _, positions := range w.Antennas {
		for i, pos1 := range positions {
			for _, pos2 := range positions[i+1:] {
				for _, node := range w.getAntinodes(pos1, pos2) {
					antinodes[node] = true
				}
			}
		}
	}
	return len(antinodes)
}

// -----------------------------------------------------------------------

// PART 2
// Now antinodes occurs at any grid position exactly in line with two antennas,
// so each antenna in a pair is also an antinode. Count the number of antinodes.
// As an example, here we have 3 antennas 'T' and we mark with '#' the antinodes,
// the number of antinodes is 9 inclusing the antennas themselves:
// T....#....
// ...T......
// .T....#...
// .........#
// ..#.......
// ..........
// ...#......
// ..........
// ....#.....
// ..........

func (w *World) getAntinodesPart2(antenna1 Pos, antenna2 Pos) []Pos {
	antinodes := []Pos{}
	distX := antenna2.X - antenna1.X
	distY := antenna2.Y - antenna1.Y
	for i := 0; ; i++ {
		pos := Pos{antenna1.X - i*distX, antenna1.Y - i*distY}
		if pos.X < 0 || pos.X > w.MaxX || pos.Y < 0 || pos.Y > w.MaxY {
			break
		} else {
			antinodes = append(antinodes, pos)
		}
	}
	for i := 0; ; i++ {
		pos := Pos{antenna2.X + i*distX, antenna2.Y + i*distY}
		if pos.X < 0 || pos.X > w.MaxX || pos.Y < 0 || pos.Y > w.MaxY {
			break
		} else {
			antinodes = append(antinodes, pos)
		}
	}
	return antinodes
}

func answer2() int {
	w := readInput()
	antinodes := map[Pos]bool{}
	for _, positions := range w.Antennas {
		for i, pos1 := range positions {
			for _, pos2 := range positions[i+1:] {
				for _, node := range w.getAntinodesPart2(pos1, pos2) {
					antinodes[node] = true
				}
			}
		}
	}
	return len(antinodes)
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 413,
	2: 1417,
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
