package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"time"
)

func setupCatReturnCron(session *discordgo.Session, appContext *AppContext) *time.Ticker {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			checkCatReturns(session, appContext)
		}
	}()

	return ticker
}

func checkCatReturns(s *discordgo.Session, appContext *AppContext) {
	cats, err := models.Cats.UpdateReturning(appContext.Store)

	if err != nil {
		logrus.WithError(err).Error("Could not retrieve returning cats")
		return
	}

	for _, cat := range cats {
		user, err := models.Users.GetById(appContext.Store, cat.OwnerId)

		if err != nil {
			logrus.WithError(err).Error("Could not retrieve cat owner")
			continue
		}

		discordUser, err := s.User(user.DiscordId)

		if err != nil {
			logrus.WithError(err).Error("Could not retrieve discord user of cat owner")
			continue
		}

		_, _ = ChannelMessageSendEmote(s, cat.AwayChannel, "🏘️",
			fmt.Sprintf(
				"%s, %s has returned!",
				discordUser.Mention(),
				cat.Name,
			))
	}
}
