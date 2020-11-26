package main

import (
	"io"
	"log"
	"strings"

	"golang.org/x/net/html"
)

func GetHtmlTitle(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Println("err: html.Parse:", err)
		return "", false
	}
	return traverse(doc)
}

func traverse(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data, true
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverse(c)
		if ok {
			return result, ok
		}
	}
	return "", false
}

func GetHtmlTitleWechat(r io.Reader) (string, bool) {
	doc, err := html.Parse(r)
	if err != nil {
		log.Println("err: html.Parse:", err)
		return "", false
	}
	return traverseWechat(doc)
}

func traverseWechat(n *html.Node) (string, bool) {
	if n.Type == html.ElementNode && n.Data == "h2" {
		for _, attr := range n.Attr {
			if attr.Key == "id" && attr.Val == "activity-name" {
				return strings.TrimSpace(n.FirstChild.Data), true
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverseWechat(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}