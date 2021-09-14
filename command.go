package main

import (
	"context"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	Userid    int
	Title     string
	Completed bool
}

func AddTodo(collection *mongo.Collection, userId int, update tgbotapi.Update) {
	title := strings.Replace(update.Message.Text, "add", "", 1)
	if title != "" {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := collection.InsertOne(ctx, Todo{userId, title, false})
		if err != nil {
			log.Fatal(err)
		}
	}
}

// func RemoveTodo(data map[int]todos, userId int, command string) string {
// 	msg := ""
// 	id, err := strconv.Atoi(command)
// 	if err != nil {
// 		msg = err.Error()
// 	} else if id <= len(data[userId]) && id > 0 {
// 		// возможно, переписать в функцию неполного смещения
// 		for _id, _ := range data[userId] {
// 			if _id > id {
// 				data[userId][id] = data[userId][_id]
// 				id = _id
// 			}
// 		}
// 		delete(data[userId], id)
// 		msg = fmt.Sprintf("<i>ℹ️ Дело %v удалено.</i>", command)
// 	} else {
// 		msg = "<i>ℹ️ Такое дело не существует.</i>"
// 	}
// 	return msg
// }

// func ToggleTodo(data map[int]todos, userId int, command string) string {
// 	// добавить проверку существования такого дела
// 	msg := ""
// 	id, err := strconv.Atoi(command)
// 	if err != nil {
// 		msg = err.Error()
// 	} else {
// 		for _id, message := range data[userId] {
// 			if _id == id {
// 				message.completed = !message.completed
// 				data[userId][_id] = message
// 			}
// 		}
// 		msg = fmt.Sprintf("<i>ℹ️ Статус дела %v изменён.</i>", id)
// 	}
// 	return msg
// }

// func CleanTodoList(data map[int]todos, userId int) string {
// 	msg := ""
// 	//сделать проверку на наличие выполненных дел
// 	data[0] = make(todos)
// 	// возможно, перенести в функцию копирования карты
// 	i := 1
// 	for _id, message := range data[userId] {
// 		if !message.completed {
// 			data[0][i] = data[userId][_id]
// 			i++
// 		}
// 	}
// 	data[userId] = make(todos)
// 	data[userId] = data[0]
// 	delete(data, 0)
// 	msg = "<i>ℹ️ Список очищен.</i>"
// 	return msg
// }

// func AllTodoList(data map[int]todos, userId int) string {
// 	//добавить проверку на непустой списка дела
// 	msg := "<b>MyTodoList</b>\n"
// 	emoji := ""
// 	title := ""
// 	for i := 1; i <= len(data[userId]); i++ {
// 		if data[userId][i].completed {
// 			emoji = "🟢"
// 			title = "<s>" + data[userId][i].title + "</s>"
// 		} else {
// 			emoji = "🔴"
// 			replacer := strings.NewReplacer("<s>", "", "</s>", "")
// 			title = replacer.Replace(data[userId][i].title)
// 		}
// 		data[userId][i] = todo{title, data[userId][i].completed}
// 		msg += fmt.Sprintf("%s %v. %s\n", emoji, i, data[userId][i].title)
// 	}
// 	return msg
// }
