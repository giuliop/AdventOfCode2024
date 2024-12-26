package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// PART 1
// You have a list of numbers, one per line of input file. Each number is a seed for a
// pseudo-random number generator. At each round the new number is generated by:
// - multiplying the previous number by 64
// - mix the result with the previous number (by bitwise XOR)
// - prune the result (by taking modul0 16777216)
// - divide the result by 32, mix and prune again
// - multiply the result by 2048, mix and prune again
// Calculate the 2000th secret number for each seed and return the sum of all of them.

func readInput() []int {
	file, err := os.Open("input/day22")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	seeds := []int{}
	for scanner.Scan() {
		seedString := scanner.Text()
		seed, err := strconv.Atoi(seedString)
		if err != nil {
			log.Fatal(err)
		}
		seeds = append(seeds, seed)
	}

	return seeds
}

// Note that 64->2^6, 32->2^5, 2048->2^11, 16777216->2^24
// In bit terms, the process looks like this:
// 1. XOR with self shifted 6 bits to the left and take lowest 24 bits
// 2. XOR with self shifted 5 bits to the right and take lowest 24 bits
// 3. XOR with self shifted 11 bits to the left and take lowest 24 bits

func nextSecret(secret int) int {
	const MODULO = 16777216 // 2^24
	secret = (secret * 64) ^ secret
	secret = secret & (MODULO - 1)
	secret = (secret / 32) ^ secret
	secret = secret & (MODULO - 1)
	secret = (secret * 2048) ^ secret
	secret = secret & (MODULO - 1)
	return secret
}

func answer1() int {
	sum := 0
	seeds := readInput()
	for _, seed := range seeds {
		for i := 0; i < 2000; i++ {
			seed = nextSecret(seed)
		}
		sum += seed
	}
	return sum
}

// -----------------------------------------------------------------------

// PART 2
// Now consider the ones digit of each of the 2000 secret numbers for each buyer.
// Consider the differences between the ones digits like this:
// 123: 3
// 15887950: 0 (-3)
// 16495136: 6 (6)
//   527345: 5 (-1)
//   704524: 4 (-1
// ...
// The first time a consecutive series of four differences is found, the corresponding
// number is the numnber of bananas the buyer will pay if given that sequance.
// Determine the sequnce of differences that will result in the highest number of bananas
// from all buyers.

func addSequencesToMemory(nums []int, memory map[[4]int]int) {
	ones := []int{}
	found := map[[4]int]bool{}
	for _, n := range nums {
		ones = append(ones, n%10)
	}
	for i := 1; i < len(ones)-3; i++ {
		key := [4]int{ones[i] - ones[i-1], ones[i+1] - ones[i],
			ones[i+2] - ones[i+1], ones[i+3] - ones[i+2]}
		if !found[key] {
			memory[key] += ones[i+3]
			found[key] = true
		}
	}
}

func answer2() int {
	seeds := readInput()
	buyersNums := [][]int{}
	for _, seed := range seeds {
		nums := []int{seed}
		for i := 0; i < 2000; i++ {
			seed = nextSecret(seed)
			nums = append(nums, seed)
		}
		buyersNums = append(buyersNums, nums)
	}
	memory := map[[4]int]int{}
	for _, nums := range buyersNums {
		addSequencesToMemory(nums, memory)
	}
	max := 0
	for _, v := range memory {
		if v > max {
			max = v
		}
	}
	return max
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 13429191512,
	2: 1582,
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