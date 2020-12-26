package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

const bannedChars = "#<>@*~_`\\/-,!|"

func comesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
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
	return s.ChannelMessageSend(channelId, formatEmojiMessage(emote, message))
}

func formatEmojiMessage(emote string, message string) string {
	return fmt.Sprintf(" %s | %s", emote, message)
}
