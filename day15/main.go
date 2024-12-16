package main

import (
	"bufio"
	"log"
	"os"
)

// PART 1
// Your input is a map with walls '#', boxes 'o', open space '.', and the robot
// starting position '@', followed by a blank line and then lines of robot's
// movements: `<,>,v,^'. Boxes are pushed by the robot and if a move would cause
// either the robot or a box to hit a wall, the move is ignored.
// Process all the moves and sum all the boxes' gps coordinates calculated as
// y*100 + x (top-left corner is 0,0). Return the sum.

type Move byte
type Pos struct {
	x, y int
}
type World struct {
	grid       [][]byte
	maxX, maxY int
	robotPos   Pos
}

const (
	robot byte = '@'
	box   byte = 'O'
	wall  byte = '#'
	empty byte = '.'
)

func readInput() (World, []Move) {
	file, err := os.Open("input/day15")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	// read world map
	var grid [][]byte
	for scanner.Scan() {
		line := []byte(scanner.Text())
		if len(line) == 0 {
			break
		}
		grid = append(grid, line)
	}
	// read robot's moves
	var moves []Move
	for scanner.Scan() {
		line := []Move(scanner.Text())
		moves = append(moves, line...)
	}
	w := World{grid, len(grid[0]) - 1, len(grid) - 1, Pos{}}
	w.markRobotStartPosition()
	return w, moves
}

func (w *World) markRobotStartPosition() Pos {
	for y, row := range w.grid {
		for x, cell := range row {
			if cell == robot {
				w.robotPos = Pos{x, y}
				w.grid[y][x] = empty
				return w.robotPos
			}
		}
	}
	panic("Robot not found in map")
}

func dir(m Move) (int, int) {
	switch m {
	case '<':
		return -1, 0
	case '>':
		return 1, 0
	case '^':
		return 0, -1
	case 'v':
		return 0, 1
	}
	return 0, 0
}

func (w *World) tile(p Pos) byte {
	if p.x < 0 || p.x > w.maxX || p.y < 0 || p.y > w.maxY {
		panic("Out of bounds")
	}
	return w.grid[p.y][p.x]
}

func (w *World) makeMove(m Move) {
	dx, dy := dir(m)
	newPos := Pos{w.robotPos.x + dx, w.robotPos.y + dy}
	switch w.tile(newPos) {
	case wall:
		return
	case box:
		newBoxPos := Pos{newPos.x + dx, newPos.y + dy}
		for w.tile(newBoxPos) != wall {
			if w.tile(newBoxPos) == empty {
				w.grid[newPos.y][newPos.x] = empty
				w.grid[newBoxPos.y][newBoxPos.x] = box
				w.robotPos = newPos
				break
			}
			newBoxPos.x += dx
			newBoxPos.y += dy
		}
	case empty:
		w.robotPos = newPos
	}
}

func (w *World) gps() int {
	sum := 0
	for y, row := range w.grid {
		for x, cell := range row {
			if cell == box {
				sum += y*100 + x
			}
		}
	}
	return sum
}

func answer1() int {
	w, moves := readInput()
	for _, m := range moves {
		w.makeMove(m)
	}
	return w.gps()
}

// -----------------------------------------------------------------------

// PART 2
// Now the map is twice as large and boxes occupy two cells. The new map can be
// obtained from the old map replacing each cell like this:
//   . -> ..
//   O -> []
//   @ -> @.
//   # -> ##
// The robot's moves stays the same and the gps is calculated as before using the
// smaller x coordinate of the box. Return the sum of the boxes' gps coordinates.

const (
	leftBoxSide  byte = '['
	rightBoxSide byte = ']'
)

func makeWorldPart2(w *World) *World {
	w2 := World{
		grid:     make([][]byte, len(w.grid)),
		maxX:     2*(w.maxX+1) - 1,
		maxY:     w.maxY,
		robotPos: Pos{w.robotPos.x * 2, w.robotPos.y},
	}
	for y, row := range w.grid {
		newRow := []byte{}
		for _, cell := range row {
			switch cell {
			case empty:
				newRow = append(newRow, empty, empty)
			case box:
				newRow = append(newRow, leftBoxSide, rightBoxSide)
			case robot:
				newRow = append(newRow, robot, empty)
			case wall:
				newRow = append(newRow, wall, wall)
			}
		}
		w2.grid[y] = newRow
	}
	return &w2
}

// moveVertically tries to move newPos vertically and returns a bool indicating
// if the move was successful and the updated list of world changes
func (w *World) moveVertically(pos Pos, dy int) (bool, map[Pos]byte) {
	// we keep track of boxes to move in toMove where we add the left side of the box
	if w.tile(pos) == rightBoxSide {
		pos.x--
	}
	toMove := []Pos{pos}
	changes := map[Pos]byte{}

	for len(toMove) > 0 {
		newToMove := map[Pos]bool{}
		for _, blockLeftSide := range toMove {
			blockRightside := Pos{blockLeftSide.x + 1, blockLeftSide.y}
			for _, blockSide := range []Pos{blockLeftSide, blockRightside} {
				changes[Pos{blockSide.x, blockSide.y + dy}] = w.tile(blockSide)
				if changes[blockSide] == 0 {
					changes[blockSide] = empty
				}
				newPos := Pos{blockSide.x, blockSide.y + dy}
				switch w.tile(newPos) {
				case wall:
					return false, nil
				case leftBoxSide:
					newToMove[newPos] = true
				case rightBoxSide:
					newToMove[Pos{newPos.x - 1, newPos.y}] = true
				}
			}
		}
		toMove = []Pos{}
		for pos := range newToMove {
			toMove = append(toMove, pos)
		}
	}
	return true, changes
}

// moveVertically tries to move newPos horizontally and returns a bool indicating
// if the move was successful and the updated list of world changes
func (w *World) moveHorizontally(pos Pos, dx int) (bool, map[Pos]byte) {
	oppositeSide := map[byte]byte{
		leftBoxSide:  rightBoxSide,
		rightBoxSide: leftBoxSide}
	changes := map[Pos]byte{}
	changes[pos] = empty
	pos.x += dx
	for w.tile(pos) != wall {
		if w.tile(pos) == empty {
			w.grid[pos.y][pos.x] = empty
			changes[pos] = oppositeSide[changes[Pos{pos.x - dx, pos.y}]]
			return true, changes
		}
		changes[pos] = oppositeSide[w.tile(pos)]
		pos.x += dx
	}
	return false, nil
}

func (w *World) makeMovePart2(m Move) {
	dx, dy := dir(m)
	newPos := Pos{w.robotPos.x + dx, w.robotPos.y + dy}
	switch w.tile(newPos) {
	case wall:
		return
	case leftBoxSide, rightBoxSide:
		var ok bool
		var changes map[Pos]byte
		if dx != 0 {
			ok, changes = w.moveHorizontally(newPos, dx)
		} else {
			ok, changes = w.moveVertically(newPos, dy)
		}
		if ok {
			w.robotPos = newPos
			for pos, tile := range changes {
				w.grid[pos.y][pos.x] = tile
			}
		}
	case empty:
		w.robotPos = newPos
	}
}

func (w *World) gpsPart2() int {
	sum := 0
	for y, row := range w.grid {
		for x, cell := range row {
			if cell == leftBoxSide {
				sum += y*100 + x
			}
		}
	}
	return sum
}

func answer2() int {
	w, moves := readInput()
	w2 := makeWorldPart2(&w)
	for _, m := range moves {
		w2.makeMovePart2(m)
	}
	return w2.gpsPart2()
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 1552463,
	2: 1554058,
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
