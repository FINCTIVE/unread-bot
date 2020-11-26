package main

import (
	"log"
	"net/http"
	"testing"
)

func TestGetHtmlTitleWechat(t *testing.T) {
	resp, err := http.Get("https://mp.weixin.qq.com/s/BIxfCgssvTs1Qp6sznexIg")
	if err != nil {
		log.Println("err: http ", err)
	}
	title, ok := GetHtmlTitleWechat(resp.Body)
	_ = resp.Body.Close()
	log.Println(title, ok)
}
