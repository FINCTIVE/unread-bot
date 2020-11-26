package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	Launch(func(bot *tb.Bot) {
		initTask()

		menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

		bot.Handle("/hello", func(m *tb.Message) {
			pass := CheckUser(m.Sender)
			if !pass {
				return
			}
			Send(m.Sender, "hello!"+m.Sender.LastName+" "+m.Sender.LastName, menu)
		})

		// add url
		bot.Handle(tb.OnText, func(m *tb.Message) {
			pass := CheckUser(m.Sender)
			if !pass {
				return
			}
			successful := false
			hasTitle := false
			var title string

			if strings.Contains(m.Text, "http") {
				// get title
				count := 10
				for {
					if count <= 0 {
						break
					}
					resp, err := http.Get(m.Text)
					if err != nil {
						log.Println("err: http ", err)
						count--
						continue
					}
					var ok bool
					// TODO: maybe we can get web content to telegraph ?
					title, ok = GetHtmlTitle(resp.Body)
					resp.Body.Close()
					if ok {
						log.Println("verbose: add:", m.Text, title)
						hasTitle = true
					} else {
						log.Println("verbose: add:", m.Text)
						log.Println("verbose: get title failed")
					}
					successful = true
					break
				}
			}
			if successful {
				addTask(m.Text, title)
				if hasTitle {
					Sendln(m.Sender, "Added: "+title, m.Text, menu)
				} else {
					Send(m.Sender, "Added: "+m.Text, "Get titles failed.", menu)
				}
			} else {
				Sendln(m.Sender, "Add URL failed.", menu)
			}
		})

		// view history
		btnHistory := menu.Text("历史记录")
		bot.Handle(&btnHistory, func(m *tb.Message) {
			bytes, err := ioutil.ReadFile("history.log")
			if err != nil {
				log.Println("verbose: history", err)
				Send(m.Sender, "历史记录为空", menu)
			} else {
				Send(m.Sender, string(bytes), menu)
			}
		})

		// view tasks
		btnView := menu.Text("未读列表")
		bot.Handle(&btnView, func(m *tb.Message) {
			output := ""
			for _, task := range Tasks {
				output += task.Title + "\n"
				output += "=> " + task.URL + "\n\n"
			}
			Send(m.Sender, output, menu, tb.NoPreview)
		})

		menu.Reply(menu.Row(btnView), menu.Row(btnHistory))
	})
}
