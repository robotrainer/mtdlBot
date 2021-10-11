package main

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Categorys struct {
	Userid   int
	Category string
}

var keyCategory = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Создать категорию"),
		tgbotapi.NewKeyboardButton("Разное"),
	),
)

var keboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Все дела"),
		tgbotapi.NewKeyboardButton("Категории"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Установить срок"),
		tgbotapi.NewKeyboardButton("Изменить статус"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Удалить дело"),
		tgbotapi.NewKeyboardButton("Удалить выполненные"),
	),
)

// var key = tgbotapi.ReplyKeyboardMarkup{}

func main() {

	client, collectionTodos, colCategory := InitMongo()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer client.Disconnect(ctx)

	bot, updates := InitBot()

	allUserId := GetAllUserId(collectionTodos, "userid") //Получить userid всех пользователей со списками дел для уведомлений.
	//Сообщение об обновлении, отправляется один раз сразу после запуска сервера.
	//При перезапуске сервера удалить строчку 22.
	// SendUpdateNotification(allUserId, bot)
	//Отправка ежедневных уведоблений.
	SendNotification(allUserId, bot)

	flag := ""

	for update := range updates {
		if update.Message == nil {
			continue
		}

		//получение времени запроса
		now := FormatTime(update.Message.Time())
		userId := update.Message.From.ID
		getMessage := update.Message.Text
		msg := ""

		// if update.Message.IsCommand() && flag == "" {
		// 	command := update.Message.Command()

		switch getMessage {
		// case "start":
		// 	msg = fmt.Sprintf("Приветсвую тебя, %s! Я бот, который поможет тебе сохранять важные дела и следить за их выполнением. Пока у меня не много функций, но ты уже можешь начать пользоватьс мной. Чтобы узнать список доступных команды, введи команду /help.", update.Message.From.FirstName)
		// case "help":
		// 	msg = "Доступные команды:\nчтобы добавить дело, просто напиши его сообщением и отправ мне\n/all - показать ваш список дел\n/toggle - изменить статус указанного дела, меняет с невыполненного, на выполненное, и наоборот\n/remove - удалить указанное дело\n/clean - удалить все выполненные дела\n/start - запуск бота\n/help - справочная информация, руководство\n/settings - настроить чат-бот"
		// case "settings":
		// 	msg = "Тут будут доступны настройки бота. Сейчас этот раздел в разработке."
		case "Удалить дело":
			msg = "Напишите номер удаляемого дела."
			flag = getMessage
		case "Изменить статус":
			msg = "Напишите номер дела."
			flag = getMessage
		case "Удалить выполненные":
			msg = CleanTodoList(collectionTodos, userId)
			msg += PrintTodoList(AllTodoList(collectionTodos, userId), now)
		case "Все дела":
			msg = PrintTodoList(AllTodoList(collectionTodos, userId), now)
		case "Установить срок":
			msg = "Напиши номер дела и дату." //сделать ввод даты выполнения с кнопки?
			flag = getMessage
		case "Категории":
			Msg := "Ваши категории:\n"
			Msg += "└ Разное" //добавить считывание категорий из БД
			m := tgbotapi.NewMessage(update.Message.Chat.ID, Msg)
			m.ReplyMarkup = keyCategory
			bot.Send(m)
			flag = getMessage
		case "Создать категорию":
			msg = "Напишите название новой категории."
			flag = getMessage
		default:
			if flag == "Удалить дело" {
				msg = RemoveTodo(collectionTodos, userId, update.Message.Text)
				msg += PrintTodoList(AllTodoList(collectionTodos, userId), now)
			} else if flag == "Изменить статус" {
				msg = ToggleTodo(collectionTodos, userId, update.Message.Text)
				msg += PrintTodoList(AllTodoList(collectionTodos, userId), now)
			} else if flag == "Установить срок" {
				msg = Deadline(collectionTodos, userId, update.Message.Text)
				msg += PrintTodoList(AllTodoList(collectionTodos, userId), now)
			} else if flag == "Категории" {
				category := update.Message.Text
				msg = PrintTodoList(CategoryTodoList(collectionTodos, category), now)
			} else if flag == "Создать категорию" {
				category := update.Message.Text
				ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
				_, err := colCategory.InsertOne(ctx, Categorys{userId, category})
				if err != nil {
					log.Fatal(err)
				}
				msg = "Категория создана."
			} else { //добавление нового дела
				msg = AddTodo(collectionTodos, userId, update.Message.Text, update.Message.Time()) //добавлять дела оп категориям
				msg += PrintTodoList(AllTodoList(collectionTodos, userId), now)
			}
			flag = ""
			// msg = AddTodo(collectionTodos, userId, update.Message.Text, update.Message.Time())
			// msg += PrintTodoList(AllTodoList(collectionTodos, userId), now)
			// msg = "<i>❗Неизвестная команда.</i>"
		}

		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		if flag == "" {
			msgParse.ReplyMarkup = keboard
		} else {
			msgParse.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
		bot.Send(msgParse)
	}
}
