package funcs

import "fmt"

func CalcLowerQuartile(speed []float64) float64 {
	n := float64(len(speed))
	b := 3 * (n + 1) / 4
	c := int(b)
	fmt.Println(b)	
	var d = b - float64(c)
	q1 := speed[c-1] + (speed[c]-speed[c-1])*d
	return q1
}
