package main

import (
	"github.com/bwmarrin/discordgo"
)

func Feed(s *discordgo.Session, m *discordgo.MessageCreate, _ []string, _ *Context) error {
	_, _ = s.ChannelMessageSend(m.ChannelID, "Not implemented")
	return nil
}
