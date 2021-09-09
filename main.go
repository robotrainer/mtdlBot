package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type todo struct {
	title     string
	completed bool
}

type todos map[int]todo

var db = map[int]todos{}

/*
- ПЕРЕНЕСТИ ВСЕ ФУНКЦИИ В ОДЕЛЬНЫЕ ПАКЕТЫ
- переименовать фу-ии команд
- переписать фe-ии rm() и tg(), убрать передачу параметра *tgbotapi.BotAPI
*/
func getToken(filename string) string {
	token, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(token)
}

func initBot() (*tgbotapi.BotAPI, tgbotapi.UpdatesChannel) {
	bot, err := tgbotapi.NewBotAPI(getToken("token.txt"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}
	return bot, updates
}

func add(data map[int]todos, userId int, update tgbotapi.Update) {
	if _, ok := data[userId]; !ok {
		data[userId] = make(todos)
	}
	messageId := len(data[userId]) + 1
	newMessage := strings.Replace(update.Message.Text, "add ", "", 1)
	data[userId][messageId] = todo{newMessage, false}
}

func rm(data map[int]todos, userId int, command string, bot *tgbotapi.BotAPI, update tgbotapi.Update) string {
	msg := ""
	id, err := strconv.Atoi(command)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
	}
	if id <= len(data[userId]) && id > 0 {
		// возможно, переписать в функцию неполного смещения
		for _id, _ := range data[userId] {
			if _id > id {
				data[userId][id] = data[userId][_id]
				id = _id
			}
		}
		delete(data[userId], id)
		msg = fmt.Sprintf("Дело %v удалено.", command)
	} else {
		msg = "Такое дело не существует."
	}
	return msg
}

func tg(data map[int]todos, userId int, command string, bot *tgbotapi.BotAPI, update tgbotapi.Update) string {
	// добавить проверку существования такого дела
	id, err := strconv.Atoi(command)
	if err != nil {
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
	}
	for _id, message := range data[userId] {
		if _id == id {
			message.completed = !message.completed
			data[userId][_id] = message
		}
	}
	msg := fmt.Sprintf("Статус дела %v изменён.", id)
	return msg
}

func cl(data map[int]todos, userId int) {
	//сделать проверку на наличие выполненных дел
	data[0] = make(todos)
	// возможно, перенести в функцию копирования карты
	i := 1
	for _id, message := range data[userId] {
		if !message.completed {
			data[0][i] = data[userId][_id]
			i++
		}
	}
	data[userId] = make(todos)
	data[userId] = data[0]
	delete(data, 0)
}

func all(data map[int]todos, userId int) string {
	//добавить проверку на непустой списка дела
	msg := ""
	for i := 1; i <= len(db[userId]); i++ {
		emoji := "🔴"
		if db[userId][i].completed {
			emoji = "🟢"
		}
		msg += fmt.Sprintf("%s %v. %s\n", emoji, i, db[userId][i].title)
	}
	return msg
}

func main() {
	bot, updates := initBot()
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		command := strings.Fields(update.Message.Text)
		userId := update.Message.From.ID

		switch command[0] {
		case "add":
			add(db, userId, update)
		case "rm":
			msg := rm(db, userId, command[1], bot, update)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg)) //сделать одинарный вызов bot.Send(...)
		case "tg":
			msg := tg(db, userId, command[1], bot, update)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "cl":
			cl(db, userId)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Список очищен."))
		case "all":
			//сделать автоматический вызов фу-ии all() после выполнения любой команды
			msg := all(db, userId)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			fmt.Println(db)
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда."))
		}
	}
}
