package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// PART 1
// input is a list like this:
// 190: 10 19
// 3267: 81 40 27
// ...
// where each line starts with a result, followed by a list of numbers with missing
// operators in between. The operators can be +, * and are applied from left to right.
// Results are valid if they can be obtained from the numbers in the list.
// Sum all valid results.

type Equation struct {
	result  int
	numbers []int
}

func readInput() []Equation {
	file, err := os.Open("input/day7")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var equations []Equation

	for scanner.Scan() {
		resStr, numbersStr, _ := strings.Cut(scanner.Text(), ": ")
		result, err := strconv.Atoi(resStr)
		if err != nil {
			log.Fatal(err)
		}
		numbers := make([]int, 0)
		for _, n := range strings.Split(numbersStr, " ") {
			num, err := strconv.Atoi(n)
			if err != nil {
				log.Fatal(err)
			}
			numbers = append(numbers, num)
		}
		equations = append(equations, Equation{result, numbers})
	}
	return equations
}

type Combination struct {
	partialResult int
	numbers       []int
}

type Op func(int, int) int

func addOP(a, b int) int  { return a + b }
func multOP(a, b int) int { return a * b }

func isValid(eq Equation, ops []Op) bool {
	if len(eq.numbers) == 1 {
		return eq.result == eq.numbers[0]
	}
	combinations := []Combination{{eq.numbers[0], eq.numbers[1:]}}
	for {
		if len(combinations) == 0 {
			break
		}
		comb := combinations[len(combinations)-1]
		combinations = combinations[:len(combinations)-1]
		if comb.partialResult > eq.result {
			continue
		}
		if len(comb.numbers) == 0 {
			if comb.partialResult == eq.result {
				return true
			}
			continue
		}
		n := comb.numbers[0]
		for _, op := range ops {
			combinations = append(combinations,
				Combination{op(comb.partialResult, n), comb.numbers[1:]})
		}
	}
	return false
}

func answer1() int {
	equations := readInput()
	res := 0
	ops := []Op{addOP, multOP}
	for _, eq := range equations {
		if isValid(eq, ops) {
			res += eq.result
		}
	}
	return res
}

// -----------------------------------------------------------------------

// PART 2
// Now there is a third operator which concatenates the numbers.
// Sum all valid results.

func concatOP(a, b int) int {
	aStr := strconv.Itoa(a)
	bStr := strconv.Itoa(b)
	concat, err := strconv.Atoi(aStr + bStr)
	if err != nil {
		log.Fatal(err)
	}
	return concat
}

func answer2() int {
	ops := []Op{addOP, multOP, concatOP}
	equations := readInput()
	res := 0
	for _, eq := range equations {
		if isValid(eq, ops) {
			res += eq.result
		}
	}
	return res
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 2437272016585,
	2: 162987117690649,
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
