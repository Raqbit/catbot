package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
)

func Profile(s *discordgo.Session, m *discordgo.MessageCreate,
	_ []string, context *CmdContext) error {

	cats, err := models.Cats.GetAllCatsOfUser(context.Store, context.User)

	if err != nil {
		return fmt.Errorf("could not get cats from db: %w", err)
	}

	catsValue := ""

	if len(cats) != 0 {
		for _, cat := range cats {
			status := "Home"
			if cat.Away {
				status = "Away"
			}
			catsValue = catsValue + fmt.Sprintf("%s - **%s**\n", cat.Name, status)
		}
	} else {
		catsValue = fmt.Sprintf("**None** :(\nBuy a cat with **%sbuy**.", context.Bot.Config.CommandPrefix)
	}

	profileDesc := ""
	if context.User.Money > 0 {
		profileDesc = fmt.Sprintf("💰 %d credits", context.User.Money)
	} else {
		profileDesc = fmt.Sprintf(
			"You don't have any money.\nUse **%sdaily** to get your daily reward.",
			context.Bot.Config.CommandPrefix,
		)
	}

	embed := &discordgo.MessageEmbed{
		Title:       m.Author.Username + "'s Profile",
		Description: profileDesc,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Cats",
				Value: catsValue,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Buy cats with %sbuy, let cats out with %sout.",
				context.Bot.Config.CommandPrefix,
				context.Bot.Config.CommandPrefix,
			),
		},
	}

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed)

	return nil
}
