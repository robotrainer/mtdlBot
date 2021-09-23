package main

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	client, collection := InitMongo()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer client.Disconnect(ctx)

	bot, updates := InitBot()

	flag := ""

	for update := range updates {
		if update.Message == nil {
			continue
		}

		userId := update.Message.From.ID
		msg := ""

		if update.Message.IsCommand() && flag == "" {
			command := update.Message.Command()

			switch command {
			case "add":
				msg = "Напишите новое дело."
				flag = command
			case "rm":
				msg = "Напишите номер удаляемого дела."
				flag = command
			case "tg":
				msg = "Напишите номер дела."
				flag = command
			case "cl":
				msg = CleanTodoList(collection, userId) //добавить возвращаемое значение в фу-ию CleanTodoList()
				msg += PrintTodoList(AllTodoList(collection, userId))
			case "all":
				msg = PrintTodoList(AllTodoList(collection, userId))
			default:
				msg = "<i>❗Неизвестная команда.</i>"
			}
		} else if flag == "add" {
			msg = AddTodo(collection, userId, update)
			msg += PrintTodoList(AllTodoList(collection, userId))
			flag = ""
		} else if flag == "rm" {
			msg = RemoveTodo(collection, userId, update.Message.Text)
			msg += PrintTodoList(AllTodoList(collection, userId))
			flag = ""
		} else if flag == "tg" {
			msg = ToggleTodo(collection, userId, update.Message.Text)
			msg += PrintTodoList(AllTodoList(collection, userId))
			flag = ""
		} else {
			msg = "<i>❗Неизвестная команда.</i>"
		}
		fmt.Println(collection)
		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
