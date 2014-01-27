package main

import (
	"errors"
	"fmt"
	"strconv"
)

// Blank data structure
// Represents a variable that may or may not be filled
type Blank struct {
	value, row, col int
}

// Returns a instance of blank
func makeBlank(value, row, col int) Blank {
	return Blank{value, row, col}
}

// Returns whether a blank is empty
func (b Blank) IsEmpty() bool {
	return b.value == 0
}

// Returns a pretty string value of the blank for printing
func (b Blank) PrettyValue() string {
	if b.IsEmpty() {
		return "⬚"
	}
	return strconv.Itoa(b.value)
}

// Sudoku data structure
// Represents an abstraction of the Sudoku board
// Tracks empty blanks
type Sudoku struct {
	blanks [9][9]Blank
	empty  []*Blank
}

// Returns a pointer to a new instance of Sudoku
func newSudoku(source string) *Sudoku {
	bytes := []byte(source)

	if len(bytes) != 81 {
		panic("Invalid source string: string must have length of 81")
	}

	sudoku := new(Sudoku)
	for index, c := range bytes {
		n, err := strconv.Atoi(string(c))
		if err != nil {
			panic("Invalid source string: " + err.Error())
		}
		row, col := index/9, index%9
		sudoku.blanks[row][col] = makeBlank(n, row, col)
	}
	return sudoku
}

// Helper function that panics if Blank index is out of bound
func assertRange(row, col int) {
	if row < 0 || col < 0 || row > 8 || col > 8 {
		panic(fmt.Sprintf("Blank is out of range: row %d, column %d", row, col))
	}
}

// Get a pointer to a Blank with index (row, col)
func (s *Sudoku) Get(row, col int) *Blank {
	assertRange(row, col)
	return &(s.blanks[row][col])
}

// Set value of (row, col) to val, updates the empty list
func (s *Sudoku) Set(row, col, val int) {
	assertRange(row, col)
	// G
	b := s.Get(row, col)
	if val == 0 {

	}
}

/*
type function func(Blank *b)

func (s *Sudoku) Map(function) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			function(s.Get(row, col))
		}
	}
}
*/

// Print the beautiful board
func (s *Sudoku) PrettyPrint() {
	edgeSep := "+-----+-----+-----+\n"
	fmt.Print(edgeSep) // top horizontal separator
	for row := 0; row < 9; row++ {
		fmt.Print("|")
		for col := 0; col < 9; col++ {
			fmt.Print(s.Get(row, col).PrettyValue())
			if (col+1)%3 == 0 {
				fmt.Print("|")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
		if row != 8 && (row+1)%3 == 0 {
			fmt.Print("|-----+-----+-----|\n") // horizontal separator
		}
	}
	fmt.Print(edgeSep) // bottom horizontal separator
}

// Check if Block b is consistent regarding Sudoku game rules
func (s *Sudoku) IsConsistent(b *Blank) bool {
	// Unfilled blank is always consistent
	if b.IsEmpty() {
		return true
	}
	// Check row, column, and box for dupe
	if s.IsRowConsistent(b.row) && s.IsColConsistent(b.col) && s.IsBoxConsistent(b.box) {
		return true
	}
	return false
}

// Check if Sudoku is solved
func (s *Sudoku) IsComplete() (bool, error) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if b := s.Get(row, col); b.IsEmpty() { // if b is empty, then it's not complete
				return false, nil
			} else if s.IsConsistent(b) == false { // if b is not consistent, then it's not complete
				return false, errors.New(fmt.Sprintf("inconsistent blank found at row %d, column %d", row, col))
			}
		}
	}
	return true, nil
}

// Main AI entrance point
func (s *Sudoku) Solve() error {
	if complete, err := s.IsComplete(); complete || err != nil {
		return err
	}
	return nil
}

func main() {
	source := "030080006500294710000300500005010804420805039108030600003007000041653002200040060"
	newSudoku(source).PrettyPrint()
}
