package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

// PART 1
// Your input is a list of numbers engraved on stones. At each blink of the eyes,
// stones' numbers change according to the first applicable rule:
// - If 0, becomes 1
// - If it has an even number of digits, the stone is replaced by two stones,
// 	 each with half the digit (the new numbers don't keep extra leading zeroes)
// - Otherwise, the number is multiplied by 2024
// How many stones are there after 25 blinks?

func readInput() []int {
	file, err := os.Open("input/day11")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// file is one line of numbers
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	line := scanner.Text()
	numbers := make([]int, 0)
	for _, n := range strings.Split(line, " ") {
		num, err := strconv.Atoi(n)
		if err != nil {
			log.Fatal(err)
		}
		numbers = append(numbers, num)
	}
	return numbers
}

// Rule takes a number and returns a slice of new
// numbers and a boolean indicating if the rule applied
type Rule func(int) ([]int, bool)

func rule1(n int) ([]int, bool) {
	if n == 0 {
		return []int{1}, true
	}
	return []int{n}, false
}

func rule2(n int) ([]int, bool) {
	nStr := strconv.Itoa(n)
	if len(nStr)%2 == 0 {
		half := len(nStr) / 2
		left, err2 := strconv.Atoi(nStr[half:])
		right, err1 := strconv.Atoi(nStr[:half])
		if err1 != nil || err2 != nil {
			log.Fatalf("Error converting string %s to two ints", nStr)

		}
		return []int{left, right}, true
	} else {
		return []int{n}, false
	}
}

func rule3(n int) ([]int, bool) {
	return []int{n * 2024}, true
}

func applyRules(n int, rules []Rule) []int {
	newStones := make([]int, 0)
	for _, rule := range rules {
		new, ok := rule(n)
		if ok {
			newStones = append(newStones, new...)
			break
		}
	}
	return newStones
}

func blink(stones []int, rules []Rule) []int {
	newStones := make([]int, 0)
	for _, s := range stones {
		newStones = append(newStones, applyRules(s, rules)...)
	}
	return newStones
}

// blinkNTimes applies the blink function n times
func blinkNTimes(stones []int, rules []Rule, n int) []int {
	for i := 0; i < n; i++ {
		stones = blink(stones, rules)
	}
	return stones
}

func answer1() int {
	rules := []Rule{rule1, rule2, rule3}
	stones := readInput()
	rounds := 25
	return len(blinkNTimes(stones, rules, rounds))
}

// -----------------------------------------------------------------------

// PART 2
// Now blink 75 times. How many stones are there now?
//
// We cannot run the naive solution for 75 rounds, it would be too slow.
// If we have the same number/rounds pair  multiple times, we can multiply the result
// by the number of times we have it

type DigitInfo struct {
	value      int
	rounds     int
	multiplier int
}

// WorkInProgress stores DigitInfos using value,rounds as key
type WorkInProgress map[[2]int]DigitInfo

func (wip WorkInProgress) pop() DigitInfo {
	var key [2]int
	for k := range wip {
		key = k
		break
	}
	d := wip[key]
	delete(wip, key)
	return d
}

func (wip WorkInProgress) add(d DigitInfo) {
	key := [2]int{d.value, d.rounds}
	if _, ok := wip[key]; ok {
		wip[key] = DigitInfo{d.value, d.rounds, wip[key].multiplier + d.multiplier}
	} else {
		wip[key] = d
	}
}

func answer2() int {
	rules := []Rule{rule1, rule2, rule3}
	stones := readInput()
	res := 0
	rounds := 75
	wip := make(WorkInProgress)
	for _, s := range stones {
		wip[[2]int{s, 0}] = DigitInfo{s, 0, 1}
	}
	for len(wip) > 0 {
		dInfo := wip.pop()
		if dInfo.rounds == rounds {
			res += dInfo.multiplier
			continue
		}
		nums := applyRules(dInfo.value, rules)
		for _, n := range nums {
			wip.add(DigitInfo{n, dInfo.rounds + 1, dInfo.multiplier})
		}
	}
	return res
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 199982,
	2: 237149922829154,
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
