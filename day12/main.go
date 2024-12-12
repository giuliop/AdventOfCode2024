package main

import (
	"bufio"
	"log"
	"os"
)

// PART 1
// Input is a grid of letters, indicating garden plots. Contiguous same letter
// plots form regions. The area of a region is the number of plots in it.
// The perimeter of a region is the number of plots' edges touching borders or
// other regions.
// The cost to fence a region is area * perimeter. Find the total cost to fence
// (note that shared borders across different regions are fenced twice).

func readInput() Garden {
	file, err := os.Open("input/day12")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var grid [][]byte

	for scanner.Scan() {
		grid = append(grid, []byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return Garden{grid, len(grid[0]) - 1, len(grid) - 1}
}

type Garden struct {
	plots      [][]byte
	maxX, maxY int
}

type Pos struct {
	x, y int
}

const outside = 0

func (g *Garden) neighbours(pos Pos) []Pos {
	res := make([]Pos, 0, 4)
	for _, dir := range []Pos{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		res = append(res, Pos{pos.x + dir.x, pos.y + dir.y})
	}
	return res
}

func (g *Garden) plotType(pos Pos) byte {
	if pos.x < 0 || pos.x > g.maxX || pos.y < 0 || pos.y > g.maxY {
		return outside
	}
	return g.plots[pos.y][pos.x]
}

// borders returns the number of neighbours of different plot type
func (g *Garden) borders(pos Pos) int {
	neighbours := g.neighbours(pos)
	borders := 0
	plotType := g.plotType(pos)
	for _, n := range neighbours {
		if g.plotType(n) != plotType {
			borders++
		}
	}
	return borders
}

// regionCost returns the regionCost to fence a region which includes start
func (g *Garden) regionCost(start Pos, visited map[Pos]bool) int {
	plotType := g.plotType(start)
	area := 0
	perimeter := 0
	toVisit := []Pos{start}
	visited[start] = true
	for len(toVisit) > 0 {
		p := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]
		area++
		perimeter += g.borders(p)
		for _, n := range g.neighbours(p) {
			if !visited[n] && g.plotType(n) == plotType {
				toVisit = append(toVisit, n)
				visited[n] = true
			}
		}
	}
	return area * perimeter
}

func answer1() int {
	garden := readInput()
	cost := 0
	visited := make(map[Pos]bool)
	for y, row := range garden.plots {
		for x := range row {
			if !visited[Pos{x, y}] {
				cost += garden.regionCost(Pos{x, y}, visited)
			}
		}
	}
	return cost
}

// -----------------------------------------------------------------------

// PART 2
// Now the perimeter counts contiguous sides as one, so for instance the
// perimeter of a 2x2 region is 4, not 8. Find the new total cost.

type Side int

const (
	Top Side = iota
	Right
	Bottom
	Left
)

var sideNames = []string{"Top", "Right", "Bottom", "Left"}

func (s Side) String() string {
	return sideNames[s]
}

type Border struct {
	pos  Pos
	side Side
}

func (p Pos) neighborAt(side Side) Pos {
	switch side {
	case Top:
		return Pos{p.x, p.y - 1}
	case Bottom:
		return Pos{p.x, p.y + 1}
	case Left:
		return Pos{p.x - 1, p.y}
	case Right:
		return Pos{p.x + 1, p.y}
	}
	return p // shouldn't happen
}

// contiguousBorders returns the list of contiguous borders for a border,
// that is with the same plot type and same side, all the way until a different
// plot type is found.
func (g *Garden) contiguousBorders(border Border) []Border {
	borders := make([]Border, 0)
	plotType := g.plotType(border.pos)
	toVisit := []Pos{border.pos}
	side := border.side
	visited := make(map[Pos]bool)
	for len(toVisit) > 0 {
		var neighbours []Pos
		pos := toVisit[len(toVisit)-1]
		toVisit = toVisit[:len(toVisit)-1]
		visited[pos] = true
		switch side {
		case Top, Bottom:
			neighbours = []Pos{
				{pos.x - 1, pos.y},
				{pos.x + 1, pos.y},
			}
		case Right, Left:
			neighbours = []Pos{
				{pos.x, pos.y - 1},
				{pos.x, pos.y + 1},
			}
		}
		for _, n := range neighbours {
			if g.plotType(n) == plotType &&
				g.plotType(n.neighborAt(side)) != plotType {
				if !visited[n] {
					toVisit = append(toVisit, n)
					visited[n] = true
				}
			}
		}
	}
	for pos := range visited {
		borders = append(borders, Border{pos, side})
	}
	return borders
}

func sideAt(direction Pos) Side {
	if direction.x == 0 {
		if direction.y == -1 {
			return Top
		}
		return Bottom
	}
	if direction.x == -1 {
		return Left
	}
	return Right
}

// bordersPart2 returns the borders with neighbours of different plots
func (g *Garden) bordersPart2(pos Pos) []Border {
	borders := make([]Border, 0)
	plotType := g.plotType(pos)
	for _, dir := range []Pos{{-1, 0}, {1, 0}, {0, -1}, {0, 1}} {
		neighbour := Pos{pos.x + dir.x, pos.y + dir.y}
		if g.plotType(neighbour) != plotType {
			side := sideAt(dir)
			borders = append(borders, Border{pos, side})
		}
	}
	return borders
}

// regionCost returns the regionCost to fence a region which includes start
func (g *Garden) regionCostPart2(start Pos, plotsVisited map[Pos]bool) int {
	plotType := g.plotType(start)
	area := 0
	perimeter := 0
	plotsToVisit := []Pos{start}
	plotsVisited[start] = true
	bordersVisited := make(map[Border]bool)
	for len(plotsToVisit) > 0 {
		p := plotsToVisit[len(plotsToVisit)-1]
		plotsToVisit = plotsToVisit[:len(plotsToVisit)-1]
		area++
		borders := g.bordersPart2(p)
		for _, b := range borders {
			if !bordersVisited[b] {
				perimeter++
				bordersVisited[b] = true
				for _, b2 := range g.contiguousBorders(b) {
					bordersVisited[b2] = true
				}
			}
		}
		for _, n := range g.neighbours(p) {
			if !plotsVisited[n] && g.plotType(n) == plotType {
				plotsToVisit = append(plotsToVisit, n)
				plotsVisited[n] = true
			}
		}
	}
	return area * perimeter
}

func answer2() int {
	garden := readInput()
	cost := 0
	visited := make(map[Pos]bool)
	for y, row := range garden.plots {
		for x := range row {
			if !visited[Pos{x, y}] {
				cost += garden.regionCostPart2(Pos{x, y}, visited)
			}
		}
	}
	return cost
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 1381056,
	2: 834828,
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
