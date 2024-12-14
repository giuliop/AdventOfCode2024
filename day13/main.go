package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
)

// PART 1
// Input is a sequence of three-line blocks separated by empty lines such as:
// Button A: X+15, Y+26
// Button B: X+39, Y+21
// Prize: X=1061, Y=6652
// Each block represents a machine with two buttons and a prize. Button A costs 3
// to press, button B costs 1. To win the prize a combination of button presses
// must exists to move from origin (0,0) to the prize location.
// Find the minimum cost to win all the winnable prizes.

type Button struct {
	dx, dy int64
}

type Pos struct {
	x, y int64
}

type Machine struct {
	a, b  Button
	prize Pos
}

type Combination struct {
	a, b int64
}

func parseButton(s string) Button {
	regexX := regexp.MustCompile(`X\+(\d+)`)
	regexY := regexp.MustCompile(`Y\+(\d+)`)
	xMatch := regexX.FindStringSubmatch(s)
	yMatch := regexY.FindStringSubmatch(s)
	dx, err := strconv.Atoi(xMatch[1])
	if err != nil {
		log.Fatal(err)
	}
	dy, err := strconv.Atoi(yMatch[1])
	if err != nil {
		log.Fatal(err)
	}
	return Button{int64(dx), int64(dy)}
}

func parsePrize(s string) Pos {
	regexX := regexp.MustCompile(`X=(\d+)`)
	regexY := regexp.MustCompile(`Y=(\d+)`)
	xMatch := regexX.FindStringSubmatch(s)
	yMatch := regexY.FindStringSubmatch(s)
	x, err := strconv.Atoi(xMatch[1])
	if err != nil {
		log.Fatal(err)
	}
	y, err := strconv.Atoi(yMatch[1])
	if err != nil {
		log.Fatal(err)
	}
	return Pos{int64(x), int64(y)}
}

func readInput() []Machine {
	file, err := os.Open("input/day13")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var machines []Machine

	for scanner.Scan() {
		buttonA := scanner.Text()
		if !scanner.Scan() {
			break
		}
		buttonB := scanner.Text()
		if !scanner.Scan() {
			break
		}
		prize := scanner.Text()

		machine := Machine{
			a:     parseButton(buttonA),
			b:     parseButton(buttonB),
			prize: parsePrize(prize),
		}
		machines = append(machines, machine)

		// Skip empty line
		scanner.Scan()
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return machines
}

// extendedEuclid computes the gcd of a and b, as well as x,y
// such that a*x + b*y = gcd(a,b).
func extendedEuclid(a, b int64) (g, x, y int64) {
	if b == 0 {
		return a, 1, 0
	}
	g, x1, y1 := extendedEuclid(b, a%b)
	return g, y1, x1 - (a/b)*y1
}

// Checks if (a1,b1,c1) and (a2,b2,c2) define proportional equations:
// a1/a2 = b1/b2 = c1/c2 using cross multiplications to avoid floats.
func isProportional(a1, a2, b1, b2, c1, c2 int64) bool {
	// If (a2,b2,c2) = (0,0,0), then (a1,b1,c1) must also be (0,0,0)
	if a2 == 0 && b2 == 0 && c2 == 0 {
		return a1 == 0 && b1 == 0 && c1 == 0
	}
	if a1*b2 != a2*b1 {
		return false
	}
	if a1*c2 != a2*c1 {
		return false
	}
	if b1*c2 != b2*c1 {
		return false
	}
	return true
}

// Solves the system of Diophantine equations:
// a1*x + b1*y = c1
// a2*x + b2*y = c2
// Returns all nonnegative integer solutions.
// Thank you ChatGPT o1 for writing this function
func buttonCombinations(machine Machine) []Combination {
	a1, b1, c1 := machine.a.dx, machine.b.dx, machine.prize.x
	a2, b2, c2 := machine.a.dy, machine.b.dy, machine.prize.y

	var result []Combination
	D := a1*b2 - a2*b1

	if D != 0 {
		// Unique solution if it exists
		xNum := c1*b2 - c2*b1
		yNum := a1*c2 - a2*c1

		if xNum%D == 0 && yNum%D == 0 {
			x := xNum / D
			y := yNum / D
			if x >= 0 && y >= 0 {
				result = append(result, Combination{x, y})
			}
		}
		return result
	}

	// D = 0, lines might be dependent
	if !isProportional(a1, a2, b1, b2, c1, c2) {
		return result // no solutions
	}

	// Reduced to a single equation: a1*x + b1*y = c1
	g, x0, y0 := extendedEuclid(a1, b1)
	if c1%g != 0 {
		return result // no solutions
	}
	x0 *= c1 / g
	y0 *= c1 / g

	bdg := b1 / g
	adg := a1 / g

	// x0 + bdg*t >= 0 => t >= -x0/bdg
	var tMin int64
	if x0 < 0 {
		// ceiling division of (-x0) by bdg
		tMin = (-x0 + bdg - 1) / bdg
	} else {
		tMin = (-x0) / bdg
	}

	// y0 - adg*t >= 0 => t <= y0/adg
	tMax := y0 / adg

	if tMin > tMax {
		return result
	}

	for t := tMin; t <= tMax; t++ {
		X := x0 + bdg*t
		Y := y0 - adg*t
		if X >= 0 && Y >= 0 {
			result = append(result, Combination{X, Y})
		}
	}

	return result
}

// Returns the minimal cost if there is a solution, or -1 if no solutions.
func minCost(combinations []Combination) int64 {
	var minCost int64 = 1<<63 - 1 // max int64
	aCost := int64(3)
	bCost := int64(1)

	for _, comb := range combinations {
		cost := comb.a*aCost + comb.b*bCost
		if cost < minCost {
			minCost = cost
		}
	}
	if minCost == 1<<63-1 {
		return 0
	}
	return minCost
}

func answer1() int {
	totalCost := int64(0)
	for _, machine := range readInput() {
		combinations := buttonCombinations(machine)
		totalCost += minCost(combinations)
	}
	return int(totalCost)
}

// -----------------------------------------------------------------------

// PART 2
// now add 10000000000000 to the X and Y position of every prize and recalculate

func answer2() int {
	totalCost := int64(0)
	const offset = 10000000000000
	for _, machine := range readInput() {
		machine.prize.x += offset
		machine.prize.y += offset
		combinations := buttonCombinations(machine)
		totalCost += minCost(combinations)
	}
	return int(totalCost)
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 38714,
	//2: 0,
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
