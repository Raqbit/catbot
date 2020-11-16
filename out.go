package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

func Out(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, context *Context) error {

	if len(parts) < 2 {
		_, _ = ChannelMesageSendError(s, m.ChannelID, "Please specify a cat to get the info of!")
		return nil
	}

	catName := strings.Join(parts[1:], " ")

	cat, err := models.Cats.GetByName(context.Store, context.User, catName)

	if err != nil {
		logrus.WithError(err).Errorln("Could not fetch cat from database")
		return err
	}

	if cat == nil {
		_, _ = ChannelMesageSendError(s, m.ChannelID, fmt.Sprintf(
			"%s, you do not have a cat with that name!",
			m.Author.Mention(),
		))
		return nil
	}

	if cat.Away {
		_, _ = ChannelMesageSendError(s, m.ChannelID, fmt.Sprintf(
			"%s, %s is already out!",
			m.Author.Mention(),
			cat.Name,
		))
		return nil
	}

	randUntil := getRandomAwayUntil(context.Config.CatMinAwayMins, context.Config.CatMaxAwayMins)

	err = cat.MarkAwayUntil(context.Store, m.ChannelID, randUntil)

	if err != nil {
		logrus.WithError(err).Errorln("Could not mark cat as away")
		return err
	}

	_, _ = ChannelMessageSendEmote(s, m.ChannelID, "💨", fmt.Sprintf(
		"%s, %s is now away!",
		m.Author.Mention(),
		cat.Name,
	))

	return nil
}

func getRandomAwayUntil(min int64, max int64) time.Time {
	now := time.Now().UTC()

	randDuration := time.Duration(rand.Int63n(max-min)+min) * time.Minute

	return now.Add(randDuration)
}
