package main

import (
	"bufio"
	"log"
	"os"
)

// PART 1
// Your input is a maze with walls '#' and open paths '.'.
// Start is 'S' and end is 'E'. There is only one path from start to end.
// You can "cheat" by removing walls for two consecutive steps, but only once.
// How many different "cheats" do save you at least 100 steps?
// Uniquely identify a cheat by its start,end pair: start is the position you are in before
// activating the cheat, and end is the first position you are in when you don't need the cheat
// anymore or when the cheat is spent.

type Pos struct {
	x, y int
}

var nullPos = Pos{-1, -1}

type World struct {
	walls [][]bool
	start Pos
	end   Pos
	maxX  int
	maxY  int
}

func readInput() World {
	file, err := os.Open("input/day20")
	// file, err := os.Open("input/day20_test")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()

	walls := [][]bool{}
	var start, end Pos
	scanner := bufio.NewScanner(file)

	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]bool, len(line))
		for x, char := range line {
			switch char {
			case '#':
				row[x] = true
			case 'S':
				start = Pos{x, y}
			case 'E':
				end = Pos{x, y}
			case '.':
				row[x] = false
			default:
				log.Fatalf("unexpected character '%c' at (%d, %d)", char, x, y)
			}
		}
		walls = append(walls, row)
		y++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	return World{walls: walls, start: start, end: end,
		maxX: len(walls[0]) - 1, maxY: len(walls) - 1}
}

func (w *World) isWall(p Pos) bool {
	return w.walls[p.y][p.x]
}

type Predecessor struct {
	pos  Pos
	cost int
}

type Path struct {
	pos         Pos
	predecessor Predecessor
	steps       int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p Pos) distanceTo(dest Pos) int {
	return abs(p.x-dest.x) + abs(p.y-dest.y)
}

type PathManager struct {
	world        *World
	paths        []Path
	predecessors map[Pos]Predecessor
}

type Route struct {
	pos  map[int]Pos // steps -> pos
	step map[Pos]int // pos -> steps
}

// route returns the route from end to start
func (pm *PathManager) route() Route {
	_, ok := pm.predecessors[pm.world.end]
	if !ok {
		log.Fatal("Best route not found yet")
	}
	positions := []Pos{}
	for pos := pm.world.end; ; pos = pm.predecessors[pos].pos {
		positions = append(positions, pos)
		if pos == pm.world.start {
			break
		}
	}
	route := Route{map[int]Pos{}, map[Pos]int{}}
	for i := 0; i < len(positions); i++ {
		route.pos[i] = positions[len(positions)-1-i]
		route.step[positions[len(positions)-1-i]] = i
	}
	return route
}

func (r Route) len() int {
	return len(r.pos) - 1
}

func (pm *PathManager) add(path Path) {
	pm.paths = append(pm.paths, path)
}

func (pm *PathManager) pop() Path {
	path := pm.paths[0]
	pm.paths = pm.paths[1:]
	return path
}

func NewPathManager(w *World) *PathManager {
	startPath := Path{
		pos:         w.start,
		steps:       0,
		predecessor: Predecessor{nullPos, 0},
	}
	paths := []Path{startPath}
	return &PathManager{w, paths, map[Pos]Predecessor{}}
}

func (pm *PathManager) findPath() Route {
	visited := map[Pos]bool{}
	for len(pm.paths) > 0 {
		path := pm.pop()
		visited[path.pos] = true
		if predecessor, ok := pm.predecessors[path.pos]; ok {
			if predecessor.cost > path.predecessor.cost {
				pm.predecessors[path.pos] = path.predecessor
			}
		} else {
			pm.predecessors[path.pos] = path.predecessor
		}
		if path.pos == pm.world.end {
			return pm.route()
		}
		for _, dir := range []Pos{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			newPos := Pos{path.pos.x + dir.x, path.pos.y + dir.y}
			if visited[newPos] || pm.world.isWall(newPos) {
				continue
			}
			predecessor := Predecessor{path.pos, path.steps}
			newPath := Path{newPos, predecessor, path.steps + 1}
			pm.add(newPath)
		}
	}
	panic("No path found")
}

// countCheats returns the number of routes that save at least 100 steps
// by removing walls for two steps
func countCheats(route Route, minStepsToSave int, cheatDuration int) int {
	cheatsCount := 0
	// cheats have to start and end in a position along the ruote.
	// to get all possible starts, we iterate over the first len(route)-minStepsToSave
	// positions in the route and check if there is another pos in the route no more than
	// cheatDuration steps away that can be the end of the cheat.
	maxSteps := route.len() - minStepsToSave
	for i := 0; i < maxSteps-1; i++ {
		cheatStart := route.pos[i]
		for j := i + minStepsToSave; j < len(route.pos); j++ {
			cheatEnd := route.pos[j]
			distance := cheatStart.distanceTo(cheatEnd)
			if distance <= cheatDuration && i+distance+route.len()-j <= maxSteps {
				cheatsCount++
			}
		}
	}
	return cheatsCount
}

func answer1() int {
	w := readInput()
	minStepsToSave := 100
	cheatDuration := 2
	pm := NewPathManager(&w)
	route := pm.findPath()
	return countCheats(route, minStepsToSave, cheatDuration)
}

// -----------------------------------------------------------------------

// PART 2
// Now cheats last 20 steps. How many different "cheats" do save you at least 100 steps?

func answer2() int {
	w := readInput()
	minStepsToSave := 100
	cheatDuration := 20
	pm := NewPathManager(&w)
	route := pm.findPath()
	return countCheats(route, minStepsToSave, cheatDuration)
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 1511,
	2: 1020507,
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
