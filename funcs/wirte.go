package funcs

import "os"

func GetWrite(f string) *os.File {
	o, err := os.OpenFile(f, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	return o
}
