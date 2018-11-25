package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"math"
	"raqb.it/catbot/utils"
	"time"
)

func Daily(s *discordgo.Session, m *discordgo.MessageCreate, _ []string, globalEnv *GlobalEnv, cmdEnv *CommandEnv) error {
	diff := time.Since(cmdEnv.User.LastDaily).Hours()

	if diff < 24 {
		utils.ChannelMesageSendError(s,
			m.ChannelID,
			fmt.Sprintf("%s, your your daily is still on cooldown! It refreshes in %.0f hours.",
				m.Author.Mention(),
				24-math.Round(diff),
			),
		)
		return nil
	}

	newAmount, err := globalEnv.Db.UserUseDaily(cmdEnv.User.ID, globalEnv.Config.CatCost)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("Failed updating money of user")
		return err
	}

	dailyTitle := fmt.Sprintf(
		"%s opens their daily...",
		m.Author.Username)

	dailyDesc := fmt.Sprintf(
		"**+%d credits** (you now have %d)",
		globalEnv.Config.CatCost,
		newAmount,
	)

	dailyFooter := fmt.Sprintf("Use %sbuy to buy a cat.",
		globalEnv.Config.CommandPrefix,
	)

	embed := &discordgo.MessageEmbed{
		Title:       dailyTitle,
		Description: dailyDesc,
		Color:       0x80df3a,
		Footer: &discordgo.MessageEmbedFooter{
			Text: dailyFooter,
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)

	return nil
}
