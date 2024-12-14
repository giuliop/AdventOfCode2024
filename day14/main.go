package main

import (
	"bufio"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
)

// PART 1
// The input is a list of strings defining robots such as:
// p=51,21 v=67,-50
// where p is the position and v is the velocity (tiles per second)
// The robots move in a 101 tiles wide and 103 tiles tall grid where
// the top left corner is (0,0) and the bottom right corner is (100,102)
// The robots move in a straight line at a constant speed and wrap around
// the edges of the grid. They can overlap without any problem.
// Simulate the robots moving for 100 seconds and count how many robots
// are in each quadrant of the grid (robots in the middle row or column are
// not counted). Return the product of the number of robots in each quadrant.

const (
	width  = 101
	height = 103
)

type Robot struct {
	x, y, vx, vy int
}

func parseRobot(s string) Robot {
	regex := regexp.MustCompile(`p=(\d+),(\d+) v=(-?\d+),(-?\d+)`)
	matches := regex.FindStringSubmatch(s)
	x, _ := strconv.Atoi(matches[1])
	y, _ := strconv.Atoi(matches[2])
	vx, _ := strconv.Atoi(matches[3])
	vy, _ := strconv.Atoi(matches[4])
	return Robot{x, y, vx, vy}
}

func readInput() []Robot {
	file, err := os.Open("input/day14")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	robots := make([]Robot, 0)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		robots = append(robots, parseRobot(scanner.Text()))
	}
	return robots
}

// positionAt returns the position of the robot at time t
func (r *Robot) positionAt(t int) (int, int) {
	x := ((r.x+r.vx*t)%width + width) % width
	y := ((r.y+r.vy*t)%height + height) % height
	return x, y
}

// quadrant returns the quadrant of the grid where the point (x,y) is
// if the point is in the middle row or column, return 0
func quadrant(x, y int) int {
	if x < width/2 && y < height/2 {
		return 1
	}
	if x > width/2 && y < height/2 {
		return 2
	}
	if x < width/2 && y > height/2 {
		return 3
	}
	if x > width/2 && y > height/2 {
		return 4
	}
	return 0
}

func answer1() int {
	robots := readInput()
	quadrantCounts := make(map[int]int)
	for _, r := range robots {
		x, y := r.positionAt(100)
		quadrant := quadrant(x, y)
		quadrantCounts[quadrant]++
	}
	return quadrantCounts[1] * quadrantCounts[2] * quadrantCounts[3] * quadrantCounts[4]
}

// -----------------------------------------------------------------------

// PART 2
// What is the fewest number of seconds that must elapse for the robots to
// arrange themselves in the shape of a Christmas tree?

// display prints the grid with the robots at time t
func display(robots []Robot, t int) {
	// Determine bounding box at time t
	minX, minY, maxX, maxY := math.MaxInt, math.MaxInt, 0, 0
	positions := make(map[[2]int]bool)

	for _, r := range robots {
		x, y := r.positionAt(t)
		positions[[2]int{x, y}] = true
		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	// Print only the bounding box region
	for y := minY; y <= maxY; y++ {
		line := make([]byte, maxX-minX+1)
		for x := minX; x <= maxX; x++ {
			if positions[[2]int{x, y}] {
				line[x-minX] = '#'
			} else {
				line[x-minX] = ' '
			}
		}
		println(string(line))
	}
}

// score calculates the size of the largest contiguous cluster
func score(positions map[[2]int]bool) int {
	visited := make(map[[2]int]bool)
	var maxClusterSize int

	// Helper function to perform depth-first search and calculate cluster size
	var dfs func(pos [2]int) int
	dfs = func(pos [2]int) int {
		if !positions[pos] || visited[pos] {
			return 0
		}
		visited[pos] = true
		clusterSize := 1
		directions := [][2]int{
			{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		}
		for _, d := range directions {
			neighbor := [2]int{pos[0] + d[0], pos[1] + d[1]}
			clusterSize += dfs(neighbor)
		}
		return clusterSize
	}
	for pos := range positions {
		if !visited[pos] && positions[pos] {
			clusterSize := dfs(pos)
			if clusterSize > maxClusterSize {
				maxClusterSize = clusterSize
			}
		}
	}
	return maxClusterSize
}

func answer2() int {
	robots := readInput()
	maxScore := 0
	minT := 0
	for t := 0; t < 10000; t++ {
		positions := make(map[[2]int]bool)
		for _, r := range robots {
			x, y := r.positionAt(t)
			positions[[2]int{x, y}] = true
		}
		score := score(positions)
		if score > maxScore {
			maxScore = score
			minT = t
		}
	}
	display(robots, minT)
	return minT
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 226548000,
	2: 7753,
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
