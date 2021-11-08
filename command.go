package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const layout = "2 Jan 2006 15:04"

type Todo struct {
	Userid     int
	Title      string
	Completed  bool
	StartTime  string
	FinishTime string
	Category   string
}

func AddTodo(collection *mongo.Collection, userId int, message string, category string, msgTime time.Time) string {
	msg := ""
	if message != "" {
		startTime := FormatTime(msgTime)
		finishTime := startTime
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := collection.InsertOne(ctx, Todo{userId, message, false, startTime, finishTime, category})
		if err != nil {
			log.Fatal(err)
		}
	} else {
		msg = "<i>❗Неверно написано дело.\n\n</i>"
	}
	return msg
}

func RemoveTodo(collection *mongo.Collection, userId int, index string, category string) string {
	// добавить проверку существования index дела
	msg := ""
	result, i := ValidityOfIndex(collection, userId, index)
	if result {
		todoList := AllTodoList(collection, userId, category)
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

func ToggleTodo(collection *mongo.Collection, userId int, index string, category string) string {
	msg := ""
	result, i := ValidityOfIndex(collection, userId, index)
	if result {
		todoList := AllTodoList(collection, userId, category)
		toggleTodo := todoList[i-1]
		filter := bson.M{"userid": userId, "title": toggleTodo.Title}
		// filter["userid"] = userId
		// filter["title"] = toggleTodo.Title
		update := bson.M{"$set": bson.M{"completed": !toggleTodo.Completed}}
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := collection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Fatal(err)
		}
		msg = "<i>Статус дела изменён.\n\n</i>"
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

func AllTodoList(collection *mongo.Collection, userId int, category string) []*Todo {
	filter := bson.M{"userid": userId, "category": category}
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

func Deadline(collection *mongo.Collection, userId int, indexAndData string, category string) string {
	msg := ""
	indexData := strings.Split(indexAndData, ". ")
	index := indexData[0]
	result, i := ValidityOfIndex(collection, userId, index)
	if result {
		data, err := ParseData(indexData[1])
		if err == nil {
			SaveFinishData(collection, userId, i, FormatTime(data), category)
			msg = "<i>Дата завершения дела установлена.\n\n</i>"
		} else {
			msg = "<i>❗Неверно указана дата.\n\n</i>"
		}
	} else {
		msg = "<i>❗Такое дело не существует.\n\n</i>"
	}
	return msg
}

// Выводить в сообщнеие статус срока выполнения дела
func PrintTodoList(todoList []*Todo, timeNow string) string {
	msg := fmt.Sprintf("<b>MyTodoList</b>\nКатегория: <b>%s</b>\n", todoList[0].Category)
	emoji := ""
	title := ""
	duration := ""
	for i := 0; i < len(todoList); i++ {
		if todoList[i].Completed {
			emoji = "🟢"
			title = "<s>" + todoList[i].Title + "</s>"
		} else {
			emoji = "🔴"
			title = todoList[i].Title
		}
		Duration := GetDuration(todoList[i].StartTime, todoList[i].FinishTime)
		DurationNow := GetDuration(timeNow, todoList[i].FinishTime)
		PersentOfDuration := DurationNow / Duration
		if PersentOfDuration <= 1 && PersentOfDuration > 0.75 {
			duration = "🌕"
		} else if PersentOfDuration <= 0.75 && PersentOfDuration > 0.5 {
			duration = "🌔"
		} else if PersentOfDuration <= 0.5 && PersentOfDuration > 0.25 {
			duration = "🌓"
		} else if PersentOfDuration <= 0.25 && PersentOfDuration > 0 {
			duration = "🌒"
		} else {
			duration = "🌑"
		}
		finish := ""
		if DurationNow > 0 {
			finish = "\n<code>" + todoList[i].FinishTime + "</code>"
		}
		msg += fmt.Sprintf("%s %s %v. %s %s\n", emoji, duration, i+1, title, finish)
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

func SaveFinishData(collection *mongo.Collection, userId int, index int, finishTime string, category string) {
	todoList := AllTodoList(collection, userId, category)
	toggleTodo := todoList[index-1]
	filter := bson.M{"userid": userId, "title": toggleTodo.Title}
	update := bson.M{"$set": bson.M{"finishtime": finishTime}}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal(err)
	}
}

func FormatTime(getTime time.Time) string {
	strTime := getTime.Format(layout)
	return strTime
}

func GetDuration(startTime string, finishTime string) float64 {
	start, _ := time.Parse(layout, startTime)
	finish, _ := time.Parse(layout, finishTime)
	duration := finish.Sub(start)
	strDuration := duration.Minutes()
	return strDuration
}

func ParseData(Data string) (time.Time, error) {
	data, err := time.Parse(layout, Data)
	return data, err
}

func ValidityOfIndex(collection *mongo.Collection, userId int, index string) (bool, int) {
	result := false
	i, _ := strconv.Atoi(index)
	count := GetCountTodos(collection, userId)
	if i > 0 && i <= int(count) {
		result = true
	}
	return result, i
}

func GetAllUserCategory(collection *mongo.Collection, userId int) []*Categorys {
	filter := bson.M{"userid": userId}
	// filter["userid"] = userId

	var results []*Categorys

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(ctx) {
		var elem Categorys
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

func PrintCategory(Category []*Categorys) string {
	msg := "<b>Ваши ктегории:</b>\n"
	for i := 0; i < len(Category); i++ {
		msg += fmt.Sprintf("%v. %s\n", i+1, Category[i].Category)
	}
	return msg
}

// Переписать функцию для удаления категории со всеми в ней делами
func RemoveCategory(collection *mongo.Collection, userId int, index string) string {
	// добавить проверку существования index дела
	msg := ""
	result, i := ValidityOfIndex(collection, userId, index)
	if result {
		category := GetAllUserCategory(collection, userId)
		removeCategory := category[i-1].Category
		filter := bson.M{"userid": userId, "category": removeCategory}
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			log.Fatal(err)
		}
		msg = "<i>Категория удалена.\n\n</i>"
	} else {
		msg = "<i>❗Такая категория не существует.\n\n</i>"
	}
	return msg
}
