package main

import (
	"fmt"
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"strings"
)

type CmdContext struct {
	Bot   *BotContext
	User  *models.User
	Store models.Queryable
}

func RegisterCommands() map[string]*Command {
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
		Description: "Let out a cat out on an adventure",
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

func handleMessage(s *discordgo.Session, m *discordgo.MessageCreate, appContext *BotContext) error {
	// Do not respond to bots
	if m.Author.Bot {
		return nil
	}

	// Do not respond to messages of myself
	if m.Author.ID == s.State.User.ID {
		return nil
	}

	// Check if message comes from DM
	fromDM, err := comesFromDM(s, m)

	if err != nil {
		return fmt.Errorf("could not determine if msg was sent in public channel: %w", err)
	}

	// Do not respond to dms
	if fromDM {
		return nil
	}

	// Split by space
	commandParts := strings.Split(m.Content, " ") // c!foo bar baz

	// Only respond if message has correct prefix
	if !strings.HasPrefix(commandParts[0], appContext.Config.CommandPrefix) {
		return nil
	}

	// Get commands label by trimming prefix
	label := commandParts[0][len(appContext.Config.CommandPrefix):]

	// Empty string label
	if len(label) == 0 {
		return nil
	}

	cmd, commandFound := appContext.Commands[label]

	// Only respond if commands is known
	if !commandFound {
		return nil
	}

	user, err := models.Users.GetOrCreate(appContext.Datastore, m.Author.ID)

	if err != nil {
		return fmt.Errorf("could not get or create user: %w", err)
	}

	tx, err := appContext.Datastore.BeginTransaction()

	if err != nil {
		return fmt.Errorf("could not create transaction: %w", err)
	}

	ctx := &CmdContext{Bot: appContext, User: user, Store: tx}

	// Execute commands
	err = cmd.Exec(s, m, commandParts, ctx)

	if err != nil {
		// Rollback changes when an error ocurred during execution of the command, but do not override err with nil
		_ = tx.Rollback()
	} else {
		err = tx.Commit()
	}

	if err != nil {
		_, _ = ChannelMesageSendError(s, m.ChannelID,
			fmt.Sprintf("Something went wrong while executing %s%s",
				appContext.Config.CommandPrefix,
				label,
			),
		)
	}

	return err
}
