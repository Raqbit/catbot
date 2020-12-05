package main

import (
	"context"
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
)

func Profile(s *discordgo.Session, m *discordgo.MessageCreate,
	_ []string, ctx *CmdContext) error {

	cat, err := models.Cats(
		models.CatWhere.OwnerID.EQ(ctx.User.ID),
		qm.Load(models.CatRels.Type),
		qm.Load(models.CatRels.CurrentActivity),
	).One(context.Background(), ctx.Store)

	if err != nil {
		return fmt.Errorf("could not get cat from db: %w", err)
	}

	catsValue := fmt.Sprintf("%s is %s\n", cat.R.Type.Name, cat.R.CurrentActivity.Description)

	profileDesc := ""
	if ctx.User.Money > 0 {
		profileDesc = fmt.Sprintf("💰 %d credits", ctx.User.Money)
	} else {
		profileDesc = fmt.Sprintf(
			"You don't have any money.\nUse **%sdaily** to get your daily reward.",
			ctx.Bot.Config.CommandPrefix,
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
				ctx.Bot.Config.CommandPrefix,
				ctx.Bot.Config.CommandPrefix,
			),
		},
	}

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, embed)

	return nil
}
