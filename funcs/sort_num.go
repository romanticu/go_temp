package funcs

import "fmt"


func SortNum() {
	nums := []int{
		1,93,90,93,20,76,76,99,99,99,3,3,199,200,66,24,90,100,299,99,
	}

	var max = -1
	var scond = -1
	

	for _, n := range nums {
		if max == -1 {
			max = n
		} else {
			if n > max {
				scond = max
				max = n
			} else {
				if scond == -1 {
					scond = n
				} else if n > scond {
					scond = n
				}
			}
		}
	}
	fmt.Println(max, scond)
}