package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type BotContext struct {
	Datastore *sqlx.DB
	Config    *Config
	Commands  map[string]*Command
}

func main() {
	cfg, err := LoadConfig()

	if err != nil {
		logrus.WithError(err).Fatal("Configuration error.")
	}

	// Seed random generator
	rand.Seed(time.Now().Unix())

	discord, err := initDiscord(cfg.BotToken)

	if err != nil {
		logrus.WithError(err).Fatal("Could not initalize Discord.")
	}

	defer discord.Close()

	db, err := NewDb(cfg.DbSrc)

	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database.")
	}

	// Register all commands
	cmds := RegisterCommands()

	botCtx := &BotContext{Config: cfg, Commands: cmds, Datastore: db}

	discord.AddHandler(messageCreate(botCtx))

	logrus.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	logrus.Println("Shutting down...")
}

func initDiscord(botToken string) (*discordgo.Session, error) {
	discord, err := discordgo.New(fmt.Sprintf("Bot %s", botToken))

	if err != nil {
		return nil, err
	}

	err = discord.Open()

	if err != nil {
		return nil, err
	}

	return discord, nil
}

func messageCreate(appContext *BotContext) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		err := handleMessage(s, m, appContext)

		if err != nil {
			logrus.WithError(err).Error("Error while handling message")
		}
	}
}
