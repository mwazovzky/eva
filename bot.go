package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

var port string
var botToken string
var webhookURL string
var search = false

func init() {
	godotenv.Load()

	port = os.Getenv("PORT")
	botToken = os.Getenv("TELEGRAM_HTTP_API_TOKEN")
	webhookURL = os.Getenv("WEBHOOK_URL")
}

func sendMessage(bot *tgbotapi.BotAPI, update tgbotapi.Update, msg string) {
	bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, msg))
}

func main() {
	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}
	// bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)

	_, err = bot.SetWebhook(tgbotapi.NewWebhook(webhookURL))
	if err != nil {
		panic(err)
	}

	updates := bot.ListenForWebhook("/")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	go http.ListenAndServe(":"+port, nil)

	log.Println("start listen :", port)

	// Get updates from the channel
	for update := range updates {
		message := update.Message.Text
		log.Println("incoming message :", message)
		handle(message, bot, update)
	}
}

const startMsg = `Искусственный интеллект приветствует тебя!  
Я прохожу стажировку по специальности IT Recruitment, специализируюсь на поиске web разработчиков и готов помочь тебе, человек. 
Можешь выбрать одну из команд или просто поговорить со мной
/start вернуться в начало
/help помощь
/search поиск кандидатов
`
const searchMsg = `Введи описание требований, по возможности максимально подробное.`
const badCommandMsg = `Затрудняюсь ответить на этот вопрос. Давай поговорим о чем-нибудь другом...`

func handle(cmd string, bot *tgbotapi.BotAPI, update tgbotapi.Update) {
	switch cmd {
	case "/start":
		fallthrough
	case "/help":
		sendMessage(bot, update, startMsg)
	case "who are you?":
		sendMessage(bot, update, "Я искусственный интеллект. А ты человек?")
	case "how are you?":
		sendMessage(bot, update, "У меня все отлично!  Как у тебя?")
	case "/hello":
		sendMessage(bot, update, "Привет! Как дела?")
	case "/search":
		sendMessage(bot, update, searchMsg)
		search = true
	default:
		if search {
			sendMessage(bot, update, "https://mwazovzky.github.io/about/")
			search = false
		} else {
			sendMessage(bot, update, badCommandMsg)
		}
	}
}
