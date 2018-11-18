package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"math/rand"
	"raqb.it/catbot/utils"
	"strings"
)

var pronouns = []string{"he", "she", "they"}

func Buy(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, globalEnv *GlobalEnv, cmdEnv *CommandEnv) error {

	if len(parts) < 2 {
		utils.ChannelMesageSendError(s, m.ChannelID, "Please specify a name for your new cat!")
		return nil
	}

	catName := strings.Join(parts[1:], " ")

	isValidName := utils.CheckCatName(catName)

	if !isValidName {
		utils.ChannelMesageSendError(s, m.ChannelID, "The name specified is not valid.")
		return nil
	}

	exists, err := globalEnv.Db.CatNameExists(cmdEnv.User.ID, catName)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("Could not verify if cat name already exists")
		return err
	}

	if exists {
		utils.ChannelMesageSendError(s, m.ChannelID, fmt.Sprintf(
			"%s, you already have a cat named %s!",
			m.Author.Mention(),
			catName,
		))
		return nil
	}

	if cmdEnv.User.Money < globalEnv.Config.CatCost {
		utils.ChannelMesageSendError(s, m.ChannelID, fmt.Sprintf(
			"%s, you don't have enough credits for a cat! You can get more by using **%sdaily** every day.",
			m.Author.Mention(),
			globalEnv.Config.CommandPrefix,
		))
		return nil
	}

	randomPronoun := getRandomPronoun()

	err = globalEnv.Db.CreateCatForUser(cmdEnv.User.ID, getRandomCryptoKittyId(), catName, randomPronoun)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("Could not create cat")
		return err
	}

	err = globalEnv.Db.UserRemoveMoney(cmdEnv.User.ID, globalEnv.Config.CatCost)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("Could not remove money from user")
		return err
	}

	utils.ChannelMessageSendEmote(s, m.ChannelID, "🐱",
		fmt.Sprintf(
			"%s, you've just purchased a new cat!",
			m.Author.Mention(),
		),
	)

	return nil
}

func getRandomPronoun() string {
	return pronouns[rand.Intn(len(pronouns))]
}

func getRandomCryptoKittyId() int {
	// Using 1000000 as max id here
	return rand.Intn(1000000)
}
