package main

import (
	"io"
	"log"
	"os"
)

// PART 1
// The world map is a grid with obstacles marked by '#', a starting position marked
// by '^', and empty spaces marked by '.'. The robot moves forward until it hits
// an obstacle, then turns right and continues. The task is to find out how many
// distinct positions the robot visits before it goes outside the grid.

// Pos is a position in the grid
type Pos = struct {
	X, Y int
}

// Obstacles[pos] is true if there is an obstacle at pos
type Obstacles = map[Pos]bool

type World struct {
	Obstacles  Obstacles
	MaxX, MaxY int
	Start      Pos
}

// Dir represents a direction
type Dir struct {
	dx int
	dy int
}

var (
	N = Dir{dx: 0, dy: -1} // North
	S = Dir{dx: 0, dy: 1}  // South
	E = Dir{dx: 1, dy: 0}  // East
	W = Dir{dx: -1, dy: 0} // West
)

func (m *World) posInFront(pos Pos, dir Dir) Pos {
	return Pos{pos.X + dir.dx, pos.Y + dir.dy}
}

func (m *World) isObstacle(pos Pos) bool {
	return m.Obstacles[pos]
}

func (m *World) isOutside(pos Pos) bool {
	return pos.X < 0 || pos.X > m.MaxX || pos.Y < 0 || pos.Y > m.MaxY
}

func readInput() World {
	file, err := os.Open("input/day6")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	obstacles := make(Obstacles)
	world := World{Obstacles: obstacles}

	y, x := 0, 0
	for {
		byteBuffer := make([]byte, 1)
		_, err := file.Read(byteBuffer)

		if err == io.EOF {
			world.MaxY = y - 1
			break
		}
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}

		switch byteBuffer[0] {
		case '\n':
			world.MaxX = x - 1
			y++
			x = 0
			continue
		case '#':
			obstacles[Pos{x, y}] = true
		case '^':
			world.Start = Pos{x, y}
		case '.':
			// empty space
		default:
			log.Fatalf("unexpected character: %c", byteBuffer[0])
		}
		x++
	}
	return world
}

func turnRight(dir Dir) Dir {
	switch dir {
	case N:
		return E
	case E:
		return S
	case S:
		return W
	case W:
		return N
	}
	log.Fatal("unexpected direction")
	return N
}

func answer1() int {
	w := readInput()
	dir := N
	visited := make(map[Pos]bool)
	for pos := w.Start; !w.isOutside(pos); {
		visited[pos] = true
		facing := w.posInFront(pos, dir)
		if w.isObstacle(facing) {
			dir = turnRight(dir)
		} else {
			pos = facing
		}
	}
	return len(visited)
}

// -----------------------------------------------------------------------

// PART 2
// Now we want to know in how many places we can put a new obstacle so that the
// robot gets stuck in a loop. We cannot put an obstacle on the starting position}

type State struct {
	Pos
	Dir
}

func clone(visited map[State]bool) map[State]bool {
	newVisited := make(map[State]bool)
	for k, v := range visited {
		newVisited[k] = v
	}
	return newVisited
}

type Path struct {
	visited     map[State]bool
	current     State
	newObstacle Pos
}

var nullPos = Pos{-1, -1}

func answer2() int {
	w := readInput()
	triedObstacles := make(map[Pos]bool)
	loopObstacles := make(map[Pos]bool)
	startPath := Path{
		visited:     make(map[State]bool),
		current:     State{Pos: w.Start, Dir: N},
		newObstacle: nullPos,
	}
	paths := []Path{startPath}
	for {
		if len(paths) == 0 {
			break
		}
		path := paths[len(paths)-1]
		paths = paths[:len(paths)-1]
		current := path.current
		if w.isOutside(current.Pos) {
			continue
		}
		if path.visited[current] {
			loopObstacles[path.newObstacle] = true
			continue
		}
		path.visited[current] = true
		facing := w.posInFront(current.Pos, current.Dir)
		if w.isObstacle(facing) || facing == path.newObstacle {
			path.current.Dir = turnRight(current.Dir)
		} else {
			if !triedObstacles[facing] && path.newObstacle == nullPos {
				triedObstacles[facing] = true
				newPath := Path{
					visited:     clone(path.visited),
					current:     current,
					newObstacle: facing,
				}
				newPath.visited[current] = false
				paths = append(paths, newPath)
			}
			path.current.Pos = facing
		}
		paths = append(paths, path)
	}
	return len(loopObstacles)
}

// ----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 4890,
	2: 1995,
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
