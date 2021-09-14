package main

import (
	"fmt"
	"strings"

	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {

	//Создаём нового клиента базы данных с поключением по указанному URL
	client, err := mongo.NewClient(options.Client().
		ApplyURI(GetURI("uri.txt")))
	if err != nil {
		log.Fatal(err)
	}

	//Подключаемся к базе данных
	fmt.Println("Client connecting...")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second) //ждём ответа в течении 10 сек
	err = client.Connect(ctx)                                           //если ответа нет, вернёт ошибку
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx) //откладываем момент отключения от базы данных
	//проверяем подклчение
	fmt.Println("Ping connecting...")
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second) //ждём ответа в течении 10 сек
	err = client.Ping(ctx, readpref.Primary())                         //отправляет сигнал ping, чтобы проверить, может ли клиент быть подключен к базе данных
	if err != nil {
		log.Fatal(err)
	}

	databases, err := client.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(databases)

	collection := client.Database("mtdlBot").Collection("todolist")

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
			AddTodo(collection, userId, update) //добавить возвращаемое значение в фу-ию AddTodo()
			// if msg != "" {
			// 	msg += "\n\n" + AllTodoList(db, userId)
			// }
		case "rm":
			RemoveTodo(collection, userId, command[1])
			// msg += "\n\n" + AllTodoList(db, userId)
		// case "tg":
		// 	msg = ToggleTodo(db, userId, command[1])
		// 	msg += "\n\n" + AllTodoList(db, userId)
		// case "cl":
		// 	msg = CleanTodoList(db, userId) //добавить возвращаемое значение в фу-ию CleanTodoList()
		// 	msg += "\n\n" + AllTodoList(db, userId)
		case "all":
			msg = PrintTodoList(AllTodoList(collection, userId))
		default:
			msg = "<i>ℹ️ Неизвестная команда.</i>"
		}
		fmt.Println(collection)
		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
