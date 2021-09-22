package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	client, collection := InitMongo()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer client.Disconnect(ctx)

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
			msg = AddTodo(collection, userId, update) //добавить возвращаемое значение в фу-ию AddTodo()
			msg += PrintTodoList(AllTodoList(collection, userId))
		case "rm":
			msg = RemoveTodo(collection, userId, command[1])
			msg += PrintTodoList(AllTodoList(collection, userId))
		case "tg":
			msg = ToggleTodo(collection, userId, command[1])
			msg += PrintTodoList(AllTodoList(collection, userId))
		case "cl":
			msg = CleanTodoList(collection, userId) //добавить возвращаемое значение в фу-ию CleanTodoList()
			msg += PrintTodoList(AllTodoList(collection, userId))
		case "all":
			msg = PrintTodoList(AllTodoList(collection, userId))
		default:
			msg = "<i>❗Неизвестная команда.</i>"
		}
		fmt.Println(collection)
		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
