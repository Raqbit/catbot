package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func Profile(s *discordgo.Session, m *discordgo.MessageCreate,
	_ []string, globalEnv *GlobalEnv, cmdEnv *CommandEnv) (error) {
	cats, err := globalEnv.Db.AllCatsOfUser(cmdEnv.User.ID)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Errorln("Could not get cats from database")
		return err
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
		catsValue = fmt.Sprintf("**None** :(\nBuy a cat with **%sbuy**.", globalEnv.Config.CommandPrefix)
	}

	profileDesc := ""
	if cmdEnv.User.Money > 0 {
		profileDesc = fmt.Sprintf("💰 %d credits", cmdEnv.User.Money)
	} else {
		profileDesc = fmt.Sprintf(
			"You don't have any money.\nUse **%sdaily** to get your daily reward.",
			globalEnv.Config.CommandPrefix,
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
				globalEnv.Config.CommandPrefix,
				globalEnv.Config.CommandPrefix,
			),
		},
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)

	return nil
}
