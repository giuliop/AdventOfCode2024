package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// PART 1
// Input is like:
// r, wr, b, g, bwu, rb, gb, br

// brwrr
// bggr
// ...

// The first line is a list of available partterns, the following lines are the
// design to form. Count how many designs can be formed from the available
// patterns (whihc can be used multiple times).

func readInput() (patterns []string, designs []string) {
	file, err := os.Open("input/day19")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	patterns = strings.Split(scanner.Text(), ", ")

	scanner.Scan()

	for scanner.Scan() {
		designs = append(designs, scanner.Text())
	}

	return
}

func canMake(d string, patterns []string, mem Memory) bool {
	if d == "" {
		return true
	}
	if can, ok := mem[d]; ok {
		return can
	}
	for _, p := range patterns {
		if strings.HasPrefix(d, p) {
			suffix := d[len(p):]
			if canMake(suffix, patterns, mem) {
				mem[d] = true
				return true
			}
		}
	}

	mem[d] = false
	return false
}

type Memory map[string]bool // design -> can make

func answer1() int {
	patterns, designs := readInput()
	mem := Memory{}
	res := 0
	for _, d := range designs {
		if canMake(d, patterns, mem) {
			res++
		}
	}
	return res
}

// -----------------------------------------------------------------------

// PART 2
// Now return the sum of all possible ways to form the designs.

type Memory2 map[string]int // design -> number of ways

func countCombinations(d string, patterns []string, mem Memory2) int {
	if d == "" {
		return 1
	}
	if count, ok := mem[d]; ok {
		return count
	}
	count := 0
	for _, p := range patterns {
		if strings.HasPrefix(d, p) {
			suffix := d[len(p):]
			suffixCount := countCombinations(suffix, patterns, mem)
			count += suffixCount
		}
	}
	mem[d] = count
	return count
}

func answer2() int {
	patterns, designs := readInput()
	mem := Memory2{}
	res := 0
	for _, d := range designs {
		res += countCombinations(d, patterns, mem)
	}
	return res
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 287,
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
