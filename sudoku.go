package main

import (
	"fmt"
	"strconv"
)

type Blank struct {
	value, row, col int
}

func makeBlank(value, row, col int) Blank {
	return Blank{value, row, col}
}

func (b Blank) IsEmpty() bool {
	return b.value == 0
}

func (b Blank) PrettyValue() string {
	if b.IsEmpty() {
		return "⬚"
	}
	return strconv.Itoa(b.value)
}

type Sudoku struct {
	blanks [9][9]Blank
	empty  []*Blank
}

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

func (s *Sudoku) Get(row, col int) *Blank {
	if row < 0 || col < 0 || row > 8 || col > 8 {
		panic("Requested blank is out of range: row " + strconv.Itoa(row) + ", column " + strconv.Itoa(col))
	}
	return &(s.blanks[row][col])
}

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

func main() {
	source := "030080006500294710000300500005010804420805039108030600003007000041653002200040060"
	newSudoku(source).PrettyPrint()
}
