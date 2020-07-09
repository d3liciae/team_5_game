package service

import (
	"errors"
	"log"
	"team_5_game/model/database"
	"team_5_game/model/telegram"
)

// RegisterUser creates the profile for the user with unique Telegram user ID
func RegisterUser(message *telegram.Message) error {
	log.Println("Start user registration")

	user, err := GetUserFromDB(message.From.ID)
	if user != nil {
		SendMessage(
			message.Chat.ID,
			"Hello "+message.From.FirstName+" you're already registered!!!",
			nil)
		err = errors.New("User is already registered")
		log.Println(err)
	} else {
		user := database.User{
			ID:            message.From.ID,
			FirstName:     message.From.FirstName,
			ClanID:        0,
			BattleCounter: 0,
			WinCounter:    0,
			CurrentBattle: 0,
			CurrentPos:    0,
			// Clan:          &database.Clan{ID: 0, Name: "NO_CLAN"},
		}
		err = SaveUserToDB(&user)
		if err == nil {
			SendMessage(
				message.Chat.ID,
				"Hello "+message.From.FirstName+" thank you for registration!!!",
				nil)
			log.Println("User successfully registered")
		} else {
			SendMessage(
				message.Chat.ID,
				"Something went wrong please contact the administrator!!!",
				nil)
			log.Println("User not registered")
		}
	}
	return err
}

// RegisterBattle creates the battle profile with the given id
func RegisterBattle(id int64) {
	log.Println("Start battle registration")

	battle, _ := GetBattleFromDB(id)
	if battle != nil {
		log.Println("Battle is already registered")
	} else {
		battle := database.Battle{
			ID: id,
		}
		for i := range battle.Sector {
			for j := range battle.Sector[i] {
				battle.Sector[i][j].ID = i*FieldWidth + j + 1
			}
		}

		err := SaveBattleToDB(&battle)
		if err != nil {
			log.Println("Battle not registered")
		} else {
			log.Println("Battle successfully registered")
		}
	}
}

// RestartGame cancel the user's current battle and clan
func RestartGame(message *telegram.Message) {
	log.Println("Restarting game")

	user, err := GetUserFromDB(message.From.ID)
	if err != nil {
		log.Println("Could not get user", err)
		return
	}

	EditMessageReplyMarkup(message.Chat.ID, user.CurrentBattlefieldMessageID, nil) // Deleting previous battlefield
	ResetCurrentBattlefieldMessageID(user)
	SendMessage(message.Chat.ID, "Previous battle is cancelled", nil)
	SetUserCurrentBattle(user, 0)
	SendClanSelectionMenu(message)
}

func ResetCurrentBattlefieldMessageID(user *database.User) {
	user.CurrentBattlefieldMessageID = 0
	err := SaveUserToDB(user)
	if err != nil {
		log.Println("Cannot save user to DB with new CurrentBattlefieldMessageID", err)
	}
}
