package main

import (
	"log"
	"os"
	"regexp"
	"strconv"
)

// PART 1

func readInput() string {
	data, err := os.ReadFile("input/day3")
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

func answer1() int {
	// input is a long string, with scattered substrings of the form "mul(x,y)"
	// return the sum of the products x*y for all substrings
	i := readInput()
	re := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	muls := re.FindAllStringSubmatch(i, -1)
	sum := 0
	for _, m := range muls {
		x, err := strconv.Atoi(m[1])
		if err != nil {
			log.Fatal(err)
		}
		y, err := strconv.Atoi(m[2])
		if err != nil {
			log.Fatal(err)
		}

		sum += x * y
	}
	return sum
}

// -----------------------------------------------------------------------

// PART 2

func answer2() int {
	// now we also have "do()" and "don't()" instructions that enable or disable the
	// following multiplication of the numbers in the following "mul()" instructions.
	// At the beginning, multiplication is enabled.
	i := readInput()
	re := regexp.MustCompile(`(do|don't)\(\)|(mul)\((\d+),(\d+)\)`)
	matches := re.FindAllStringSubmatch(i, -1)
	sum := 0
	mulEnabled := true
	for _, m := range matches {
		switch m[1] {
		case "do":
			mulEnabled = true
		case "don't":
			mulEnabled = false
		// if we matched a "mul(x,y)" m[1] will be "" and m[2] will be "mul"
		case "":
			if mulEnabled {
				x, err := strconv.Atoi(m[3])
				if err != nil {
					log.Fatal(err)
				}
				y, err := strconv.Atoi(m[4])
				if err != nil {
					log.Fatal(err)
				}
				sum += x * y
			}
		}
	}
	return sum
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 181345830,
	2: 98729041,
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
