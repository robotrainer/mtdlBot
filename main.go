package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/*
- переписать фe-ии rm() и tg(), убрать передачу параметра *tgbotapi.BotAPI
*/

func main() {
	var db = map[int]todos{}
	bot, updates := InitBot()
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		command := strings.Fields(update.Message.Text)
		userId := update.Message.From.ID
		msg := ""

		switch command[0] {
		case "add":
			AddTodo(db, userId, update)
			msg = AllTodoList(db, userId)
		case "rm":
			msg = RemoveTodo(db, userId, command[1])
			msg += "\n" + AllTodoList(db, userId)
		case "tg":
			msg = ToggleTodo(db, userId, command[1])
			msg += "\n" + AllTodoList(db, userId)
		case "cl":
			CleanTodoList(db, userId) //добавить возвращаемое значение в фу-ию CleanTodoList()
			msg = AllTodoList(db, userId)
		case "all":
			//сделать автоматический вызов фу-ии all() после выполнения любой команды
			msg = AllTodoList(db, userId)
		default:
			msg = "<i>Неизвестная команда.</i>"
		}
		fmt.Println(db)
		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
