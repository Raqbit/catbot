package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
)

// Max Crypto Kitty ID used to generate random ones
const MaxCryptoKittyID = 1000000

var pronouns = []string{"he", "she", "they"}

func Buy(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, ctx *CmdContext) error {
	if len(parts) < 2 {
		_, _ = ChannelMesageSendError(s, m.ChannelID, "Please specify a name for your new cat!")
		return nil
	}

	catName := strings.Join(parts[1:], " ")

	if !isValidCatName(catName) {
		_, _ = ChannelMesageSendError(s, m.ChannelID, "The name specified is not valid.")
		return nil
	}

	exists, err := models.Cats.CatNameExists(ctx.Store, ctx.User, catName)

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

	if ctx.User.Money < ctx.Bot.Config.CatCost {
		_, _ = ChannelMesageSendError(s, m.ChannelID, fmt.Sprintf(
			"%s, you don't have enough credits for a cat! You can get more by using **%sdaily** every day.",
			m.Author.Mention(),
			ctx.Bot.Config.CommandPrefix,
		))
		return nil
	}

	randomPronoun := getRandomPronoun()
	cryptoKitty := getRandomCryptoKittyId()

	if err = models.Cats.CreateForUser(ctx.Store, ctx.User, cryptoKitty, catName, randomPronoun); err != nil {
		logrus.WithError(err).Errorln("Could not create cat")
		return err
	}

	if err = ctx.User.ModifyMoney(ctx.Store, -ctx.Bot.Config.CatCost); err != nil {
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
	return rand.Intn(MaxCryptoKittyID)
}
