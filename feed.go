package main

import (
	"github.com/bwmarrin/discordgo"
)

func Feed(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, globalEnv *GlobalEnv, cmdEnv *CommandEnv) error {
	s.ChannelMessageSend(m.ChannelID, "Not implemented")
	return nil
}
