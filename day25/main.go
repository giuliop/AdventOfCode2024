package main

import (
	"bufio"
	"log"
	"os"
)

// PART 1
// You have keys and locks in your input file, like this (lock left, key right):
// #####       .....
// .####       #....
// .####       #....
// .####       #...#
// .#.#.       #.#.#
// .#...       #.###
// .....       #####
// Locks are filled in the top row, keys in the bottom row.
// Count the key/lock pairs that don't overlap with each other (e.g., in the example
// above, the key and lock overlap in the last column, so that pair doesn't count).

// each int is the number of filled cells in a column
type Lock []int
type Key []int

const (
	lock int = 0
	key  int = 1
)

func readInput() (locks []Lock, keys []Key) {
	file, err := os.Open("input/day25")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	schemaColumns := 5
	schemaRows := 7
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		schemaType := key
		if line == "#####" {
			schemaType = lock
		}
		schematic := make([]int, schemaColumns)
		for i := 0; i < schemaRows; i++ {
			for j := 0; j < schemaColumns; j++ {
				if line[j] == '#' {
					schematic[j]++
				}
			}
			scanner.Scan()
			line = scanner.Text()
		}
		if schemaType == lock {
			locks = append(locks, Lock(schematic))
		} else {
			keys = append(keys, Key(schematic))
		}
	}
	return locks, keys
}

func overlap(lock Lock, key Key) bool {
	for i := 0; i < len(lock); i++ {
		if lock[i]+key[i] > 7 {
			return true
		}
	}
	return false
}

func answer1() int {
	locks, keys := readInput()
	res := 0
	for _, lock := range locks {
		for _, key := range keys {
			if !overlap(lock, key) {
				res++
			}
		}
	}
	return res
}

// -----------------------------------------------------------------------

// PART 2

func answer2() int {
	return 0
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 3619,
	2: 0,
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
