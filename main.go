package main

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strconv"
	"strings"
)

type todo struct {
	title     string
	completed bool
}

type todos map[int]todo

var db = map[int]todos{}

func main() {
	bot, err := tgbotapi.NewBotAPI("")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		word := strings.Fields(update.Message.Text)
		userId := update.Message.From.ID

		switch word[0] {
		case "ADD":
			if _, ok := db[userId]; !ok {
				db[userId] = make(todos)
			}
			messageId := len(db[userId]) + 1
			newMsg := strings.Replace(update.Message.Text, "ADD ", "", 1)
			db[userId][messageId] = todo{newMsg, false}
			// bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда ADD"))
		case "RM":
			id, err := strconv.Atoi(word[1])
			if err != nil {
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, err.Error()))
			}
			delete(db[userId], id)
			for _id, _ := range db[userId] {
				if _id > id {
					db[userId][id] = db[userId][_id]
					id = _id
				}
			}
			// bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда RM"))
		case "TG":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда TG"))
		case "CL":
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда CL"))
		case "ALL":
			msg := ""
			for id, message := range db[userId] {
				msg += fmt.Sprintf("%v %s %t\n", id, message.title, message.completed)
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			fmt.Println(db)
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
		}
	}
}
