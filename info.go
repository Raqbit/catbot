package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const ckBaseUrl = "https://storage.googleapis.com/ck-kitty-image/0x06012c8cf97bead5deae237070f9587f8e7a266d"

func Info(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, ctx *CmdContext) error {
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

	_, _ = s.ChannelMessageSendEmbed(m.ChannelID, createCatProfileEmbed(cat))

	return nil
}
