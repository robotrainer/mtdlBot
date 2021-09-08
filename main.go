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

func main() {
	bot, updates := initBot()
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		command := strings.Fields(update.Message.Text)
		userId := update.Message.From.ID
		msg := ""

		switch command[0] {
		case "add":
			add(db, userId, update)
		case "rm":
			msg := rm(db, userId, command[1], bot, update)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "tg":
			id, err := strconv.Atoi(command[1])
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}
			for _id, message := range db[userId] {
				if _id == id {
					message.completed = !message.completed
					db[userId][_id] = message
				}
			}
			msg = fmt.Sprintf("Дело %v выполнено.", id)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "cl":
			db[0] = make(todos)
			i := 1
			for _id, message := range db[userId] {
				if !message.completed {
					db[0][i] = db[userId][_id]
					i++
				}
			}
			db[userId] = make(todos)
			db[userId] = db[0]
			delete(db, 0)
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Список очищен."))
		case "all":
			msg = ""
			for i := 1; i <= len(db[userId]); i++ {
				emoji := "🔴"
				if db[userId][i].completed {
					emoji = "🟢"
				}
				msg += fmt.Sprintf("%s %v. %s\n", emoji, i, db[userId][i].title)
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			fmt.Println(db)
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда."))
		}
	}
}
