package main

import (
	"github.com/Raqbit/catbot/models"
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type AppContext struct {
	Store    models.Querier
	Config   *Config
	Commands map[string]*Command
}

func main() {
	cfg, err := LoadConfig()

	if err != nil {
		logrus.WithError(err).Fatal("Configuration error.")
	}

	// Seed random generator
	rand.Seed(time.Now().Unix())

	discord, err := initDiscord(cfg.BotToken)
	defer discord.Close()

	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to Discord.")
	}

	db, err := NewDb(cfg.DbSrc)

	if err != nil {
		logrus.WithError(err).Fatal("Failed to connect to database.")
	}

	// Register all commands
	cmds := RegisterCommands()

	appContext := &AppContext{Store: db, Config: cfg, Commands: cmds}

	// Setup cat return cron
	catReturnTicker := setupCatReturnCron(discord, appContext)
	defer catReturnTicker.Stop()

	discord.AddHandler(messageCreate(appContext))

	logrus.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	logrus.Println("Shutting down...")
}

func initDiscord(botToken string) (*discordgo.Session, error) {
	discord, err := discordgo.New("Bot " + botToken)

	if err != nil {
		return nil, err
	}

	err = discord.Open()

	if err != nil {
		return nil, err
	}

	return discord, nil
}

func messageCreate(appContext *AppContext) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		HandleMessage(s, m, appContext)
	}
}
