package bot

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type NotificationBot struct {
	bot *tgbotapi.BotAPI
}

func NewNotificationBot(config *Config) (*NotificationBot, error) {
	bot, err := tgbotapi.NewBotAPI(config.TgbotapiKey)
	if err != nil {
		return nil, err
	}

	return &NotificationBot{
			bot: bot,
		},
		nil
}

func (n *NotificationBot) SendNotification(chatID int64, message string) error {
	_, err := n.bot.Send(tgbotapi.NewMessage(chatID, message))
	if err != nil {
		log.Println("Can't send message %v", err)
		return err
	}
	return nil
}

func (n *NotificationBot) Bot() error {

	fmt.Println("Authorized on account %s", n.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := n.bot.GetUpdatesChan(u)
	if err != nil {
		fmt.Println("Can't start channel for getting updates %v", err)
	}

	for update := range updates {
		if err != nil {
			fmt.Println("Can't handle updates %v", err)
		}
		reply := "I do not know how to answer"
		if update.Message == nil {
			continue
		}

		fmt.Println(update.Message.From.UserName, update.Message.From.ID, update.Message.Text)

		switch update.Message.Command() {
		case "start":
			reply = "Your chatID: " + strconv.FormatInt(update.Message.Chat.ID, 10) + ". I will send you notifications about Statistic Service there."
		}

		err := n.SendNotification(update.Message.Chat.ID, reply)
		if err != nil {
			fmt.Println("Can't send message %v", err)
			return err
		}
	}
	return nil
}
