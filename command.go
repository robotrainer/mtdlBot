package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Todo struct {
	Userid    int
	Title     string
	Completed bool
}

func AddTodo(collection *mongo.Collection, userId int, message string) string {
	msg := ""
	if message != "" {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := collection.InsertOne(ctx, Todo{userId, message, false})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg = "<i>❗Неверно написано дело.\n\n</i>"
	}
	return msg
}

func RemoveTodo(collection *mongo.Collection, userId int, index string) string {
	// добавить проверку существования index дела
	msg := ""
	count := GetCountTodos(collection, userId)
	i, _ := strconv.Atoi(index)
	if i > 0 && i <= int(count) {
		todoList := AllTodoList(collection, userId)
		removeTodo := todoList[i-1].Title
		filter := bson.M{"userid": userId, "title": removeTodo}
		// filter["userid"] = userId
		// filter["title"] = removeTodo
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		msg = "<i>Дело удалено.\n\n</i>"
	} else {
		msg = "<i>❗Такое дело не существует.\n\n</i>"
	}
	return msg
}

func ToggleTodo(collection *mongo.Collection, userId int, index string) string {
	// добавить проверку существования index дела
	msg := ""
	count := GetCountTodos(collection, userId)
	i, _ := strconv.Atoi(index)
	if i > 0 && i <= int(count) {
		todoList := AllTodoList(collection, userId)
		toggleTodo := todoList[i-1]
		filter := bson.M{"userid": userId, "title": toggleTodo.Title}
		// filter["userid"] = userId
		// filter["title"] = toggleTodo.Title
		update := bson.M{"$set": bson.M{"completed": !toggleTodo.Completed}}
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		result, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		msg = "<i>Статус дела изменён.\n\n</i>"
		fmt.Printf("%v %v\n", result.MatchedCount, result.ModifiedCount)
	} else {
		msg = "<i>❗Такое дело не существует.\n\n</i>"
	}
	return msg
}

func CleanTodoList(collection *mongo.Collection, userId int) string {
	filter := bson.M{"userid": userId, "completed": true}
	// filter["userid"] = userId
	// filter["completed"] = true
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.DeleteMany(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	msg := "<i>TodoList очищен.\n\n</i>"
	return msg
}

func AllTodoList(collection *mongo.Collection, userId int) []*Todo {
	filter := bson.M{"userid": userId}
	// filter["userid"] = userId

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

func GetCountTodos(collection *mongo.Collection, userId int) int64 {
	filter := bson.M{"userid": userId}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	return count
}
