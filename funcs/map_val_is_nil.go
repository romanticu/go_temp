package funcs

import "fmt"

func TestMapValIsNil() {
	m := make(map[string][]string)

	m["1"] = append(m["1"], "test1")
	m["1"] = append(m["1"], "test2")

	fmt.Println(m["1"])
}