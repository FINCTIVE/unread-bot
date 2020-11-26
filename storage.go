package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Task struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

var Tasks []Task

func initTask() {
	bytes, err := ioutil.ReadFile("tasks.json")
	if err == nil {
		_ = json.Unmarshal(bytes, &Tasks)
	}
}

func updateStorage() {
	bytes, err := json.Marshal(Tasks)
	if err != nil {
		log.Println("err: storage: ", err)
	}
	err = ioutil.WriteFile("tasks.json", bytes, os.ModePerm)

	if err != nil {
		log.Println("err: storage: ", err)
	}
}

func getTask(url string) (Task, bool) {
	for _, task := range Tasks {
		if task.URL == url {
			return task, true
		}
	}
	return Task{}, false
}

func addTask(url, title string) {
	for _, task := range Tasks {
		if task.URL == url {
			return
		}
	}
	Tasks = append(Tasks, Task{
		URL:   url,
		Title: title,
	})
	updateStorage()
}

func setFirstTask(url string) {
	for i, task := range Tasks {
		if task.URL == url {
			var t Task = Tasks[0]
			Tasks[0] = task
			Tasks[i] = t
		}
	}
	updateStorage()
}

func setLastTask(url string) {
	for i, task := range Tasks {
		if task.URL == url {
			var t Task = Tasks[len(Tasks)-1]
			Tasks[len(Tasks)-1] = task
			Tasks[i] = t
		}
	}
	updateStorage()
}

// save to history.log
// remove from Tasks
func finishTask(url string) {
	task, ok := getTask(url)
	if !ok {
		return
	}

	rmTask(url)
	// log file
	f, err := os.OpenFile("history.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("error: history:", err)
	}
	defer f.Close()
	if _, err := f.WriteString(url + " " + task.Title + "\n"); err != nil {
		log.Println("error: history:", err)
	}
}

// remove from Tasks
func rmTask(url string) {
	var newTasks []Task
	for _, task := range Tasks {
		if task.URL != url {
			newTasks = append(newTasks, task)
		}
	}
	Tasks = newTasks
	updateStorage()
}
