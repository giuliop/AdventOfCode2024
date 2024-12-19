package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// PART 1
// You are inside a 71x71 menory grid, starting at the top-left corner 0,0.
// Your input is a list of lines such:
// 16,23
// representing corrupted memory cells which cannot be traversed.
// Considering the first 1024 input lines, what is the minimum number of
// steps to reach the bottom-right corner 69,69?

type Pos struct {
	x, y int
}

type World struct {
	corrupted  map[Pos]bool
	maxX, maxY int
	start      Pos
	end        Pos
}

func readInput() []Pos {
	file, err := os.Open("input/day18")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	corrupted := []Pos{}

	for scanner.Scan() {
		xStr, yStr, ok := strings.Cut(scanner.Text(), ",")
		if !ok {
			log.Fatal("Invalid input line: ", scanner.Text())
		}
		x, err := strconv.Atoi(xStr)
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(yStr)
		if err != nil {
			log.Fatal(err)
		}
		corrupted = append(corrupted, Pos{x, y})
	}
	return corrupted
}

type Predecessor struct {
	pos  Pos
	cost int
}

type Path struct {
	pos           Pos
	predecessor   Predecessor
	steps         int
	estimatedCost int
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (p Pos) estimateCost(end Pos) int {
	return abs(p.x-end.x) + abs(p.y-end.y)
}

type Paths map[int][]Path // steps -> paths

type PathManager struct {
	world        *World
	costsHeap    PriorityQueue
	paths        Paths
	predecessors map[Pos]Predecessor
}

// bestRoute returns the best route from start to end
// in reverse order
func (pm *PathManager) bestRoute() []Pos {
	_, ok := pm.predecessors[pm.world.end]
	if !ok {
		log.Fatal("Best route not found yet")
	}
	route := []Pos{}
	for pos := pm.world.end; ; pos = pm.predecessors[pos].pos {
		route = append(route, pos)
		if pos == pm.world.start {
			break
		}
	}
	return route
}

func (pm *PathManager) add(path Path) {
	cost := path.estimatedCost
	_, ok := pm.paths[cost]
	if !ok {
		pm.paths[cost] = []Path{}
		heap.Push(&pm.costsHeap, cost)
	}
	pm.paths[cost] = append(pm.paths[cost], path)
}

func (pm *PathManager) pop() Path {
	cost, ok := pm.costsHeap.Peek()
	if !ok {
		panic("No more paths")
	}
	path := pm.paths[cost][len(pm.paths[cost])-1]
	pm.paths[cost] = pm.paths[cost][:len(pm.paths[cost])-1]
	if len(pm.paths[cost]) == 0 {
		heap.Pop(&pm.costsHeap)
		delete(pm.paths, cost)
	}

	return path
}

func NewPathManager(w *World) *PathManager {
	startEstimatedCost := Pos{0, 0}.estimateCost(w.end)
	pq := PriorityQueue{startEstimatedCost}
	heap.Init(&pq)
	startPredecessor := Predecessor{Pos{-1, -1}, 0}
	startPath := Path{w.start, startPredecessor, 0, startEstimatedCost}
	paths := Paths{startEstimatedCost: []Path{startPath}}
	return &PathManager{w, pq, paths, map[Pos]Predecessor{}}
}

func (pm *PathManager) findPath() []Pos {
	visited := map[Pos]bool{}
	end := pm.world.end
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
		if path.pos == end {
			return pm.bestRoute()
		}
		for _, dir := range []Pos{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
			newPos := Pos{path.pos.x + dir.x, path.pos.y + dir.y}
			if newPos.x < 0 || newPos.x > pm.world.maxX ||
				newPos.y < 0 || newPos.y > pm.world.maxY ||
				pm.world.corrupted[newPos] ||
				visited[newPos] {
				continue
			}
			estimatedCost := path.steps + 1 + newPos.estimateCost(end)
			predecessor := Predecessor{path.pos, path.steps}
			newPath := Path{newPos, predecessor, path.steps + 1, estimatedCost}
			pm.add(newPath)
		}
	}
	return nil
}

func answer1() int {
	w := &World{
		corrupted: map[Pos]bool{}, maxX: 70, maxY: 70, start: Pos{0, 0}, end: Pos{70, 70},
	}
	corrupted := readInput()
	for i := 0; i < 1024; i++ {
		w.corrupted[corrupted[i]] = true
	}
	pm := NewPathManager(w)
	bestRoute := pm.findPath()
	return len(bestRoute) - 1
}

// -----------------------------------------------------------------------

// PART 2
// Now consider the other lines of the input, which is the first additional corrupted
// cell that cause the end to be unreachable?

func answer2() int {
	w := &World{
		corrupted: map[Pos]bool{}, maxX: 70, maxY: 70, start: Pos{0, 0}, end: Pos{70, 70},
	}
	corrupted := readInput()
	for i := 0; i < 1024; i++ {
		w.corrupted[corrupted[i]] = true
	}
	for i := 1024; i < len(corrupted); i++ {
		w.corrupted[corrupted[i]] = true
		pm := NewPathManager(w)
		bestRoute := pm.findPath()
		if bestRoute == nil {
			fmt.Printf("%d,%d\n", corrupted[i].x, corrupted[i].y)
			return corrupted[i].x*100 + corrupted[i].y
		}
	}
	panic("No solution found")
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 290,
	2: 6454,
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
