package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"time"
)

func setupCatReturnCron(session *discordgo.Session, botContext *BotContext) *time.Ticker {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		for range ticker.C {
			err := checkCatReturns(session, botContext)

			if err != nil {
				logrus.WithError(err).Error("Could not check cat returns")
			}
		}
	}()

	return ticker
}

func checkCatReturns(s *discordgo.Session, botContext *BotContext) error {
	cats, err := models.Cats.UpdateReturning(botContext.Datastore)

	if err != nil {
		return fmt.Errorf("could not retrieve returning cats: %w", err)
	}

	for _, cat := range cats {
		user, err := models.Users.GetById(botContext.Datastore, cat.OwnerId)

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

	return nil
}
