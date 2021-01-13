package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/net/context"
	tb "gopkg.in/tucnak/telebot.v2"
	"gopkg.in/yaml.v2"
)

var globalBot *tb.Bot
var GlobalConfig Config

// Launch loads the yaml configuration file, and start the bot.
func Launch(load func(bot *tb.Bot)) {
	configBytes, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = yaml.Unmarshal(configBytes, &GlobalConfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	globalBot, err = tb.NewBot(tb.Settings{
		Token:  GlobalConfig.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	load(globalBot)

	log.Println("verbose: bot started")
	globalBot.Start()
}

// if the message is too long, cut it into pieces
// and send seperately
const LongMessageLength = 4000
const maxRetry = 1000

// Send sends message. If failed, retry until it's successful.
// (to deal with poor network problem ...)
// Also, Send split long message to small pieces. (Telegram has message length limit.)
// and send them seperately.
func Send(sender *tb.User, message string, options ...interface{}) {
	lines := strings.Split(message, "\n")
	if message == "" || len(lines) == 0 {
		log.Println("verbose: sending empty message, quit...")
		return
	}

	var msgs []string
	splitDone := false
	for {
		var newMessage string
		for {
			if len(newMessage)+len(lines[0])+1 < LongMessageLength {
				newMessage += lines[0] + "\n"
				lines = lines[1:]
				if len(lines) == 0 {
					msgs = append(msgs, newMessage)
					splitDone = true
					break
				}
			} else {
				msgs = append(msgs, newMessage)
				break
			}
		}
		if splitDone {
			break
		}
	}

	retryCounter := 0
	for {
		_, err := globalBot.Send(sender, msgs[0], options...)
		if err != nil {
			log.Println("err: send:", msgs[0])
			log.Println(err)
			retryCounter++
			if retryCounter >= maxRetry {
				log.Println("err: send: tried ", maxRetry, " times. Give it up.")
				// for errors not related with network
				_, _ = globalBot.Send(sender, "Messages not sent, pls check your terminal log.", options...)
				break
			}
		} else {
			if len(msgs) == 1 {
				break
			} else {
				msgs = msgs[1:]
			}
		}
	}
}

// Sendln wraps Send for terminal message
// Also log to the terminal
func Sendln(sender *tb.User, v ...interface{}) {
	message := fmt.Sprintln(v...)
	log.Println(v...)
	Send(sender, message, tb.NoPreview)
}

// Sendf wraps Send for terminal message
// Also log to the terminal
func Sendf(sender *tb.User, format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Printf(format, v...)
	Send(sender, message, tb.NoPreview)
}

// CheckUser will check whether the username is in the config.yaml
// if the users property in config.yaml is set to be "*" or empty,
// all users will pass the check
func CheckUser(sender *tb.User) (pass bool) {
	if len(GlobalConfig.Users) == 0 || GlobalConfig.Users[0] == "*" {
		return true
	}

	pass = false
	for _, username := range GlobalConfig.Users {
		if username == sender.Username {
			pass = true
			break
		}
	}
	if pass == false {
		Send(sender, "Sorry, it's a bot for private usage.")
	}
	log.Println("verbose: check ", sender.Username, " - pass:", pass)
	return pass
}

// RunCommand runs a command in the backgroud and capture its output bytes in real time,
// combining stdout and stderr output,
// send struct{} to stop channel will stop the command,
// listen for done channel to wait the command finish.
// note: the output slice has no lock.
func RunCommand(name string, arg ...string) (output *[]byte, stop *chan struct{}, done *chan error) {
	var outputBytes []byte
	var stopChan = make(chan struct{}, 1)
	var doneChan = make(chan error)

	var cmdDone = make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, name, arg...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	r := io.MultiReader(stdout, stderr)
	cmd.Start()

	go func() {
		for {
			var buffer []byte = make([]byte, 1024)
			n, err := r.Read(buffer)
			if n > 0 {
				// log.Println(string(buffer))
				outputBytes = append(outputBytes, buffer...)
			}
			if err != nil {
				if err == io.EOF {
					break
				} else {
					// ignore
					log.Println(err)
					break
				}
			}
		}

		cmdDone <- cmd.Wait()
	}()

	go func() {
		select {
		case err := <-cmdDone:
			doneChan <- err
		case <-stopChan:
			cancel()
			doneChan <- nil
		}
	}()

	return &outputBytes, &stopChan, &doneChan
}
