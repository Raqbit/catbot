package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
)

func Out(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, globalEnv *GlobalEnv, cmdEnv *CommandEnv) error {

	if len(parts) < 2 {
		_, _ = ChannelMesageSendError(s, m.ChannelID, "Please specify a cat to get the info of!")
		return nil
	}

	catName := strings.Join(parts[1:], " ")

	cat, err := globalEnv.Db.GetCatByName(cmdEnv.User.ID, catName)

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

	randUntil := getRandomAwayUntil(globalEnv.Config.CatMinAwayMins, globalEnv.Config.CatMaxAwayMins)

	err = globalEnv.Db.MarkCatAwayUntil(cat.ID, m.ChannelID, randUntil)

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
