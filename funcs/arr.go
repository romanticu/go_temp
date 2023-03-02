package funcs

import "fmt"

type TModel struct {
	ID          *string
	MyName      string
	ServicePort string
	Age         int
}

func TestArr() {
	type Arr struct {
		A []string
	}

	var a Arr

	arr2 := make([]string, 0)

	arr2 = append(arr2, a.A...)
	fmt.Println(arr2)
}
