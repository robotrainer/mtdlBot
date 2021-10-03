package main

import (
	"context"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	client, collection := InitMongo()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer client.Disconnect(ctx)

	bot, updates := InitBot()

	allUserId := GetAllUserId(collection, "userid") //Получить userid всех пользователей со списками дел для уведомлений.
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
		msg := ""

		if update.Message.IsCommand() && flag == "" {
			command := update.Message.Command()

			//добавить команды /help и /settings
			//настроить комманду /start
			switch command {
			case "start":
				msg = fmt.Sprintf("Приветсвую тебя, %s! Я бот, который поможет тебе сохранять важные дела и следить за их выполнением. Пока у меня не много функций, но ты уже можешь начать пользоватьс мной. Чтобы узнать список доступных команды, введи команду /help.", update.Message.From.FirstName)
				// добавить запоминание Chat.ID в базу данных
			case "help":
				msg = "Доступные команды:\nчтобы добавить дело, просто напиши его сообщением и отправ мне\n/all - показать ваш список дел\n/toggle - изменить статус указанного дела, меняет с невыполненного, на выполненное, и наоборот\n/remove - удалить указанное дело\n/clean - удалить все выполненные дела\n/start - запуск бота\n/help - справочная информация, руководство\n/settings - настроить чат-бот"
			case "settings":
				msg = "Тут будут доступны настройки бота. Сейчас этот раздел в разработке."
			case "remove":
				msg = "Напишите номер удаляемого дела."
				flag = command
			case "toggle":
				msg = "Напишите номер дела."
				flag = command
			case "clean":
				msg = CleanTodoList(collection, userId) //добавить возвращаемое значение в фу-ию CleanTodoList()
				msg += PrintTodoList(AllTodoList(collection, userId), now)
			case "all":
				msg = PrintTodoList(AllTodoList(collection, userId), now)
			case "deadline":
				msg = "Напиши номер дела и дату." //сделать ввод даты выполнения с кнопки
				flag = command
			default:
				msg = "<i>❗Неизвестная команда.</i>"
			}
		} else if flag == "remove" {
			msg = RemoveTodo(collection, userId, update.Message.Text)
			msg += PrintTodoList(AllTodoList(collection, userId), now)
			flag = ""
		} else if flag == "toggle" {
			msg = ToggleTodo(collection, userId, update.Message.Text)
			msg += PrintTodoList(AllTodoList(collection, userId), now)
			flag = ""
		} else if flag == "deadline" {
			//написать алгоритм действий
			// сохранять время deadline в БД
			//добавить команду добавления или изменения начала выполнения дела
			// ПЕРЕПИСАТЬ и ОФОРМИТЬ В ФУ-ИЮ
			msg = Deadline(collection, userId, update.Message.Text)
			msg += PrintTodoList(AllTodoList(collection, userId), now)
			flag = ""
		} else { //добавление нового дела
			msg = AddTodo(collection, userId, update.Message.Text, update.Message.Time())
			msg += PrintTodoList(AllTodoList(collection, userId), now)
		}
		fmt.Println(collection)
		msgParse := tgbotapi.NewMessage(update.Message.Chat.ID, msg)
		msgParse.ParseMode = tgbotapi.ModeHTML
		bot.Send(msgParse)
	}
}
