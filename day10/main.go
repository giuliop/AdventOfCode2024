package main

import (
	"bufio"
	"log"
	"os"
)

// PART 1
// The input is a grid of digits indicating the slope of a mountain.
// 0s represent trailheads and valid trails move N, E, S, W increasing
// the slope by 1 at each step until a 9. The score of a trailhead is
// number of 9s reachable from it. Sum all trailhead scores.

type World struct {
	grid [][]int
}

func readInput() *World {
	file, err := os.Open("input/day10")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	grid := make([][]int, 0)
	for scanner.Scan() {
		var row []int
		for _, c := range scanner.Text() {
			row = append(row, int(c-'0'))
		}
		grid = append(grid, row)
	}

	return &World{grid: grid}
}

type Pos struct {
	x, y int
}

func (p Pos) isValid(w *World) bool {
	return p.x >= 0 && p.x < len(w.grid[0]) && p.y >= 0 && p.y < len(w.grid)
}

func (w *World) slope(p Pos) int {
	return w.grid[p.y][p.x]
}

func (w *World) trailHeads() []Pos {
	var trailHeads []Pos
	for y := 0; y < len(w.grid); y++ {
		for x := 0; x < len(w.grid[0]); x++ {
			pos := Pos{x, y}
			if w.slope(pos) == 0 {
				trailHeads = append(trailHeads, pos)
			}
		}
	}
	return trailHeads
}

var directions = []Pos{{0, -1}, {1, 0}, {0, 1}, {-1, 0}}

func (w *World) nextSteps(p Pos) []Pos {
	steps := make([]Pos, 0)
	slope := w.slope(p)
	if slope == 9 {
		return steps
	}
	for _, d := range directions {
		next := Pos{p.x + d.x, p.y + d.y}
		if next.isValid(w) && w.slope(next) == slope+1 {
			steps = append(steps, next)
		}
	}
	return steps
}

// trails returns the number of of trails that reach 9s
// and the number of distinct 9s reached
func (w *World) trails(trailhead Pos) (trailCount, nineCount int) {
	if w.slope(trailhead) != 0 {
		log.Fatal("Not a trailhead")
	}
	trails := []Pos{trailhead}
	nines := make(map[Pos]bool)

	for len(trails) > 0 {
		trail := trails[len(trails)-1]
		trails = trails[:len(trails)-1]
		for _, next := range w.nextSteps(trail) {
			if w.slope(next) == w.slope(trail)+1 {
				if w.slope(next) == 9 {
					nines[next] = true
					trailCount++
				} else {
					trails = append(trails, next)
				}
			}
		}
	}
	nineCount = len(nines)
	return trailCount, nineCount
}

func answer1() int {
	w := readInput()
	score := 0
	for _, trailhead := range w.trailHeads() {
		_, nineCount := w.trails(trailhead)
		score += nineCount
	}
	return score
}

// -----------------------------------------------------------------------

// PART 2
// The rating of a trailhead is the number of distinct valid trails that start
// from it and end at a 9. Sum all trailhead ratings.

func answer2() int {
	w := readInput()
	rating := 0
	for _, trailhead := range w.trailHeads() {
		trailCount, _ := w.trails(trailhead)
		rating += trailCount
	}
	return rating
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 789,
	2: 1735,
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
