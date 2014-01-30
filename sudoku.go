package main

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	useForwardChecking bool = false
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
			b.row != newB.row &&
			b.col != newB.col {
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

// Returns a blank that is unassigned
func (s *Sudoku) unassignedBlock() *blank {
	// TODO: heuristic
	for row := 0; row < 9; row++ {
		for col := 0; col < 9; col++ {
			if b := s.Get(row, col); b.isEmpty() {
				return b
			}
		}
	}
	return nil
}

// Returns true if assigning value to blank b is consistent
func (s *Sudoku) isAssignable(b *blank, value int) bool {
	oldval := b.value
	b.value = value
	assignable := s.isConsistent(b)
	b.value = oldval
	return assignable
}

// Returns a slice of ints that is ordered by some heuristics
func (s *Sudoku) orderDomainValues(b *blank) []int {
	// TODO
	return []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
}

// Infer on particular assignments, return whether inference is sucess,
// and a map of inferred values
func (s *Sudoku) inference(b *blank, value int) (bool, map[*blank]int) {
	return true, make(map[*blank]int)
}

// Implements a recursive backtrack search algorithm
func (s *Sudoku) backTrack() (bool, error) {
	// check completeness and correctness
	if done, err := s.IsComplete(); done || err != nil {
		return done, err
	}
	// get any unassigned block
	unassignedBlock := s.unassignedBlock()
	for _, val := range s.orderDomainValues(unassignedBlock) {
		if s.isAssignable(unassignedBlock, val) {
			// add {var = value} to assignment
			unassignedBlock.value = val
			// make inference
			success, inferences := s.inference(unassignedBlock, val)
			if success {
				// add inferences inferences to assignments
				for b, v := range inferences {
					b.value = v
				}
				// recursive call
				if done, _ := s.backTrack(); done { // return if done
					return true, nil
				}
				// If we reach here, the particular value is not correct
				// remove inferences
				for b, _ := range inferences {
					b.value = 0
				}
			}
			// remove assignment
			unassignedBlock.value = 0
		}
	}
	return false, errors.New("no consistent domain value available")
}

// Solves puzzle by performing backtracking search
// Pass Sudoku by value, returns a new completed
// sudoku or nil, as well as any error that occured
func (s *Sudoku) Solve() (*Sudoku, error) {
	// make a copy of puzzle
	newS := *s
	// run recursive algorithm on puzzle
	done, err := newS.backTrack()
	if !done {
		return nil, err
	}
	return &newS, err
}

func main() {
	//p1 := "030080006500294710000300500005010804420805039108030600003007000041653002200040060"
	//p2 := "308296000040008000502100087013000000780000035000000410120007803000800020000542106"
	p3 := "700000000600410250013095000860000000301000405000000086000840530042036007000000009"

	s := NewSudoku(p3)

	newS, err := s.Solve()
	if err != nil {
		panic(err.Error())
	}

	// output solved result
	fmt.Println("Original Sudoku puzzle:")
	s.PrettyPrint()
	fmt.Println("\nSolution:")
	newS.PrettyPrint()
}
