package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"math"
	"time"
)

func Daily(s *discordgo.Session, m *discordgo.MessageCreate, _ []string, context *Context) error {
	diff := time.Since(context.User.LastDaily).Hours()

	if diff < 24 {
		_, _ = ChannelMesageSendError(s,
			m.ChannelID,
			fmt.Sprintf("%s, your your daily is still on cooldown! It refreshes in %.0f hours.",
				m.Author.Mention(),
				24-math.Round(diff),
			),
		)
		return nil
	}

	// TODO: Add randomness?
	newAmount, err := context.User.UseDaily(context.Store, context.Config.CatCost)

	if err != nil {
		logrus.WithError(err).Errorln("Failed updating money of user")
		return err
	}

	dailyTitle := fmt.Sprintf(
		"%s opens their daily...",
		m.Author.Username)

	dailyDesc := fmt.Sprintf(
		"**+%d credits** (you now have %d)",
		context.Config.CatCost,
		newAmount,
	)

	dailyFooter := fmt.Sprintf("Use %sbuy to buy a cat.",
		context.Config.CommandPrefix,
	)

	embed := &discordgo.MessageEmbed{
		Title:       dailyTitle,
		Description: dailyDesc,
		Color:       0x80df3a,
		Footer: &discordgo.MessageEmbedFooter{
			Text: dailyFooter,
		},
	}

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed)

	return nil
}
