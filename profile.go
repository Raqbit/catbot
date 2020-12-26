package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func Profile(s *discordgo.Session, m *discordgo.MessageCreate,
	_ []string, ctx *CmdContext) error {

	if ctx.Cat == nil {
		_, _ = ChannelMesageSendError(s, m.ChannelID, fmt.Sprintf("You don't seem to have a cat yet! Get one with **%sstart**!", ctx.Bot.Config.CommandPrefix))
		return nil
	}

	catsValue := fmt.Sprintf("%s is %s\n", ctx.Cat.R.Type.Name, ctx.Cat.R.CurrentActivity.Description)

	moneyDesc := fmt.Sprintf(
		"You don't have any money.\nUse **%sdaily** to get your daily reward.",
		ctx.Bot.Config.CommandPrefix,
	)

	if ctx.User.Money > 0 {
		moneyDesc = fmt.Sprintf("You have %d 💰", ctx.User.Money)
	}

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       m.Author.Username + "'s Cat",
		Description: catsValue,
		Image: &discordgo.MessageEmbedImage{
			URL: ctx.Cat.R.Type.AvatarURL,
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: moneyDesc,
		},
	})

	return nil
}
