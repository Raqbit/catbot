package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"time"
)

func setupCatReturnCron(session *discordgo.Session, env *GlobalEnv) *time.Ticker {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			checkCatReturns(session, env)
		}
	}()

	return ticker
}

func checkCatReturns(s *discordgo.Session, env *GlobalEnv) {
	cats, err := env.Db.UpdateReturningCats()

	if err != nil {
		logrus.WithError(err).Error("Could not retrieve returning cats")
		return
	}

	for _, cat := range cats {
		user, err := env.Db.GetUserById(cat.OwnerId)

		if err != nil {
			logrus.WithError(err).Error("Could not retrieve cat owner")
			continue
		}

		discordUser, err := s.User(user.DiscordId)

		if err != nil {
			logrus.WithError(err).Error("Could not retrieve discord user of cat owner")
			continue
		}

		_, _ = ChannelMessageSendEmote(s, cat.AwayChannel, "",
			fmt.Sprintf(
				"%s, %s has returned!",
				discordUser.Mention(),
				cat.Name,
			))
	}
}
