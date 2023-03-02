package main

import (
	"fmt"
	"strings"
)

func main() {
	var s = "abc"
	fmt.Println(s)
	for _, c := range s {
		fmt.Println(c)
	}
}
func solveNQueens(n int) [][]string {
	res = make([][]string, 0)
	var m = make([][]string, 0)
	for i := 0; i < n; i++ {
		var row []string
		for j := 0; j < n; j++ {
			row = append(row, ".")
		}
		m = append(m, row)
	}
	backtrack(m, 0)
	fmt.Println("res ----- ")
	fmt.Println(res)
	return res
}

var res [][]string

func backtrack(m [][]string, row int) {
	if row == len(m) {
		var strs []string
		for _, item := range m {
			strs = append(strs, strings.Join(item, ""))
		}
		fmt.Println(m)
		res = append(res, strs)
		return
	}
	n := len(m)
	for i := 0; i < n; i++ {
		if !isValid(m, row, i) {
			continue
		}
		m[row][i] = "Q"
		backtrack(m, row+1)
		m[row][i] = "."
	}

}

func isValid(m [][]string, row, col int) bool {
	n := len(m)
	fmt.Println("row ", row, " col ", col)
	for i := 0; i < row; i++ {
		if m[i][col] == "Q" {
			return false
		}
	}

	for i, j := row-1, col-1; i >= 0 && j >= 0; i, j = i-1, j-1 {
		if m[i][j] == "Q" {
			return false
		}
	}
	// left
	for i, j := row-1, col+1; i >= 0 && j < n; i, j = i-1, j+1 {
		fmt.Println(i, j)
		if m[i][j] == "Q" {
			return false
		}
	}

	return true
}
