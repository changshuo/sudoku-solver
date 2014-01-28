package main

import (
	"errors"
	"fmt"
	"strconv"
)

// blank data structure
// Represents a variable that may or may not be filled
type blank struct {
	value, row, col int
}

// Returns a instance of blank
func makeBlank(value, row, col int) blank {
	return blank{value, row, col}
}

// Returns whether a blank is empty
func (b blank) isEmpty() bool {
	return b.value == 0
}

// Returns a pretty string value of the blank for printing
func (b blank) prettyValue() string {
	if b.isEmpty() {
		return "⬚"
	}
	return strconv.Itoa(b.value)
}

// box data structure
// Represents a box region in the puzzle.
// Is an array to pointers of blanks that are in the box
type box [9]*blank

// Sudoku data structure
// Represents an abstraction of the Sudoku board
// Tracks empty blanks
type Sudoku struct {
	blanks [9][9]blank
	boxes  [3][3]box
}

// Returns a pointer to a new instance of Sudoku
func NewSudoku(source string) *Sudoku {
	bytes := []byte(source)

	if len(bytes) != 81 {
		panic("Invalid source string: string must have length of 81")
	}

	// setup blanks
	sudoku := new(Sudoku)
	for index, c := range bytes {
		n, err := strconv.Atoi(string(c))
		if err != nil {
			panic("Invalid source string: " + err.Error())
		}
		row, col := index/9, index%9
		sudoku.blanks[row][col] = makeBlank(n, row, col)

		// compute boxes
		boxRow, boxCol := row/3, col/3
		inBoxRow, inBoxCol := row%3, col%3
		sudoku.boxes[boxRow][boxCol][inBoxRow*3+inBoxCol] = &sudoku.blanks[row][col]
	}

	return sudoku
}

// Helper function that panics if blank index is out of bound
func assertRange(row, col int) {
	if row < 0 || col < 0 || row > 8 || col > 8 {
		panic(fmt.Sprintf("Blank is out of range: row %d, column %d", row, col))
	}
}

// Get a pointer to a blank with index (row, col)
func (s *Sudoku) Get(row, col int) *blank {
	assertRange(row, col)
	return &(s.blanks[row][col])
}

// Set value of (row, col) to val
func (s *Sudoku) Set(row, col, val int) {
	assertRange(row, col)
	s.Get(row, col).value = val
}

// Print the beautiful board
func (s *Sudoku) PrettyPrint() {
	edgeSep := "+-----+-----+-----+\n"
	fmt.Print(edgeSep) // top horizontal separator
	for row := 0; row < 9; row++ {
		fmt.Print("|")
		for col := 0; col < 9; col++ {
			fmt.Print(s.Get(row, col).prettyValue())
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

func (s *Sudoku) isRowConsistent(b *blank) bool {
	for _, newB := range s.blanks[b.row] {
		if newB.value == b.value && newB.col != b.col {
			return false
		}
	}
	return true
}
func (s *Sudoku) isColConsistent(b *blank) bool {
	col := b.col
	for row := 0; row < 9; row++ {
		newB := s.blanks[row][col]
		if newB.value == b.value && row != b.row {
			return false
		}
	}
	return true
}
func (s *Sudoku) getBox(b *blank) *box {
	return &(s.boxes[b.row/3][b.col/3])
}
func (s *Sudoku) isBoxConsistent(b *blank) bool {
	for _, newB := range s.getBox(b) {
		if b.value == newB.value &&
			b.row == newB.row &&
			b.col == newB.col {
			return false
		}
	}
	return true
}

// Check if Block b is consistent regarding Sudoku game rules
func (s *Sudoku) isConsistent(b *blank) bool {
	// Unfilled blank is always consistent
	if b.isEmpty() {
		return true
	}
	// Check row, column, and box for dupe
	if s.isRowConsistent(b) && s.isColConsistent(b) && s.isBoxConsistent(b) {
		return true
	}
	return false
}

// Check if Sudoku is solved
func (s *Sudoku) IsComplete() (bool, error) {
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if b := s.Get(row, col); b.isEmpty() { // if b is empty, then it's not complete
				return false, nil
			} else if s.isConsistent(b) == false { // if b is not consistent, then it's not complete
				return false, errors.New(fmt.Sprintf("inconsistent blank found at row %d, column %d", row, col))
			}
		}
	}
	return true, nil
}

// Solves puzzle by performing backtracking search
func (s *Sudoku) Solve() error {
	// check completeness and correctness
	if complete, err := s.IsComplete(); complete || err != nil {
		return err
	}

	return nil
}

func main() {
	source := "030080006500294710000300500005010804420805039108030600003007000041653002200040060"
	s := NewSudoku(source)
	err := s.Solve()
	if err != nil {
		panic(err.Error())
	}

	// output solved result
	s.PrettyPrint()
	/*
		s.PrettyPrint()

		// test: output all boxes
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				var b box = s.boxes[i][j]
				for k := 0; k < 9; k++ {
					fmt.Print(b[k].prettyValue())
				}
				fmt.Println()
			}
		}
	*/
}
