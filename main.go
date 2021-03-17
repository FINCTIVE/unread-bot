package main

import (
	"log"
	"mvdan.cc/xurls/v2"
	"net/http"
	"strings"

	tb "gopkg.in/tucnak/telebot.v2"
)

// If the received message has more than 5 urls,
// the bot will reply each of the urls one by one,
// because this is regarded as an importing action.
const startReplyURL = 5

// can only review the latest history
// but all history urls will be stored in the file.
const historyLength = 60

func main() {
	Launch(func(bot *tb.Bot) {
		initStorage()

		menu := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

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
				count := 3
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

					if strings.Contains(url, "mp.weixin.qq.com") {
						title, hasTitle = GetHtmlTitleWechat(resp.Body)
					} else {
						title, hasTitle = GetHtmlTitle(resp.Body)
					}
					_ = resp.Body.Close()

					successful = true
					break
				}

				if !successful {
					log.Println("error: the link can not be reached.")
				}

				addLink(url, title)
				log.Println("verbose: add:", url, title)
				if !hasTitle {
					log.Println("verbose: get title failed")
				}

				if len(urls) > startReplyURL {
					Send(m.Sender, url, menu)
				}
			}
		})

		bot.Handle("/history", func(m *tb.Message) {
			pass := CheckUser(m.Sender)
			if !pass {
				return
			}
			output := "Latest history: ...\n"
			startIndex := len(Links) - historyLength
			if startIndex < 0 {
				startIndex = 0
			}
			for i := startIndex; i <= len(Links)-1; i++ {
				output += Links[i].Title + "\n"
				output += "=> " + Links[i].URL + "\n\n"
			}
			Send(m.Sender, output, menu, tb.NoPreview)
		})

		bot.Handle("/hello", func(m *tb.Message) {
			pass := CheckUser(m.Sender)
			if !pass {
				return
			}
			Send(m.Sender, "hello!"+m.Sender.LastName+" "+m.Sender.LastName, menu)
		})

		err := bot.SetCommands(
			[]tb.Command{
				{"/hello", "hello"},
				{"/history", "view url history"},
			})
		if err != nil {
			log.Println(err)
		}
	})
}
