package main

import (
	"context"
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"math"
	"time"
)

func Daily(s *discordgo.Session, m *discordgo.MessageCreate, _ []string, ctx *CmdContext) error {
	diff := time.Since(ctx.User.LastDaily).Hours()

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

	ctx.User.LastDaily = time.Now()
	// TODO: Add randomness?
	ctx.User.Money -= ctx.Bot.Config.CatCost

	_, err := ctx.User.Update(
		context.Background(),
		ctx.Store,
		boil.Whitelist(
			models.UserColumns.LastDaily,
			models.UserColumns.Money,
		),
	)

	if err != nil {
		return fmt.Errorf("failed using user's daily: %w", err)
	}

	dailyTitle := fmt.Sprintf(
		"%s opens their daily reward...",
		m.Author.Username)

	dailyDesc := fmt.Sprintf(
		"**+%d credits** (you now have %d)",
		ctx.Bot.Config.CatCost,
		ctx.User.Money,
	)

	dailyFooter := fmt.Sprintf("Use %sbuy to buy a cat.",
		ctx.Bot.Config.CommandPrefix,
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
