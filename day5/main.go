package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// PART 1

type Rules map[string][]string
type Update []string

func readInput() (Rules, []Update) {
	file, err := os.Open("input/day5")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	rules := make(Rules)
	var updates []Update

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		if x, y, ok := strings.Cut(line, "|"); !ok {
			log.Fatal("Invalid rule: ", line)
		} else {
			rules[x] = append(rules[x], y)
		}
	}
	for scanner.Scan() {
		line := scanner.Text()
		updates = append(updates, strings.Split(line, ","))
	}

	return rules, updates
}

func middleValue(u Update) int {
	if len(u)%2 == 0 {
		log.Fatal("Update has even length: ", u)
	}
	value, err := strconv.Atoi(u[len(u)/2])
	if err != nil {
		log.Fatal("Invalid middle value: ", u)
	}
	return value
}

func isValid(update Update, rules Rules) bool {
	for i, elem := range update[1:] {
		for _, prev := range update[:i] {
			if notBefore, ok := rules[elem]; ok {
				for _, invalidElem := range notBefore {
					if prev == invalidElem {
						return false
					}
				}
			}
		}
	}
	return true
}

func answer1() int {
	// input is like this:
	// 81|51
	// ...
	//
	// 79,64,35,74,22,94,19
	// ...
	// The first section is a list of pairs of page numbers X|Y, specifyng that
	// page X must come before page Y. After a blank line, there is a list of
	// page updates, where each line is a list of page numbers separated by commas.
	// An update is valid if its page order respects the rules of the first section.
	// Sum the middle page numbers of all valid updates.
	rules, updates := readInput()
	sum := 0
	for _, u := range updates {
		if isValid(u, rules) {
			sum += middleValue(u)
		}
	}
	return sum
}

// -----------------------------------------------------------------------

// PART 2

// reorderElement moves the element at posToMove to posDest in the update,
// shifting right the elements between posDest and posToMove.
func reorderElement(update Update, posToMove, posDest int) {
	if posToMove <= posDest {
		log.Fatal("posToMove must be greater than posDest")
	}
	elem := update[posToMove]
	copy(update[posDest+1:posToMove+1], update[posDest:posToMove])
	update[posDest] = elem
}

func reorder(update Update, rules Rules) (Update, bool) {
	reordered := false
	for pos := 1; pos < len(update); pos++ {
		elem := update[pos]
	outer:
		for prevPos, prev := range update[:pos] {
			if rule, exist := rules[elem]; exist {
				for _, notBefore := range rule {
					if prev == notBefore {
						reorderElement(update, pos, prevPos)
						reordered = true
						break outer
					}
				}
			}
		}
	}
	return update, reordered
}

func answer2() int {
	// now reorder all the invalid updates so that they become valid and sum their middle values
	rules, updates := readInput()
	sum := 0
	for _, u := range updates {
		if newU, reordered := reorder(u, rules); reordered {
			sum += middleValue(newU)
		}
	}
	return sum
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 4637,
	2: 6370,
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
