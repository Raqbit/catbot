package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const bannedChars = "#<>@*~_`\\/-,!|"

func CheckCatName(name string) bool {
	if len(name) > 20 {
		return false
	}
	if strings.ContainsAny(name, bannedChars) {
		return false
	}

	return true
}

func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

func ChannelMesageSendError(s *discordgo.Session, channelId string, message string) (*discordgo.Message, error) {
	return ChannelMessageSendEmote(s, channelId, "🚫", message)
}

func ChannelMessageSendEmote(s *discordgo.Session, channelId string, emote string, message string) (*discordgo.Message, error) {
	return s.ChannelMessageSend(channelId, fmt.Sprintf(" %s | %s", emote, message))
}
