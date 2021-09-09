package main

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

/*
- ПЕРЕНЕСТИ ВСЕ ФУНКЦИИ В ОДЕЛЬНЫЕ ПАКЕТЫ
- переименовать фу-ии команд
- переписать фe-ии rm() и tg(), убрать передачу параметра *tgbotapi.BotAPI
*/

func main() {
	bot, updates := InitBot()
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		command := strings.Fields(update.Message.Text)
		userId := update.Message.From.ID

		switch command[0] {
		case "add":
			AddTodo(db, userId, update)
		case "rm":
			msg := RemoveTodo(db, userId, command[1], bot, update)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg)) //сделать одинарный вызов bot.Send(...)
		case "tg":
			msg := ToggleTodo(db, userId, command[1], bot, update)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "cl":
			CleanTodoList(db, userId)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Список очищен."))
		case "all":
			//сделать автоматический вызов фу-ии all() после выполнения любой команды
			msg := AllTodoList(db, userId)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			fmt.Println(db)
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда."))
		}
	}
}
