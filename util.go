package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"strings"
)

const bannedChars = "#<>@*~_`\\/-,!|"

func isValidCatName(name string) bool {
	if len(name) > 20 {
		return false
	}
	if strings.ContainsAny(name, bannedChars) {
		return false
	}

	return true
}

func createCatProfileEmbed(cat *models.Cat) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title: fmt.Sprintf("Your %s", cat.R.Type.Name),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: cat.R.Type.AvatarURL,
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Status",
				Value: cat.R.CurrentActivity.Description,
			},
		},
	}
}

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
