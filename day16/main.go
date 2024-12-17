package main

import (
	"container/heap"
	"io"
	"log"
	"math"
	"os"
)

// PART 1
// We have a maze with walls '#' and open spaces '.', start 'S' and end 'E'.
// We start facing east and need to find the lowest cost path to the end.
// A move forward costs 1, a 90-degree turn costs 1000.

type Pos struct {
	x, y int
}

type Walls map[Pos]bool

type World struct {
	walls Walls
	start Pos
	end   Pos
}

func readInput() World {
	file, err := os.Open("input/day16")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	walls := Walls{}
	world := World{walls: walls}

	y, x := 0, 0
	for {
		b := make([]byte, 1)
		_, err := file.Read(b)

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}

		switch b[0] {
		case '\n':
			y++
			x = 0
			continue
		case '#':
			walls[Pos{x, y}] = true
		case 'S':
			world.start = Pos{x, y}
		case 'E':
			world.end = Pos{x, y}
		}
		x++
	}
	return world
}

type PathState struct {
	pos Pos
	dir Pos
}

type Path struct {
	pos     Pos
	dir     Pos
	visited map[Pos]bool
}

var (
	north = Pos{0, -1}
	east  = Pos{1, 0}
	south = Pos{0, 1}
	west  = Pos{-1, 0}
)

type Paths map[PathState]map[Pos]bool // pathState -> visited

type PathManager struct {
	world     *World
	costsHeap PriorityQueue
	paths     map[int]Paths     // cost -> (pathState -> visited)
	costs     map[PathState]int // pathState -> cost
}

func NewPathManager(w *World) *PathManager {
	pq := PriorityQueue{0}
	heap.Init(&pq)
	initialPathState := PathState{w.start, east}
	paths := map[int]Paths{0: {PathState{w.start, east}: map[Pos]bool{}}}
	costs := map[PathState]int{initialPathState: 0}
	return &PathManager{w, pq, paths, costs}
}

func (pm *PathManager) add(path Path, cost int) {
	pathState := PathState{path.pos, path.dir}
	if currentCost, ok := pm.costs[pathState]; ok {
		if cost > currentCost {
			return
		}
		if cost == currentCost {
			// save the visited tiles for part 2
			for pos := range path.visited {
				pm.paths[currentCost][pathState][pos] = true
			}
			return
		}
		// we need to remove the old path
		delete(pm.costs, pathState)
		delete(pm.paths[currentCost], pathState)
		if len(pm.paths[currentCost]) == 0 {
			delete(pm.paths, currentCost)
			heap.Remove(&pm.costsHeap, pm.costsHeap.FindIndex(currentCost))
		}
	}
	_, ok := pm.paths[cost]
	if !ok {
		heap.Push(&pm.costsHeap, cost)
		pm.paths[cost] = map[PathState]map[Pos]bool{}
	}
	pm.paths[cost][pathState] = path.visited
	pm.costs[pathState] = cost
}

func (pm *PathManager) pop() (Path, int, bool) {
	cost, ok := pm.costsHeap.Peek()
	if !ok {
		return Path{}, cost, false
	}
	if len(pm.paths[cost]) == 0 {
		panic("No path found for existing cost")
	}
	var pathState PathState
	var visited map[Pos]bool
	for k, v := range pm.paths[cost] {
		pathState = k
		visited = v
		break
	}
	path := Path{pathState.pos, pathState.dir, visited}
	delete(pm.paths[cost], pathState)
	if len(pm.paths[cost]) == 0 {
		delete(pm.paths, cost)
		heap.Pop(&pm.costsHeap)
	}
	delete(pm.costs, PathState{path.pos, path.dir})
	return path, cost, true
}

const (
	moveCost = 1
	turnCost = 1000
)

func clone(visited map[Pos]bool) map[Pos]bool {
	newVisited := map[Pos]bool{}
	for k, v := range visited {
		newVisited[k] = v
	}
	return newVisited
}

func answer1() int {
	w := readInput()
	pm := NewPathManager(&w)
	for len(pm.paths) > 0 {
		path, cost, _ := pm.pop()
		if path.pos == w.end {
			return cost
		}
		path.visited[path.pos] = true
		for _, dir := range []Pos{north, east, south, west} {
			newPos := Pos{path.pos.x + dir.x, path.pos.y + dir.y}
			newCost := cost + moveCost
			if dir != path.dir {
				newCost += turnCost
			}
			if !(w.walls[newPos]) && !path.visited[newPos] {
				newPath := Path{newPos, dir, clone(path.visited)}
				pm.add(newPath, newCost)
			}
		}
	}
	log.Fatal("No path found")
	return -1
}

// -----------------------------------------------------------------------

// PART 2
// now find all the tiles that are part of at least one of the optimal
// paths (i.e., lowest cost) from start to end

func answer2() int {
	w := readInput()
	pm := NewPathManager(&w)
	bestCost := math.MaxInt
	bestTiles := map[Pos]bool{}
	for len(pm.paths) > 0 {
		path, cost, _ := pm.pop()
		path.visited[path.pos] = true
		if path.pos == w.end {
			bestCost = cost
			for pos := range path.visited {
				bestTiles[pos] = true
			}
		}
		for _, dir := range []Pos{north, east, south, west} {
			newPos := Pos{path.pos.x + dir.x, path.pos.y + dir.y}
			newCost := cost + moveCost
			if dir != path.dir {
				newCost += turnCost
			}
			if !(w.walls[newPos]) && !path.visited[newPos] && !(newCost > bestCost) {
				newPath := Path{newPos, dir, clone(path.visited)}
				pm.add(newPath, newCost)
			}
		}
	}
	return len(bestTiles)
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 102488,
	2: 559,
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
