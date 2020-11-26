package main

import (
	"io/ioutil"
	"log"
	"mvdan.cc/xurls/v2"
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

			rx := xurls.Relaxed()
			urls := rx.FindAllString(m.Text, -1)
			log.Println("verbose: ", urls)

			for _, url := range urls {
				if url[0:4] != "http" {
					url = "http://" + url
				}
				successful := false
				hasTitle := false
				var title string

				// get title
				count := 10
				for {
					if count <= 0 {
						break
					}
					resp, err := http.Get(url)
					if err != nil {
						log.Println("err: http ", err)
						count--
						continue
					}
					var ok bool
					// TODO: maybe we can get web content to telegraph ?

					if strings.Contains(url, "mp.weixin.qq.com"){
						title, ok = GetHtmlTitleWechat(resp.Body)
					} else {
						title, ok = GetHtmlTitle(resp.Body)
					}
					resp.Body.Close()
					if ok {
						hasTitle = true
					}
					successful = true
					break
				}
				if successful {
					addTask(url, title)
					if hasTitle {
						log.Println("verbose: add:", url, title)
						Send(m.Sender, "Added: "+title+" "+url, menu, tb.NoPreview)
					} else {
						log.Println("verbose: add:", url)
						log.Println("verbose: get title failed")
						Send(m.Sender, "Added: "+url, "Get titles failed.", menu, tb.NoPreview)
					}
				} else {
					Sendln(m.Sender, "Add URL failed.", url)
				}
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
				Send(m.Sender, string(bytes), menu, tb.NoPreview)
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

		// manage tasks
		manageTasksMarkup := &tb.ReplyMarkup{}
		var rmBtn tb.Btn = manageTasksMarkup.Data("删除", "remove")
		var bottomBtn tb.Btn = manageTasksMarkup.Data("置底", "bottom")
		manageTasksMarkup.Inline(
			manageTasksMarkup.Row(rmBtn, bottomBtn),
		)

		btnManage := menu.Text("队列管理")
		bot.Handle(&btnManage, func(m *tb.Message) {
			for _, task := range Tasks {
				Send(m.Sender, task.URL, manageTasksMarkup)
			}
		})

		bot.Handle(&rmBtn, func(c *tb.Callback) {
			// Always respond!
			bot.Respond(c, &tb.CallbackResponse{})

			url := strings.Trim(c.Message.Text, " ")
			finishTask(url)
			Sendln(c.Sender, "rm:", url)
		})

		bot.Handle(&bottomBtn, func(c *tb.Callback) {
			// Always respond!
			bot.Respond(c, &tb.CallbackResponse{})

			url := strings.Trim(c.Message.Text, " ")
			setLastTask(url)
			Sendln(c.Sender, "pinned:", url)
		})

		menu.Reply(menu.Row(btnView), menu.Row(btnManage), menu.Row(btnHistory))
	})
}
