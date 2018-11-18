package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"raqb.it/catbot/models"
	"raqb.it/catbot/utils"
	"strings"
)

type CommandEnv struct {
	User *models.User
}

func RegisterCommands() (map[string]*Command) {
	cmds := make(map[string]*Command)

	addCommand(cmds, &Command{
		Name:        "Buy",
		Description: "Buy cats",
		Aliases:     []string{"buy"},
		Admin:       false,
		Exec:        Buy,
	})

	addCommand(cmds, &Command{
		Name:        "Out",
		Description: "Let out a cat",
		Aliases:     []string{"out"},
		Admin:       false,
		Exec:        Out,
	})

	addCommand(cmds, &Command{
		Name:        "Feed",
		Description: "Feed a cat",
		Aliases:     []string{"feed"},
		Admin:       false,
		Exec:        Feed,
	})

	addCommand(cmds, &Command{
		Name:        "Profile",
		Description: "Get your profile",
		Aliases:     []string{"profile"},
		Admin:       false,
		Exec:        Profile,
	})

	addCommand(cmds, &Command{
		Name:        "Daily",
		Description: "Get your daily reward",
		Aliases:     []string{"daily"},
		Admin:       false,
		Exec:        Daily,
	})

	addCommand(cmds, &Command{
		Name:        "Info",
		Description: "Get info about one of your cats",
		Aliases:     []string{"info", "catinfo"},
		Admin:       false,
		Exec:        Info,
	})

	return cmds
}

func addCommand(cmds map[string]*Command, command *Command) {
	for _, alias := range command.Aliases {

		_, commandFound := cmds[alias]

		// Give error if label was already registered
		if commandFound {
			logrus.WithFields(logrus.Fields{
				"alias":       alias,
				"old_command": cmds[alias].Name,
				"new_command": command.Name,
			}).Errorln("Duplicate commands alias detected, overwriting!")
		}

		// Overwite old commands
		cmds[alias] = command
	}
}

func HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate, globalEnv *GlobalEnv) {
	// Do not respond to messages of myself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Check if message comes from DM
	comesFromDM, err := utils.ComesFromDM(s, m)

	if err != nil {
		logrus.Error("Could not determine if message was sent in public channel")
		return
	}

	// Do not respond to dms
	if comesFromDM {
		return
	}

	// Do not respond to bots
	if m.Author.Bot {
		return
	}

	// Split by space
	commandParts := strings.Split(m.Content, " ") // c!foo bar baz

	// Only respond if message has correct prefix
	if !strings.HasPrefix(commandParts[0], globalEnv.Config.CommandPrefix) {
		return
	}

	// Get commands label by trimming prefix
	label := commandParts[0][len(globalEnv.Config.CommandPrefix):]

	// Empty string label
	if len(label) == 0 {
		return
	}

	cmd, commandFound := globalEnv.Commands[label]

	// Only respond if commands is known
	if !commandFound {
		return
	}

	user, err := globalEnv.Db.GetUserOrCreate(m.Author.ID)

	if err != nil {
		logrus.Error("Could not fetch or create user!")
		return
	}

	cmdEnv := &CommandEnv{User: user}

	// Execute commands
	err = cmd.Exec(s, m, commandParts, globalEnv, cmdEnv)

	if err != nil {
		utils.ChannelMesageSendError(s, m.ChannelID,
			fmt.Sprintf("Something went wrong while executing %s%s",
				globalEnv.Config.CommandPrefix,
				label,
			),
		)
		return
	}
}
