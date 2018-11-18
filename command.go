package main

import (
	"github.com/bwmarrin/discordgo"
)

type RunFunc func(s *discordgo.Session, m *discordgo.MessageCreate, parts []string, globalEnv *GlobalEnv, cmdEnv *CommandEnv) error

type Command struct {
	Name        string
	Description string
	Aliases     []string
	Admin       bool
	Exec        RunFunc
}

func (command *Command) AddAlias(alias string) {
	command.Aliases = append(command.Aliases, alias)
}
