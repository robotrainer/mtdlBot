package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"go.mongodb.org/mongo-driver/bson"
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

func AllTodoList(collection *mongo.Collection, userId int) []*Todo {
	filter := bson.M{}
	filter["userid"] = userId

	var results []*Todo

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		var elem Todo
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(ctx)
	return results
}

func PrintTodoList(todoList []*Todo) string {
	msg := "<b>MyTodoList</b>\n"
	emoji := ""
	title := ""
	for i := 0; i < len(todoList); i++ {
		if todoList[i].Completed {
			emoji = "🟢"
			title = "<s>" + todoList[i].Title + "</s>"
		} else {
			emoji = "🔴"
			title = todoList[i].Title
		}
		msg += fmt.Sprintf("%s %v. %s\n", emoji, i+1, title)
	}
	return msg
}
