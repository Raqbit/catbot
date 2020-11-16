package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
)

var pronouns = []string{"he", "she", "they"}

func Buy(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, context *Context) error {

	if len(parts) < 2 {
		_, _ = ChannelMesageSendError(s, m.ChannelID, "Please specify a name for your new cat!")
		return nil
	}

	catName := strings.Join(parts[1:], " ")

	isValidName := CheckCatName(catName)

	if !isValidName {
		_, _ = ChannelMesageSendError(s, m.ChannelID, "The name specified is not valid.")
		return nil
	}

	exists, err := models.Cats.CatNameExists(context.Store, context.User, catName)

	if err != nil {
		logrus.WithError(err).Errorln("Could not verify if cat name already exists")
		return err
	}

	if exists {
		_, _ = ChannelMesageSendError(s, m.ChannelID, fmt.Sprintf(
			"%s, you already have a cat named %s!",
			m.Author.Mention(),
			catName,
		))
		return nil
	}

	if context.User.Money < context.Config.CatCost {
		_, _ = ChannelMesageSendError(s, m.ChannelID, fmt.Sprintf(
			"%s, you don't have enough credits for a cat! You can get more by using **%sdaily** every day.",
			m.Author.Mention(),
			context.Config.CommandPrefix,
		))
		return nil
	}

	randomPronoun := getRandomPronoun()
	cryptoKitty := getRandomCryptoKittyId()

	// TODO: USE TRANSACTION
	err = models.Cats.CreateForUser(context.Store, context.User, cryptoKitty, catName, randomPronoun)

	if err != nil {
		logrus.WithError(err).Errorln("Could not create cat")
		return err
	}

	err = context.User.ModifyMoney(context.Store, -context.Config.CatCost)

	if err != nil {
		logrus.WithError(err).Errorln("Could not remove money from user")
		return err
	}

	_, _ = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content: formatEmojiMessage("🐱", fmt.Sprintf(
			"%s, you've just purchased a new cat!",
			m.Author.Mention(),
		)),
		Embed: createCatProfileEmbed(&models.Cat{
			Name:          catName,
			CryptoKittyID: cryptoKitty,
			Pronoun:       randomPronoun,
			Hunger:        100,
		}),
	})

	return nil
}

func getRandomPronoun() string {
	return pronouns[rand.Intn(len(pronouns))]
}

func getRandomCryptoKittyId() int {
	// Using 1000000 as max id here
	return rand.Intn(1000000)
}
