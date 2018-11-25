package main

import (
	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"raqb.it/catbot/config"
	"raqb.it/catbot/models"
	"syscall"
	"time"
)

type GlobalEnv struct {
	Db       models.Datastore
	Config   *config.Config
	Commands map[string]*Command
}

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Configuration error.")
	}

	if cfg.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	// Seed random generator
	rand.Seed(time.Now().Unix())

	discord, err := initDiscord(cfg.BotToken)
	defer discord.Close()

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to connect to Discord.")
	}

	db, err := models.NewDb(cfg.DbSrc)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to connect to database.")
	}

	// Register all commands
	cmds := RegisterCommands()

	globalEnv := &GlobalEnv{Db: db, Config: cfg, Commands: cmds}

	discord.AddHandler(messageCreate(globalEnv))

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

func messageCreate(globalEnv *GlobalEnv) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		HandleMessage(s, m, globalEnv)
	}
}
