package funcs

import (
	"fmt"
	"regexp"
	"strings"
)


func PWDRegx() {
	pwdReg := regexp.MustCompile(`[0-9a-zA-Z\~\!\@\#\$%\^\&\*\(\)\_\+\=\-\{\}\\\|\[\]\:\'\,\.\/\?\"\<\>]`)

	pwd := "0934qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM<>"
	for _, ch := range strings.Split(pwd, "") {
		fmt.Println(ch)
		fmt.Println(pwdReg.MatchString(ch))
	}
}