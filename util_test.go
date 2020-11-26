package main

// test RunCommand interactively
// import (
// 	"fmt"
// 	"log"
// )

// func main() {
// 	bytes, stop, done := RunCommand("ping", "baidu.com")

// 	for {
// 		var inputStr string
// 		fmt.Scan(&inputStr)
// 		if inputStr == "stop" {
// 			*stop <- struct{}{}
// 		} else if inputStr == "l" {
// 			fmt.Println(string(*bytes))
// 		}

// 		select {
// 		case err := <-*done:
// 			if err != nil {
// 				log.Println(err)
// 			}
// 			break
// 		default:
// 		}
// 	}
// }
