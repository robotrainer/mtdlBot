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

func main() {
	token, err := ioutil.ReadFile("token.txt")
	if err != nil {
		panic(err)
	}
	bot, err := tgbotapi.NewBotAPI(string(token))
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

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		word := strings.Fields(update.Message.Text)
		userId := update.Message.From.ID
		msg := ""

		switch word[0] {
		case "ADD":
			if _, ok := db[userId]; !ok {
				db[userId] = make(todos)
			}
			messageId := len(db[userId]) + 1
			newMsg := strings.Replace(update.Message.Text, "ADD ", "", 1)
			db[userId][messageId] = todo{newMsg, false}
		case "RM":
			id, err := strconv.Atoi(word[1])
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}
			if id <= len(db[userId]) && id > 0 {
				for _id, _ := range db[userId] {
					if _id > id {
						db[userId][id] = db[userId][_id]
						id = _id
					}
				}
				delete(db[userId], id)
				msg = fmt.Sprintf("Дело %v удалено.", word[1])
			} else {
				msg = "Такое дело не существует."
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
		case "TG":
			id, err := strconv.Atoi(word[1])
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
		case "CL":
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
		case "ALL":
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
