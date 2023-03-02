package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
	"unicode/utf8"

	_ "net/http/pprof"

	"golang.org/x/crypto/bcrypt"
)

func RandInt(rangeInt int) int {
	rand.Seed(time.Now().UnixNano())
	for {
		num := rand.Intn(10)
		fmt.Println(num)
		if num != 0 {
			return num
		}
	}
}

type Deal struct {
	dealSize int
}

func (d Deal) Add() {
	d.dealSize = d.dealSize + 1
}
func (d Deal) Sub() {
	d.dealSize = d.dealSize - 2
}
func (d Deal) GetDealSize() int {
	return d.dealSize
}

func addQuery(eq map[string]interface{}) {
	eq["query"] = "bbbbbb"
}

//20010801

func isRecentDate(d string) bool {
	timeString := "20060102"
	lastYear := time.Now().AddDate(-1, 0, 0)
	lastYearYesterday := lastYear.AddDate(0, 0, -1)
	lastYearTomorrow := lastYear.AddDate(0, 0, 1)
	if len(d) == 8 {

		if d == lastYear.Format(timeString) || d == lastYearYesterday.Format(timeString) ||
			d == lastYearTomorrow.Format(timeString) {
			return true
		}
	}

	return false
}

func sortarr(n []int) {
	for i := 0; i < len(n)/2; i++ {
		n[i], n[len(n)-i-1] = n[len(n)-i-1], n[i]
	}
}

type fc func()

func Group(relativePath string) {
	// for _, item := range handlers {
	// 	item()
	// }
}

func f1() {
	fmt.Println("f1")
}
func f2() {
	fmt.Println("f2")
}

type D struct {
	Name        string        `json:"name"`
	DisplayName string        `json:"display_name"`
	Description string        `json:"description"`
	Properties  []interface{} `json:"properties"`
}

func FilterEmoji(content string) string {
	new_content := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			new_content += string(value)
		}
	}
	return new_content
}

func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

func minWindow(s string, t string) string {
	var left, right int
	var win = make(map[byte]int)
	var need = make(map[byte]int)
	var count int
	for i := range t {
		need[t[i]] = need[t[i]] + 1
	}

	var start = -1
	var end = len(s)
	//var l = len(s)
	for right < len(s) {
		c := s[right]
		right++
		if need[c] != 0 {
			win[c] = win[c] + 1
			if win[c] == need[c] {
				count++
			}
		}

		for count == len(need) {
			//fmt.Println(left, right, " == ", right-left)
			if right-left < end-start {
				fmt.Print("before >> ", right-start, " right ", right, " start ", start, " ")
				end = right
				start = left
				fmt.Println(left, right, start, " >> ", right-left)
			}
			//if right-left < l {
			//	l = right - left
			//	start = left
			//	fmt.Println(left, right, " ==*** ", right-left)
			//}
			d := s[left]
			if need[d] != 0 {
				if win[d] == need[d] {
					count--
				}
				win[d] = win[d] - 1
			}

			left++
		}

	}

	if start == -1 {
		return ""
	}
	fmt.Println(start, end, s[start:end])
	//fmt.Println(start, l, s[start:start+l])
	//return s[start : start+l]
	return s[start:end]
}

func main() {
	//"cabwefgewcwaefgcf"
	//"cae"
	fmt.Println(minWindow("cabwefgewcwaefgcf", "cae"))
}

// CompareHashAndPassword ...
func CompareHashAndPassword(hashed, password string) bool {
	if ok := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); ok == nil {
		return true
	}

	return false
}

// HashPassword ...
func HashPassword(password string) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(string(hashed))
}
