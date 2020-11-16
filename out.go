package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"math/rand"
	"strings"
	"time"
)

func Out(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, ctx *CmdContext) error {
	if len(parts) < 2 {
		_, _ = ChannelMesageSendError(s, m.ChannelID, "Please specify a cat to get the info of!")
		return nil
	}

	catName := strings.Join(parts[1:], " ")

	cat, err := models.Cats.GetByName(ctx.Store, ctx.User, catName)

	if err != nil {
		return fmt.Errorf("could not fetch cat from db: %w", err)
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

	randUntil := getRandomAwayUntil(ctx.Bot.Config.CatMinAwayMins, ctx.Bot.Config.CatMaxAwayMins)

	err = cat.MarkAwayUntil(ctx.Store, m.ChannelID, randUntil)

	if err != nil {
		return fmt.Errorf("could not mark cat as away: %w", err)
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
