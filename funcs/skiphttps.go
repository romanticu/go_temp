package funcs

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
)

func NewBranch() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	var url = "https://gitlab.mvalley.com/api/v4/projects/111/repository/branches?branch=newbranch2&ref=dev1"
	req, _ := http.NewRequest("POST", url, nil)
	req.Header.Add("PRIVATE-TOKEN", "bwcxLHH82vyGQTMFJyUR")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func GetTags() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	var url = "https://gitlab.mvalley.com/api/v4/projects/530/repository/tags"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("PRIVATE-TOKEN", "bwcxLHH82vyGQTMFJyUR")
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}