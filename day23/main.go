package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// PART 1
// Your input is a list of two letters computer ids linked togeterh, e,g,
//     kh-tc
// Find and count all the triples of computer ids that are linked together
// where at least one starts with 't'

type Graph map[string][]string // computer id -> list of linked computer ids

func readInput() Graph {
	file, err := os.Open("input/day23")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	graph := Graph{}

	for scanner.Scan() {
		c1, c2 := scanner.Text()[:2], scanner.Text()[3:]
		graph[c1] = append(graph[c1], c2)
		graph[c2] = append(graph[c2], c1)
	}

	return graph
}

func answer1() int {
	graph := readInput()
	triples := map[string]bool{}
	for c1, linked := range graph {
		for _, c2 := range linked {
			for _, c3 := range graph[c2] {
				if c3 != c1 {
					if c1[0] == 't' || c2[0] == 't' || c3[0] == 't' {
						for _, c1Wanted := range graph[c3] {
							if c1Wanted == c1 {
								set := []string{c1, c2, c3}
								sort.Strings(set)
								triples[set[0]+set[1]+set[2]] = true
							}
						}
					}
				}
			}
		}
	}
	return len(triples)
}

// -----------------------------------------------------------------------

// PART 2
// now find the largest set of interconnected computer ids and return their names
// sorted alphabetically and separated by commas

func connectedToAll(pc string, network map[string]bool, graph Graph) bool {
	for c := range network {
		if c != pc && !contains(graph[c], pc) {
			return false
		}
	}
	return true
}

func contains(slice []string, item string) bool {
	for _, element := range slice {
		if element == item {
			return true
		}
	}
	return false
}

func answer2() int {
	graph := readInput()
	visited := map[string]bool{}
	maxNetworkSize := 0
	maxNetwork := map[string]bool{}
	for pc := range graph {
		network := map[string]bool{pc: true}
		toVisit := []string{pc}
		for len(toVisit) > 0 {
			c := toVisit[len(toVisit)-1]
			toVisit = toVisit[:len(toVisit)-1]
			if visited[c] {
				continue
			}
			visited[c] = true
			for _, linked := range graph[c] {
				if connectedToAll(linked, network, graph) {
					network[linked] = true
					toVisit = append(toVisit, linked)
				}
			}
		}
		if len(network) > maxNetworkSize {
			maxNetworkSize = len(network)
			maxNetwork = network
		}
	}
	keys := []string{}
	for k := range maxNetwork {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	fmt.Println(strings.Join(keys, ","))
	return maxNetworkSize
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 1368,
	// dd,ig,il,im,kb,kr,pe,ti,tv,vr,we,xu,zi
	2: 13,
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
