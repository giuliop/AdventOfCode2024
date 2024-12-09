package main

import (
	"io"
	"log"
	"os"
	"strconv"
)

// PART 1
// The input is a string of numbers like "2333133121414131402" where each odd digit
// represents a file length in blocks and each even digit represents empty blocks.
// Each file has an ID number, which starts with 0 for the leftmost file (the first
// digit) and increases by 1 for each file.
// The task is to first compact the disk by moving file blocks starting from the far
// right to empty blocks on the left. Then, calculate the checksum of the disk by
// multiplying for each file block its ID with its position in the disk and summing
// all these products.

const empty = -1

// Blocks is a struct to represent `Lenght` contiguous blocks of file `ID` on disk.
// If `ID` is `empty`, it represents `Length` empty blocks.
type Blocks struct {
	id     int
	length int
	before *Blocks
	after  *Blocks
}

type Disk struct {
	first      *Blocks
	last       *Blocks
	firstEmpty *Blocks
}

// mergeEmpty merges b (if it's an empty block) with the next/previous empty block
// if they exist.
func (d *Disk) mergeEmpty(b *Blocks) {
	if !b.isEmpty() {
		return
	}
	if b.before != nil && b.before.isEmpty() {
		b.length += b.before.length
		d.remove(b.before)
	}
	if b.after != nil && b.after.isEmpty() {
		b.length += b.after.length
		d.remove(b.after)
	}
}

// append a block to the disk
func (d *Disk) append(b *Blocks) {
	if d.first == nil {
		d.first = b
		d.last = b
	} else {
		d.last.after = b
		b.before = d.last
		d.last = b
	}
	if b.isEmpty() {
		d.mergeEmpty(b)
		if d.firstEmpty == nil {
			d.firstEmpty = b
		}
	}
}

// nextEmpty returns the next empty block after `b` in the disk
// or nil if there is none before the until block.
func (b *Blocks) nextEmpty(until *Blocks) *Blocks {
	if b == until {
		return nil
	}
	for b = b.after; b != until && b != nil; b = b.after {
		if b.isEmpty() {
			return b
		}
	}
	return nil
}

// isEmpty returns true if `b` is an empty block.
func (b *Blocks) isEmpty() bool {
	return b.id == empty
}

// insertBefore inserts `b` before `beforeThis` in the disk.
func (d *Disk) insertBefore(b *Blocks, beforeThis *Blocks) {
	if beforeThis == b.after {
		return
	}
	if beforeThis == nil {
		// only valid if disk is empty
		if d.first != nil {
			log.Fatal("insertBefore: beforeThis is nil but disk is not empty")
		}
		d.append(b)
		return
	}
	if d.last == b {
		d.last = b.before
	}
	if d.first == b {
		d.first = b.after
	}
	if beforeThis.before != nil {
		beforeThis.before.after = b
	}
	b.before = beforeThis.before
	b.after = beforeThis
	beforeThis.before = b
	if d.first == beforeThis {
		d.first = b
	}
	if b.isEmpty() {
		d.mergeEmpty(b)
		// we need to recompute firstEmpty
		d.firstEmpty = d.first.nextEmpty(nil)
	}
}

// remove a block from the disk
func (d *Disk) remove(b *Blocks) {
	if d.firstEmpty == b {
		d.firstEmpty = b.nextEmpty(nil)
	}
	if b.before != nil {
		b.before.after = b.after
	}
	if b.after != nil {
		b.after.before = b.before
	}
	if d.first == b {
		d.first = b.after
	}
	if d.last == b {
		d.last = b.before
	}
	if b.before != nil && b.after != nil {
		d.mergeEmpty(b.before)
	}
}

// isCompacted returns true if the disk has no empty space (except at the end)
func (d *Disk) isCompacted() bool {
	return d.firstEmpty == nil || d.firstEmpty == d.last
}

// checkSum returns the checksum of the disk.
func (d *Disk) checkSum() int {
	checksum := 0
	position := 0
	for b := d.first; b != nil; b = b.after {
		if b.isEmpty() {
			position += b.length
			continue
		}
		for j := 0; j < b.length; j++ {
			checksum += b.id * position
			position++
		}
	}
	return checksum
}

// moveBlocks moves blocks (fully or in part) from `from` to `to` overwriting 'to'.
// Doesn't move more that 'to' can hold
func (d *Disk) moveBlocks(from, to *Blocks) {
	if from == to {
		return
	}
	if to.length >= from.length {
		before := from.before
		after := from.after
		if after != nil {
			if after.isEmpty() {
				after.length += from.length
			} else {
				d.insertBefore(&Blocks{id: empty, length: from.length}, after)
			}
		}
		d.remove(from)
		d.insertBefore(from, to)
		if before != nil && before.isEmpty() {
			d.mergeEmpty(before)
		}
		to.length -= from.length
		if to.length == 0 {
			before := to.before
			d.remove(to)
			if before != nil && before.isEmpty() {
				d.mergeEmpty(before)
			}
		}
	} else {
		from.length -= to.length
		d.insertBefore(&Blocks{id: from.id, length: to.length}, to)
		d.remove(to)
		if from.after != nil {
			if from.after.isEmpty() {
				from.after.length += to.length
			} else {
				d.insertBefore(&Blocks{id: empty, length: to.length}, from.after)
			}
		}
	}
}

// compactWithFragmentation compacts the disk allowing fragmentation of same blocks.
func (d *Disk) compactWithFragmentation() {
	for !d.isCompacted() {
		for d.last.isEmpty() {
			d.remove(d.last)
		}
		d.moveBlocks(d.last, d.firstEmpty)
	}
}

func readInput() Disk {
	file, err := os.Open("input/day9")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	d := Disk{}

	isSpace := false
	for id := 0; ; isSpace = !isSpace {
		b := make([]byte, 1)
		_, err := file.Read(b)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}
		if b[0] == '\n' {
			break
		}
		if b[0] < '0' || b[0] > '9' {
			log.Fatalf("invalid character in input: %v", b[0])
		}
		length := int(b[0] - '0')
		if isSpace {
			d.append(&Blocks{id: -1, length: length})
		} else {
			d.append(&Blocks{id: id, length: length})
			id++
		}
	}
	return d
}

func answer1() int {
	disk := readInput()
	disk.compactWithFragmentation()
	return disk.checkSum()
}

// -----------------------------------------------------------------------

// PART 2
// Now don't move file blocks but only entire files if they can fit in the empty blocks.
// As before, start from the rightmost file and try to move to the leftmost empty space.
// Return the new checksum.

// moveLeftmost moves the full blocks `b` to the leftmost empty space that fits them
// Return a boolean indicating if the blocks were moved.
func (d *Disk) moveLeftmost(b *Blocks) bool {
	if b.isEmpty() {
		return false
	}

	// we cannot use d.firstEmpty because it might be after b
	var firstEmpty *Blocks
	if d.first.isEmpty() {
		firstEmpty = d.first
	} else {
		firstEmpty = d.first.nextEmpty(b)
	}

	for empty := firstEmpty; empty != nil; {
		nextEmpty := empty.nextEmpty(b)
		if empty.length >= b.length {
			d.moveBlocks(b, empty)
			return true
		}
		empty = nextEmpty
	}
	return false
}

// compactWithoutFragmentation compacts the disk without fragmentation by moving entire
// files only, starting from the rightmost file and moving to the leftmost empty space.
func (d *Disk) compactWithoutFragmentation() {
	checked := make(map[*Blocks]bool)

	for file := d.last; file != nil; {
		nextFile := file.before
		if !checked[file] {
			d.moveLeftmost(file)
			checked[file] = true
		}
		file = nextFile
	}
}

func answer2() int {
	disk := readInput()
	disk.compactWithoutFragmentation()
	return disk.checkSum()
}

// printDisk prints the disk to stdout
func printDisk(d *Disk) {
	for b := d.first; b != nil; b = b.after {
		id := strconv.Itoa(b.id)
		if b.isEmpty() {
			id = "."
		}
		for i := 0; i < b.length; i++ {
			print(id)
		}
	}
	println()
}

// -----------------------------------------------------------------------

var correctAnswers = map[int]int{
	1: 6435922584968,
	2: 6469636832766,
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
