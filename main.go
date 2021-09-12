package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	var db = map[int]todos{}
	bot, updates := InitBot()
	for update := range updates {
		if update.Message == nil {
			continue
		}

		command := strings.Fields(update.Message.Text)
		userId := update.Message.From.ID
		msg := ""

		switch command[0] {
		case "add":
			msg = AddTodo(db, userId, update) //добавить возвращаемое значение в фу-ию AddTodo()
			if msg != "" {
				msg += "\n\n" + AllTodoList(db, userId)
			}
		case "rm":
			msg = RemoveTodo(db, userId, command[1])
			msg += "\n\n" + AllTodoList(db, userId)
		case "tg":
			msg = ToggleTodo(db, userId, command[1])
			msg += "\n\n" + AllTodoList(db, userId)
		case "cl":
			msg = CleanTodoList(db, userId) //добавить возвращаемое значение в фу-ию CleanTodoList()
			msg += "\n\n" + AllTodoList(db, userId)
		case "all":
			msg = AllTodoList(db, userId)
		default:
			msg = "<i>ℹ️ Неизвестная команда.</i>"
		}
		fmt.Println(db)
		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
