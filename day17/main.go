package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// PART 1
// Run a computer described below and return its outputs concatenated with commas as a string
// The computer has three registers, A, B, C which can hold any integer.
// It accepts intructions in the form of <opcode operand>.
// The instruction pointer starts at 0 and increases by 2 after each instruction.
//
// Operands, depending on the opcode, can be a literal integer or a "combo" operand, tha is:
// - Combo operands 0 through 3 represent literal values 0 through 3.
// - Combo operand 4 represents the value of register A.
// - Combo operand 5 represents the value of register B.
// - Combo operand 6 represents the value of register C.
// - Combo operand 7 is reserved and will not appear in valid programs.
//
// Opcodes are:
//  0. adv - performs division. Numerator: in A. Denominator: 2 to the power of combo operand.
//     The result is truncated to an integer and then written to the A register.
//  1. bxl - bitwise XOR of register B and the instruction's literal operand, result in B.
//  2. bst - combo operand modulo 8 (thereby keeping only its lowest 3 bits), result in B .
//  3. jnz - nothing if the A register is 0. If the A register is not zero, jump by setting
//     the instruction pointer to the value of its literal operand (no incrementing by 2)
//  4. bxc - bitwise XOR of register B and register C, then result in B (operand ignored)
//  5. out - combo operand modulo 8, then outputs value (if multiple values, separated by commas)
//  6. bdv - like adv but result stored in B register (numerator still read from A)
//  7. cdv - like adv but result stored in C register (numerator still read from A)
type InitialState struct {
	A, B, C int
	program []int
}

func readInput() InitialState {
	// input file format:
	// Register A: 47006051
	// Register B: 0
	// Register C: 0
	// (empty line)
	// Program: 2,4,1,3,7,5,1,5,0,3,4,3,5,5,3,0
	file, err := os.Open("input/day17")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var is InitialState
	scanner := bufio.NewScanner(file)

	for _, register := range []*int{&is.A, &is.B, &is.C} {
		scanner.Scan()
		_, regText, ok := strings.Cut(scanner.Text(), ": ")
		if !ok {
			log.Fatal("Invalid register")
		}
		*register, err = strconv.Atoi(regText)
		if !ok || err != nil {
			log.Fatal("Invalid register value")
		}
	}
	scanner.Scan()
	scanner.Scan()
	_, programText, ok := strings.Cut(scanner.Text(), ": ")
	if !ok {
		log.Fatal("Invalid program")
	}
	for _, instr := range strings.Split(programText, ",") {
		i, err := strconv.Atoi(instr)
		if err != nil {
			log.Fatal("Invalid instruction")
		}
		is.program = append(is.program, i)
	}
	return is
}

type Computer struct {
	A, B, C int
	program []int
	pc      int
	output  []int
}

func initializeComputer(is InitialState) Computer {
	return Computer{
		A:       is.A,
		B:       is.B,
		C:       is.C,
		program: is.program,
		pc:      0,
		output:  []int{},
	}
}

func (c *Computer) combo(operand int) int {
	if operand < 4 {
		return operand
	}
	switch operand {
	case 4:
		return c.A
	case 5:
		return c.B
	case 6:
		return c.C
	}
	log.Fatal("Invalid combo operand")
	return 0
}

func (c *Computer) step() {
	instr := c.program[c.pc]
	operand := c.program[c.pc+1]
	c.pc += 2
	switch instr {
	case 0:
		c.A /= 1 << c.combo(operand)
	case 1:
		c.B ^= operand
	case 2:
		c.B = c.combo(operand) % 8
	case 3:
		if c.A != 0 {
			c.pc = operand
		}
	case 4:
		c.B ^= c.C
	case 5:
		c.output = append(c.output, c.combo(operand)%8)
	case 6:
		c.B = c.A / (1 << c.combo(operand))
	case 7:
		c.C = c.A / (1 << c.combo(operand))
	default:
		log.Fatal("Invalid instruction")
	}
}

func (c *Computer) run() string {
	for c.pc < len(c.program) {
		c.step()
	}
	return c.getOutput()
}

func (c *Computer) getOutput() string {
	if len(c.output) == 0 {
		return ""
	}
	strOutputs := make([]string, len(c.output))
	for i, o := range c.output {
		strOutputs[i] = strconv.Itoa(o)
	}
	return strings.Join(strOutputs, ",")
}

func answer1() int {
	initialState := readInput()
	c := initializeComputer(initialState)
	output := c.run()
	fmt.Println(output)
	// return output as an int removing all commas
	outputInt, err := strconv.Atoi(strings.ReplaceAll(output, ",", ""))
	if err != nil {
		log.Fatal("Invalid output")
	}
	return outputInt
}

// -----------------------------------------------------------------------

// PART 2
// Now find the minimum value for register A in the initial state that makes the computer
// output  the same initial program.

// Let's analyze the program:
// 2,4 1,3 7,5 1,5 0,3 4,3 5,5 3,0
// 0 - 2,4 - bst A -> B = A % 8
// 1 - 1,3 - bxl 3 -> B ^= 3
// 2 - 7,5 - cdv 5 -> C = A / 2 **  B
// 3 - 1,5 - bxl 5 -> B ^= 5
// 4 - 0,3 - adv 3 -> A /= 2 ** 3
// 5 - 4,3 - bxc   -> B ^= C
// 6 - 5,5 - out 5 -> output B % 8
// 7 - 3,0 - jnz 0 -> if A != 0 jump to 0

// The program is equivalent to:
// output ((((A % 8) ^ 3) ^ 5) ^ (A / 2 ** ((A % 8) ^ 3))) % 8
// A = A / 8
// if A != 0 jump to 0

// At the last iteration A must be less than 8 and the output expression above must
// be equal to the program's last value. I can use that to find the last value of A.
// Then at the run before that, A shall be so that A / 8 gives me the last value
// of A and so on until I find the initial value of A

func outputForA(A int) int {
	return ((((A % 8) ^ 3) ^ 5) ^ (A/(1<<((A%8)^3)))%8)
}

func minInt(slice []int) int {
	if len(slice) == 0 {
		log.Fatal("Cannot find the minimum of an empty slice")
	}
	min := slice[0]
	for _, val := range slice[1:] {
		if val < min {
			min = val
		}
	}
	return min
}

func answer2() int {
	initialState := readInput()
	lenProgram := len(initialState.program)
	As := []int{0}
	for i := lenProgram - 1; i >= 0; i-- {
		programValue := initialState.program[i]
		newAs := []int{}
		for _, A := range As {
			for remainder := 0; remainder < 8; remainder++ {
				candidate := A*8 + remainder
				if outputForA(candidate) == programValue {
					newAs = append(newAs, candidate)
				}
			}
			As = newAs
		}
	}
	if len(As) == 0 {
		log.Fatal("No solution found")
	}
	return minInt(As)
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 627231605,
	2: 236548287712877,
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
