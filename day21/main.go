package main

import (
	"bufio"
	"log"
	"os"
)

// PART 1
// We have numerical (left) and directional (left) keypads:

// +---+---+---+           +---+---+
// | 7 | 8 | 9 |           | ^ | A |
// +---+---+---+       +---+---+---+
// | 4 | 5 | 6 |       | < | v | > |
// +---+---+---+       +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
//     | 0 | A |
//     +---+---+

// You are typing on a dir keypad controlling robot 1.
// Robor 1 is typing on a dir keypad controlling robot 2.
// Robot 2 is typing on a dir keypad controlling robot 3.
// Robot 3 is typing on a num keypad controlling a locked door.

// Dir keypads control the robot's arm they are linked to, whrere 'A' makes
// the robot puch the button of the pad they are contrlling. Each robot's arm
// starts aiming at the 'A' button of the keypad they are controlling.
// Note that robots cannot ever aim at a "empty" button, lest they panic.

// Each line in your input (e.g., "029A") is a code to type in the num keypad.
// For each code, find the shortest sequence you can type to cause robot 3 to
// type the code. The `complexity` of a code is the product of the length of
// the shortest sequence and the numerical part of the code (ignoring leading 0s).
// Return the sum of the complexities of all codes.

type Code []int

// numerical returns the numerical part of the code
func (c Code) numerical() int {
	result := 0
	multiplier := 1
	for i := len(c) - 1; i >= 0; i-- {
		switch c[i] {
		case A:
			continue
		default:
			result += multiplier * (c[i])
			multiplier *= 10
		}
	}
	return result
}

const A int = 10

func readInput() []Code {
	file, err := os.Open("input/day21")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var codes []Code

	for scanner.Scan() {
		code := make([]int, 0)
		for _, char := range scanner.Text() {
			if char == 'A' {
				code = append(code, A)
			} else {
				// char is 0-9
				code = append(code, int(char)-48)
			}
		}
		codes = append(codes, code)
	}

	return codes
}

const (
	up int = iota
	down
	left
	right
)

type Pos struct {
	x, y int
}
type Keypad []Pos

const empty int = 11

var numKeypad = Keypad{
	// 0 1 2 3 4 5 6 7 8 9 A, empty is 11
	{1, 3}, {0, 2}, {1, 2}, {2, 2}, {0, 1}, {1, 1},
	{2, 1}, {0, 0}, {1, 0}, {2, 0}, {2, 3}, {0, 3},
}

var dirKeypad = Keypad{
	// up down left right ... A is index 10, empty 11, the same as numKeypad
	{1, 0}, {1, 1}, {0, 1}, {2, 1}, {9, 9}, {9, 9},
	{9, 9}, {9, 9}, {9, 9}, {9, 9}, {2, 0}, {0, 0},
}

func (k Keypad) isNum() bool {
	return k[empty] == Pos{0, 3}
}

func (k Keypad) isDir() bool {
	return k[empty] == Pos{0, 0}
}

const (
	numKeypadType int = iota
	dirKeypadType
)

// move returns the moves to go from one key to another
func (k Keypad) moveSets(from, to int) []int {
	distX := k[to].x - k[from].x
	distY := k[to].y - k[from].y

	var horizontalMoveSet []int
	for dist := distX; dist != 0; {
		if distX > 0 {
			horizontalMoveSet = append(horizontalMoveSet, right)
			dist--
		} else {
			horizontalMoveSet = append(horizontalMoveSet, left)
			dist++
		}
	}
	var verticalMoveSet []int
	for dist := distY; dist != 0; {
		if distY > 0 {
			verticalMoveSet = append(verticalMoveSet, down)
			dist--
		} else {
			verticalMoveSet = append(verticalMoveSet, up)
			dist++
		}
	}
	var res []int
	if horizontalMoveSet == nil || verticalMoveSet == nil {
		res = append(append(horizontalMoveSet, verticalMoveSet...), A)
		// check if one set traverses the empty space, only one ordering option
	} else if k.isNum() && k[to].x == 0 && k[from].y == 3 {
		res = append(append(verticalMoveSet, horizontalMoveSet...), A)
	} else if k.isNum() && k[to].y == 3 && k[from].x == 0 {
		res = append(append(horizontalMoveSet, verticalMoveSet...), A)
	} else if k.isDir() && k[to].x == 0 && k[from].y == 0 {
		res = append(append(verticalMoveSet, horizontalMoveSet...), A)
	} else if k.isDir() && k[to].y == 0 && k[from].x == 0 {
		res = append(append(horizontalMoveSet, verticalMoveSet...), A)
	} else {
		// if not, move horizontally first if going left, the most expensive move,
		// this optimizes the number of moves down the road
		if distX < 0 {
			res = append(append(horizontalMoveSet, verticalMoveSet...), A)
		} else {
			res = append(append(verticalMoveSet, horizontalMoveSet...), A)
		}
	}
	return res
}

type Memory map[[2]int][]int // [from, to] -> moves

// precomputeMoves returns a map of all possible movesets for all pairs of keys
// for the given keypad type
func precomputeMoves(keypadType int) Memory {
	mem := Memory{}
	var keypad Keypad
	if keypadType == numKeypadType {
		keypad = numKeypad
	} else {
		keypad = dirKeypad
	}
	for from := 0; from < 11; from++ {
		for to := 0; to < 11; to++ {
			if (keypad[from] == Pos{9, 9}) ||
				(keypad[to] == Pos{9, 9}) {
				continue
			}
			mem[[2]int{from, to}] = keypad.moveSets(from, to)
		}
	}
	return mem
}

// precomputeDirMoves returns a map of the moves to do on a dir keypad
// to move from a key to another (and press it) following the movies of the targetMem
func precomputeDirMoves(targetMem, dirMem Memory) Memory {
	mem := Memory{}
	for fromTo, moves := range targetMem {
		dirMoves := []int{}
		pos := A
		for _, move := range moves {
			dirMoves = append(dirMoves, dirMem[[2]int{pos, move}]...)
			pos = move
		}
		mem[fromTo] = dirMoves
	}
	return mem
}

func answer1() int {
	numpadMoves := precomputeMoves(numKeypadType)
	dirpadMoves := precomputeMoves(dirKeypadType)
	dirToNumMoves := precomputeDirMoves(numpadMoves, dirpadMoves)
	dirToDirToNumMoves := precomputeDirMoves(dirToNumMoves, dirpadMoves)

	res := 0
	codes := readInput()
	for _, code := range codes {
		sequenceLen := 0
		pos := A
		for _, key := range code {
			sequenceLen += len(dirToDirToNumMoves[[2]int{pos, key}])
			pos = key
		}
		res += code.numerical() * sequenceLen
	}
	return res
}

// -----------------------------------------------------------------------

// PART 2
// Now instead of 2 directional robots, we have 25 of them controlling each other.
// Find the new sum of complexities of all codes.
func answer2() int {
	numpadMoves := precomputeMoves(numKeypadType)
	dirpadMoves := precomputeMoves(dirKeypadType)
	dirToNumMoves := precomputeDirMoves(numpadMoves, dirpadMoves)

	res := 0
	codes := readInput()
	for _, code := range codes {
		sequenceLen := 0
		movesToCount := map[[2]int]int{}
		keyPos := A
		for _, key := range code {
			moves := dirToNumMoves[[2]int{keyPos, key}]
			movePos := A
			for _, move := range moves {
				movesToCount[[2]int{movePos, move}]++
				movePos = move
			}
			keyPos = key
		}

		for i := 1; i < 25; i++ {
			newMovesToCount := map[[2]int]int{}
			for fromTo, count := range movesToCount {
				moves := dirpadMoves[fromTo]
				movePos := A
				for _, move := range moves {
					newMovesToCount[[2]int{movePos, move}] += count
					movePos = move
				}
			}
			movesToCount = newMovesToCount
		}
		for _, count := range movesToCount {
			sequenceLen += count
		}
		res += code.numerical() * sequenceLen
	}
	return res
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 157908,
	2: 196910339808654,
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
