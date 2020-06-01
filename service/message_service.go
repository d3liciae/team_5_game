package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"team_5_game/config"
	"team_5_game/model/telegram"
)

func ProcessWebhookMessage(update *telegram.Update) {
	log.Println("Processing webhook message:", convertToString(update))
	message := update.Message
	callbackQuery := update.CallbackQuery

	if message != nil {
		if message.Text == "/start" && isCommand(message) {
			RegisterUser(message)
			СhooseАction(message)
			return
		}
	}

	if callbackQuery != nil {
		if strings.HasPrefix(callbackQuery.Data, "CLAN_SELECT") {
			ProcessClanSelection(callbackQuery)
		}
	}

	if callbackQuery != nil {
		if strings.HasPrefix(callbackQuery.Data, "START_BATTLE") {
			ProcessBattleStarting(callbackQuery)
		}
	}
}

func SendMessage(chatID int64, message string, replyMarkup *telegram.InlineKeyboardMarkup) {
	err := sendMessage(chatID, message, replyMarkup)
	if err != nil {
		log.Println("Error in sending message:", err)
	}
}

func EditMessageReplyMarkup(chatID int64, messageID int64, replyMarkup *telegram.InlineKeyboardMarkup) {
	err := editMessageReplyMarkup(chatID, messageID, replyMarkup)
	if err != nil {
		log.Println("Error in editing message reply markup:", err)
	}
}

func convertToString(update *telegram.Update) string {
	out, err := json.Marshal(update)
	if err != nil {
		log.Println("Could not marshal update message", err)
		return "[Unable to convert to string]"
	}

	return string(out)
}

func sendMessage(chatID int64, message string, replyMarkup *telegram.InlineKeyboardMarkup) error {
	log.Println("Sending message to the chat:", chatID, " message: ", message)
	reqBody := &telegram.NewMessage{
		ChatID:      chatID,
		Text:        message,
		ReplyMarkup: replyMarkup,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(
		"https://api.telegram.org/bot"+config.BotToken()+"/sendMessage",
		"application/json",
		bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	log.Println("Message sent successfully")
	return nil
}

func editMessageReplyMarkup(chatID int64, messageID int64, replyMarkup *telegram.InlineKeyboardMarkup) error {
	log.Println("Editing message reply markup chat:", chatID, " message: ", messageID)
	reqBody := &telegram.EditMessageReplyMarkup{
		ChatID:      chatID,
		MessageID:   messageID,
		ReplyMarkup: replyMarkup,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(
		"https://api.telegram.org/bot"+config.BotToken()+"/editMessageReplyMarkup",
		"application/json",
		bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	log.Println("Message changed successfully")
	return nil
}

func isCommand(message *telegram.Message) bool {
	if message.Entities == nil || len(*message.Entities) == 0 {
		return false
	}

	entity := (*message.Entities)[0]
	return entity.Offset == 0 && entity.Type == "bot_command"
}
