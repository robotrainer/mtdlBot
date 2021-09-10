package main

import (
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type todo struct {
	title     string
	completed bool
}

type todos map[int]todo

func AddTodo(data map[int]todos, userId int, update tgbotapi.Update) {
	if _, ok := data[userId]; !ok {
		data[userId] = make(todos)
	}
	messageId := len(data[userId]) + 1
	newMessage := strings.Replace(update.Message.Text, "add ", "", 1)
	data[userId][messageId] = todo{newMessage, false}
}

func RemoveTodo(data map[int]todos, userId int, command string) string {
	msg := ""
	id, err := strconv.Atoi(command)
	if err != nil {
		msg = err.Error()
	} else if id <= len(data[userId]) && id > 0 {
		// возможно, переписать в функцию неполного смещения
		for _id, _ := range data[userId] {
			if _id > id {
				data[userId][id] = data[userId][_id]
				id = _id
			}
		}
		delete(data[userId], id)
		msg = fmt.Sprintf("Дело %v удалено.", command)
	} else {
		msg = "Такое дело не существует."
	}
	return msg
}

func ToggleTodo(data map[int]todos, userId int, command string) string {
	// добавить проверку существования такого дела
	msg := ""
	id, err := strconv.Atoi(command)
	if err != nil {
		msg = err.Error()
	} else {
		for _id, message := range data[userId] {
			if _id == id {
				message.completed = !message.completed
				data[userId][_id] = message
			}
		}
		msg = fmt.Sprintf("Статус дела %v изменён.", id)
	}
	return msg
}

func CleanTodoList(data map[int]todos, userId int) {
	//сделать проверку на наличие выполненных дел
	data[0] = make(todos)
	// возможно, перенести в функцию копирования карты
	i := 1
	for _id, message := range data[userId] {
		if !message.completed {
			data[0][i] = data[userId][_id]
			i++
		}
	}
	data[userId] = make(todos)
	data[userId] = data[0]
	delete(data, 0)
}

func AllTodoList(data map[int]todos, userId int) string {
	//добавить проверку на непустой списка дела
	msg := "<s>MyTodoList!</s>\n"
	for i := 1; i <= len(data[userId]); i++ {
		emoji := "🔴"
		if data[userId][i].completed {
			emoji = "🟢"
		}
		msg += fmt.Sprintf("%s %v. %s\n", emoji, i, data[userId][i].title)
	}
	return msg
}
