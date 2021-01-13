package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Link struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

var Links []Link

const filename string = "links.json"

func initStorage() {
	bytes, err := ioutil.ReadFile(filename)
	if err == nil {
		_ = json.Unmarshal(bytes, &Links)
	}
}

func addLink(url, title string) {
	//for _, task := range Links {
	//	if task.URL == url {
	//		return
	//	}
	//}
	Links = append(Links, Link{
		URL:   url,
		Title: title,
	})

	bytes, err := json.Marshal(Links)
	if err != nil {
		log.Println("err: storage: ", err)
	}
	err = ioutil.WriteFile(filename, bytes, os.ModePerm)

	if err != nil {
		log.Println("err: storage: ", err)
	}
}