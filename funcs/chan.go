package funcs

import (
	"fmt"
	"time"
)

var msg chan string

func TestChan() {
	msg := make(chan string, 10)
	go func() {
		for  {
			str := <- msg
			fmt.Println(str)
		}
	}()
	for i := 0; i < 20; i++ {
		strContent := fmt.Sprintf("%v - test", i)
		msg <- strContent
	}

	time.Sleep(time.Second*10)
}

func TestGoroutine() {
	fmt.Println("start")
	chs := make([]chan int, 10)
	ret := make([]chan int, 10)

	for i := 0; i < 10; i++ {
		chs[i] = make(chan int)
		ret[i] = make(chan int)
	}

	for i := 0; i < 10; i++ {
		go func (i int)  {
			defer func (i int)  {
				ret[i] <-1
			}(i)
			chs[i] <- 1
			fmt.Println(i)
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-chs[i]
		<-ret[i]
	}
	select{}
}