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
	ferstStart := true

	for update := range updates {
		if ferstStart {
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Обновление @mtdlBot:\n- для добавления нового дела теперь не надо вводить команду /add, просто напиши дело и отправь его боту;\n- при вводе команды /start, бот приветсвует пользователя."))
			ferstStart = false
		}

		if update.Message == nil {
			continue
		}

		userId := update.Message.From.ID
		msg := ""

		if update.Message.IsCommand() && flag == "" {
			command := update.Message.Command()

			switch command {
			case "start":
				msg = fmt.Sprintf("Приветсвую тебя, %s! Пока тут нет руководства, но оно скоро появится.(надеюсь)", update.Message.From.FirstName)
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
		} else if flag == "rm" {
			msg = RemoveTodo(collection, userId, update.Message.Text)
			msg += PrintTodoList(AllTodoList(collection, userId))
			flag = ""
		} else if flag == "tg" {
			msg = ToggleTodo(collection, userId, update.Message.Text)
			msg += PrintTodoList(AllTodoList(collection, userId))
			flag = ""
		} else {
			msg = AddTodo(collection, userId, update.Message.Text)
			msg += PrintTodoList(AllTodoList(collection, userId))
		}
		fmt.Println(collection)
		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
