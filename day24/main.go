package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unique"
)

// PART 1
// Your input is a list of wire names and initial value (0 or 1), e.g.
//     x00: 1
// followed by a blank line and a list of gates such as:
//    x14 AND y14 -> cwj
// indicating how the input wires generate the output wires.
// Gates can be AND, OR, XOR.
// Simulate the circuit and return the number represented by the bits generated
// by the wires starting with 'z', where z00 is the least significant bit.

type Op func(int, int) int

func and(a, b int) int {
	return a & b
}

func or(a, b int) int {
	return a | b
}

func xor(a, b int) int {
	return a ^ b
}

type Gate struct {
	op     Op
	opName string
	input1 unique.Handle[string]
	input2 unique.Handle[string]
	output unique.Handle[string]
}

func (g *Gate) compute(wires System) bool {
	in1 := wires[g.input1]
	in2 := wires[g.input2]
	out := wires[g.output]
	if in1.value == undefined || in2.value == undefined {
		out.value = undefined
		return false
	} else {
		out.value = Output(g.op(int(in1.value), int(in2.value)))
		return true
	}
}

type Output int

const (
	zero Output = iota
	one
	undefined
)

type Wire struct {
	name     unique.Handle[string]
	value    Output
	inputTo  []*Gate
	outputOf *Gate
}

type System map[unique.Handle[string]]*Wire

type NamesToOutput map[unique.Handle[string]]Output

func readInput() (s System, initializedWires NamesToOutput) {
	file, err := os.Open("input/day24")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	initializedWires = NamesToOutput{}
	s = System{}

	scanner := bufio.NewScanner(file)
	// read the initialized wires
	regex := regexp.MustCompile(`^(\w+): (\d+)$`)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // blank line, end of initialized wires
			break
		}
		matches := regex.FindStringSubmatch(line)
		if matches == nil {
			log.Fatal("Invalid input: ", line)
		}
		wireName := unique.Make(matches[1])
		valueString := matches[2]
		value, err := strconv.Atoi(valueString)
		if err != nil {
			log.Fatal("Invalid input: ", line)
		}
		initializedWires[wireName] = Output(value)
	}

	// read the gates
	regex = regexp.MustCompile(`^(\w+) (\w+) (\w+) -> (\w+)$`)
	for scanner.Scan() {
		line := scanner.Text()
		matches := regex.FindStringSubmatch(line)
		if matches == nil {
			log.Fatal("Invalid input: ", line)
		}
		in1 := unique.Make(matches[1])
		opName := matches[2]
		in2 := unique.Make(matches[3])
		out := unique.Make(matches[4])
		op := and
		if opName == "OR" {
			op = or
		} else if opName == "XOR" {
			op = xor
		}
		gate := Gate{op, opName, in1, in2, out}
		s.initializeGate(&gate, initializedWires)
	}

	return s, initializedWires
}

func (s System) initializeGate(gate *Gate, initializedWires NamesToOutput) {
	for _, wireName := range []unique.Handle[string]{gate.input1, gate.input2} {
		if _, ok := s[wireName]; !ok {
			value, ok := initializedWires[wireName]
			if !ok {
				value = undefined
			}
			s[wireName] = &Wire{wireName, value, []*Gate{gate}, nil}
		} else {
			s[wireName].inputTo = append(s[wireName].inputTo, gate)
		}
	}
	if _, ok := s[gate.output]; !ok {
		s[gate.output] = &Wire{gate.output, undefined, []*Gate{}, gate}
	} else {
		s[gate.output].outputOf = gate
	}
}

// value returns the integer value of all the bits with the given name, which can be used with
// input x, y, or z to get the value of the corresponding input or ouput.
func (s System) value(name byte) int {
	if name != 'x' && name != 'y' && name != 'z' {
		panic("Invalid byte " + string(name))
	}
	res := 0
	wires := []unique.Handle[string]{}
	for _, wire := range s {
		if wire.name.Value()[0] == name {
			wires = append(wires, wire.name)
		}
	}
	sort.Slice(wires, func(i, j int) bool {
		num1, _ := strconv.Atoi(strings.TrimPrefix(wires[i].Value(), string(name)))
		num2, _ := strconv.Atoi(strings.TrimPrefix(wires[j].Value(), string(name)))
		return num1 < num2
	})
	for i, wireName := range wires {
		res += int(s[wireName].value) << i
	}
	return res
}

func (s System) run(initializedWires NamesToOutput) int {
	newWires := []unique.Handle[string]{}
	for wireName := range initializedWires {
		newWires = append(newWires, wireName)
	}
	for len(newWires) > 0 {
		wire := newWires[len(newWires)-1]
		newWires = newWires[:len(newWires)-1]
		gates := s[wire].inputTo
		for _, gate := range gates {
			didCompute := gate.compute(s)
			if didCompute {
				newWires = append(newWires, gate.output)
			}
		}
	}
	return s.value('z')
}

func answer1() int {
	s, initializedWires := readInput()
	return s.run(initializedWires)
}

// -----------------------------------------------------------------------

// PART 2
// Now we need to fix the system so that if we interpret the x inputs as one binary numbers
// and the y inputs as another binary number, the output z is the sum of the two numbers.
// For x, y, z as before we have that x00, y00, z00 are the least significant bits and so on.
// To fix the systems there are four pair of gates that need their output wires swapped.
// Find them and return the list of the eight output wires that need to be swapped in
// lexicographical order and separated by commas.

// Generating all the possible swap combinations is too large.
// We use the 'displaySubCircuit' function to examine how each output wire is connected.
// We can see that the adder is built like this:
//    z[n] = XOR
//              XOR(x[n], y[n])
//              carry[n-1] )
//    carry[n-1] = OR
//                   AND(x[n-1], y[n-1])
//                   AND
//						carry[n-2]
//                      XOR(x[n-1], y[n-1])
// We'll trasvere the circuit and rename the output wires of the carries to e.g., c00.
// Then we'll inspect the circuit to easily find the wires that need to be swapped.

func (s System) swap(w1, w2 *Wire) {
	gate1, gate2 := w1.outputOf, w2.outputOf
	gate1.output, gate2.output = w2.name, w1.name
	w1.outputOf, w2.outputOf = w2.outputOf, w1.outputOf
}

// displaySubCircuit displays the subcircuit of the wire out, showing how all inputs
// / wires combine to generate the output wire
func (s System) displaySubCircuit(out *Wire, levels int) {
	fmt.Println(out.name.Value())
	var helper func(w *Wire, indent string, levels int)
	helper = func(w *Wire, indent string, levels int) {
		if levels == 0 {
			return
		}
		gate := w.outputOf
		if gate == nil {
			return
		}
		indentNext := indent + "   " // 3 spaces per level
		in1 := s[gate.input1]
		in2 := s[gate.input2]

		fmt.Println(indentNext + in1.name.Value())
		if in1.name.Value()[:3] != "car" && in1.name.Value()[:3] != "xor" && in1.name.Value()[:3] != "and" && in1.name.Value()[:3] != "cnd" {
			helper(in1, indentNext, levels-1)
		}

		fmt.Println(indentNext + " " + gate.opName)

		fmt.Println(indentNext + in2.name.Value())
		if in2.name.Value()[:3] != "car" && in2.name.Value()[:3] != "xor" && in2.name.Value()[:3] != "and" && in2.name.Value()[:3] != "cnd" {
			helper(in2, indentNext, levels-1)
		}
	}
	helper(out, "", levels)
}

func (s System) rename(w *Wire, newName unique.Handle[string]) error {
	if w.name.Value()[0] == 'z' {
		return fmt.Errorf("cannot rename z wires: %s", w.name.Value())
	}
	if w.outputOf == nil {
		return fmt.Errorf("cannot rename wires that are not outputs of gates: %s", w.name.Value())
	}
	if len(w.name.Value()) > 3 {
		return fmt.Errorf("cannot rename wires that are not named with 3 characters: %s",
			w.name.Value())
	}
	oldName := w.name
	w.name = newName
	delete(s, oldName)
	s[newName] = w
	w.outputOf.output = newName
	for _, gate := range w.inputTo {
		if gate.input1 == oldName {
			gate.input1 = newName
		}
		if gate.input2 == oldName {
			gate.input2 = newName
		}
	}
	return nil
}

// wiresNamed returns the wires with the given prefix name, sorted by the number in the name
func (s System) wiresNamed(prefix string) []*Wire {
	wires := []*Wire{}
	for _, wire := range s {
		if strings.HasPrefix(wire.name.Value(), prefix) {
			wires = append(wires, wire)
		}
	}
	sort.Slice(wires, func(i, j int) bool {
		num1, _ := strconv.Atoi(strings.TrimPrefix(wires[i].name.Value(), string(prefix)))
		num2, _ := strconv.Atoi(strings.TrimPrefix(wires[j].name.Value(), string(prefix)))
		return num1 < num2
	})
	return wires
}

// gates returns all the gates in the system
func (s System) gates() []*Gate {
	gates := []*Gate{}
	for _, wire := range s {
		if wire.outputOf != nil {
			gates = append(gates, wire.outputOf)
		}
	}
	return gates
}

func (w *Wire) number() int {
	for i, c := range w.name.Value() {
		if c >= '0' && c <= '9' {
			num, _ := strconv.Atoi(w.name.Value()[i:])
			return num
		}
	}
	panic("Invalid wire name " + w.name.Value())
}

func numberToWireNumber(num int) string {
	numString := strconv.Itoa(num)
	if num < 10 {
		numString = "0" + numString
	}
	return numString
}

func (s System) labelCarriesFrom(prevCarry *Wire) error {
	for _, andGate := range prevCarry.inputTo {
		otherInput := andGate.input1
		if otherInput == prevCarry.name {
			otherInput = andGate.input2
		}
		xorGate := s[otherInput].outputOf
		if xorGate == nil {
			return fmt.Errorf("could not find XOR gate for %s", prevCarry.name.Value())
		}
		num := numberToWireNumber(prevCarry.number() + 1)
		if andGate.opName == "AND" &&
			strings.HasPrefix(xorGate.output.Value(), "xor") {
			err := s.rename(s[andGate.output], unique.Make("cnd"+num))
			if err != nil {
				return fmt.Errorf("could not rename %s: %v - %s", andGate.output.Value(),
					err, prevCarry.name.Value())
			}
			for _, orGate := range s[andGate.output].inputTo {
				otherInput := orGate.input1
				if otherInput == s[andGate.output].name {
					otherInput = orGate.input2
				}
				if orGate.opName == "OR" &&
					strings.HasPrefix(s[otherInput].name.Value(), "and") {
					err := s.rename(s[orGate.output], unique.Make("carry"+num))
					if err != nil {
						return fmt.Errorf("could not rename %s: %v - %s",
							orGate.output.Value(), err, prevCarry.name.Value())
					}
					return s.labelCarriesFrom(s[orGate.output])
				}
			}
		}
	}
	return fmt.Errorf("could not find the next carry after %s", prevCarry.name.Value())
}

// hasInputBits returns true if the gate inputs are the same x and y bit (e.g., x01 and y01)
func (g *Gate) hasInputBits() bool {
	return ((g.input1.Value()[0] == 'x' && g.input2.Value()[0] == 'y') ||
		(g.input1.Value()[0] == 'y' && g.input2.Value()[0] == 'x')) &&
		(g.input1.Value()[1:] == g.input2.Value()[1:])
}

// renameWires renames the input wires to z00, z01, ... into c00, ... for the carries and
// s00, ... for the sum bits. We identify the sum and carries by looking at the z-wire gate.
// If it's well formed the gate should be a XOR with inputs that should be a XOR gate (sum)
// and an OR gate (carry)
func (s System) renameWires() string {
	// identify XOR (x XOR y) and AND (x AND y ) of input bits and rename them with xor, and
	for _, gate := range s.gates() {
		// don't rename z wires
		if strings.HasPrefix(gate.output.Value(), "z") {
			continue
		}
		if gate.opName == "XOR" && gate.hasInputBits() {
			s.rename(s[gate.output], unique.Make("xor"+gate.input1.Value()[1:]))
		} else if gate.opName == "AND" && gate.hasInputBits() {
			s.rename(s[gate.output], unique.Make("and"+gate.input1.Value()[1:]))
		}
	}
	// the first carry is the output of the first AND gate
	firstCarry := (s[unique.Make("and00")])
	err := s.labelCarriesFrom(firstCarry)
	if err != nil {
		fmt.Println(err)
	}
	// could not find the next carry after carry04
	// examining z05 and z06 we find the first swap: z05 and dkr
	res := []string{"z05", "dkr"}
	s.swap(s[unique.Make("z05")], s[unique.Make("dkr")])
	s.rename(s[unique.Make("bvc")], unique.Make("carry05"))
	err = s.labelCarriesFrom(s[unique.Make("carry05")])
	if err != nil {
		fmt.Println(err)
	}
	// next up carry15; we need to swap z15 and htp
	res = append(res, "z15", "htp")
	s.swap(s[unique.Make("z15")], s[unique.Make("htp")])
	err = s.labelCarriesFrom(s[unique.Make("carry14")])
	if err != nil {
		fmt.Println(err)
	}
	// next up carry 20; we need to swap z20 with hhh
	res = append(res, "z20", "hhh")
	s.swap(s[unique.Make("z20")], s[unique.Make("hhh")])
	s.rename(s[unique.Make("hhh")], unique.Make("carry20"))
	err = s.labelCarriesFrom(s[unique.Make("carry20")])
	if err != nil {
		fmt.Println(err)
	}
	// next up carry 35; we need to swap xor36(rhv) with and36(ggk)
	res = append(res, "rhv", "ggk")
	s.swap(s[unique.Make("xor36")], s[unique.Make("and36")])
	s.rename(s[unique.Make("gqf")], unique.Make("carry36"))
	err = s.labelCarriesFrom(s[unique.Make("carry36")])
	if err != nil {
		fmt.Println(err)
	}
	// wires := s.wiresNamed("z")
	// for _, wire := range wires {
	// 	s.displaySubCircuit(wire, 3)
	// }
	sort.Strings(res)
	return strings.Join(res, ",")
}

func answer2() int {
	s, _ := readInput()
	fmt.Println(s.renameWires())

	return 0
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 53755311654662,
	// dkr,ggk,hhh,htp,rhv,z05,z15,z20
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
