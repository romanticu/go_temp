package funcs

import "fmt"

func TestSlice() {
	s := []int{5}
	s = append(s, 7)
	s = append(s, 8)
	x := append(s, 11)
	y := append(s, 12)
	fmt.Println(x, y)
	fmt.Println(s)
}

func TestSlice2() {
	s := []int{}
	x := append(s, 11)
	y := append(s, 12)
	fmt.Println(x, y)
	fmt.Println(s)
}