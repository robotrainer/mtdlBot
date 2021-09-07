package main

import (
	"fmt"
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
	bot, err := tgbotapi.NewBotAPI("1993669332:AAECB5_FyzH0RUpn_Md9dVBw9Fwh2yk6BHI")
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
		msg := ""

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
			// bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "🤩"))
		case "CL":
			for _id, message := range db[userId] {
				if message.completed {
					delete(db[userId], _id)
				}
			}
			i := 1
			for _id, _ := range db[userId] {
				db[userId][i] = db[userId][_id]
				i++
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Команда CL"))
		case "ALL":
			for id, message := range db[userId] {
				emoji := "🔴"
				if message.completed {
					emoji = "🟢"
				}
				msg += fmt.Sprintf("%s %v. %s\n", emoji, id, message.title)
			}
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
			fmt.Println(db)
		default:
			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда"))
		}
	}
}
