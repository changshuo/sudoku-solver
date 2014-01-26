package main

import (
	"fmt"
	"strconv"
)

type Sudoku struct {
	vars  [9][9]int
	empty []*int
}

func buildSudoku(source string) *Sudoku {
	bytes := []byte(source)

	if len(bytes) != 81 {
		panic("Invalid source string: string must have length of 81")
	}

	node := new(Sudoku)
	for index, c := range bytes {
		n, err := strconv.Atoi(string(c))
		if err != nil {
			panic("Invalid source string: " + err.Error())
		}
		node[index/9][index%9] = n
	}
	return node
}

func (n *Sudoku) PrettyPrint() {
	edgeSep := "+-----+-----+-----+\n"
	fmt.Print(edgeSep) // top horizontal separator
	for row := 0; row < 9; row++ {
		fmt.Print("|")
		for col := 0; col < 9; col++ {
			fmt.Print(n[row][col])
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
	buildSudoku(source).PrettyPrint()
}
