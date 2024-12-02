package main

import (
	"bufio"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

// PART 1

func readInput() (leftNumbers, rightNumbers []int) {
	file, err := os.Open("input/day1")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		numStrings := strings.Fields(scanner.Text())
		l, _ := strconv.Atoi(numStrings[0])
		r, _ := strconv.Atoi(numStrings[1])
		leftNumbers = append(leftNumbers, l)
		rightNumbers = append(rightNumbers, r)
	}
	return leftNumbers, rightNumbers
}

func answer1() int {
	// each line of input is a string like this: "69214   60950"
	// we need to order the leftNumbers and rightNumbers numbers in each line
	// and add all the differences between the rightNumbers and leftNumbers numbers
	leftNumbers, rightNumbers := readInput()
	slices.Sort(leftNumbers)
	slices.Sort(rightNumbers)
	sum := 0
	for i := 0; i < len(leftNumbers); i++ {
		diff := leftNumbers[i] - rightNumbers[i]
		if diff < 0 {
			diff = -diff
		}
		sum += diff
	}
	return sum
}

// -----------------------------------------------------------------------

// PART 2

func answer2() int {
	// compute the number of times each left number appears among the right numbers
	// sum all the left numbers times the number of times they appear among the right numbers
	leftNumbers, rightNumbers := readInput()
	m := make(map[int]int)
	for _, r := range rightNumbers {
		m[r]++
	}
	sum := 0
	for _, l := range leftNumbers {
		sum += l * m[l]
	}
	return sum
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 1879048,
	2: 21024792,
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
