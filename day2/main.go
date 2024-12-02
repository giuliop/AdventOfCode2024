package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// PART 1

func readInput() [][]int {
	file, err := os.Open("input/day2")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var reports [][]int

	for scanner.Scan() {
		numbersString := scanner.Text()
		levelStrings := strings.Split(numbersString, " ")
		levels := make([]int, len(levelStrings))
		for i, level := range levelStrings {
			levels[i], _ = strconv.Atoi(level)
		}
		reports = append(reports, levels)
	}
	return reports
}

func isSafe(report []int) bool {
	dir := report[0] < report[1]
	for i := 0; i < len(report)-1; i++ {
		diff := report[i] - report[i+1]
		if (diff < 0) != dir || diff == 0 || diff > 3 || diff < -3 {
			return false
		}
	}
	return true
}

func answer1() int {
	// each input line is like "7 6 4 2 1"
	// each line is a "report" and each number is a "level"
	// a report is "safe" if
	//  - the levels are either all increasing or all decreasing.
	//  - any two adjacent levels differ by at least one and at most three
	// Return how many reports are safe
	reports := readInput()
	sum := 0
	for _, report := range reports {
		if isSafe(report) {
			sum += 1
		}
	}
	return sum
}

// -----------------------------------------------------------------------

// PART 2

// we can tolerate one bad level in a report by removing it
func isSafeWithTolerance(report []int) bool {
	for skip := 0; skip < len(report); skip++ {
		newReport := make([]int, 0, len(report))
		newReport = append(newReport, report[:skip]...)
		newReport = append(newReport, report[skip+1:]...)
		if isSafe(newReport) {
			return true
		}
	}
	return false
}

func answer2() int {
	// a report is still safe if we remove one "bad" level and it becomes safe
	reports := readInput()
	sum := 0
	for _, report := range reports {
		if isSafeWithTolerance(report) {
			sum += 1
		}
	}
	return sum
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 224,
	2: 293,
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
