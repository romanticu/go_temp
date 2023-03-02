package funcs

import (
	"fmt"
	"io/ioutil"
	"log"
)

func ReadFile1() {
	content, err := ioutil.ReadFile("a.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(content))

}

func WriteFile1() {
	err := ioutil.WriteFile("./output2.txt", []byte("测试文件2"), 0666) //写入文件(字节数组)
	if err != nil {
		log.Fatal(err.Error())
	}
}
