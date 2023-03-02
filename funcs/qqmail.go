package funcs

import (
	"log"
	"math/rand"
	"time"

	"gopkg.in/gomail.v2"
)

func TakeAPicture() {
	for {
		var h, min int
		var get bool
		if time.Now().Hour() == 9 && !get {
			h, min = getTime()
			get = true
		}
		var t = time.Now()
		if t.Hour() == h && t.Minute() == min {
			sendEMail()
			get = false
		}
		time.Sleep(time.Minute)
	}
}

func getTime() (int, int) {
	var h int
	var min int
	rand.Seed(time.Now().UnixNano())
	for {
		h = rand.Intn(23)
		if h >= 10 {
			break
		}
	}

	for {
		min = rand.Intn(59)
		break
	}

	return h, min
	
}

func sendEMail() {
	sender := "mrfgq@qq.com"        //发送者qq邮箱
	authCode := "cxnaxfdkqziubjaf"  //qq邮箱授权码
	mailTitle := "拍照片啦"             //邮件标题
	mailBody := "干嘛呢？干嘛呢？干嘛呢？该拍照片啦" //邮件内容,可以是html

	//接收者邮箱列表
	mailTo := []string{
		"lijie@rimepevc.com",
	}

	m := gomail.NewMessage()
	m.SetHeader("From", sender)       //发送者腾讯企业邮箱账号
	m.SetHeader("To", mailTo...)      //接收者邮箱列表
	m.SetHeader("Subject", mailTitle) //邮件标题
	m.SetBody("text/html", mailBody)  //邮件内容,可以是html

	//发送邮件服务器、端口、发送者qq邮箱、qq邮箱授权码
	//服务器地址和端口是腾讯的
	d := gomail.NewDialer("smtp.qq.com", 587, sender, authCode)
	if err := d.DialAndSend(m); err != nil {
		log.Println("send mail failed", err)
		return
	}

	log.Println("success")
}
