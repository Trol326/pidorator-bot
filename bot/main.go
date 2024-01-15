package bot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"pidorator-bot/bot/trigger"
	"pidorator-bot/database"
	"pidorator-bot/database/mongodb"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

type Client struct {
	Session *discordgo.Session
	DB      database.Database
	Log     *zerolog.Logger
}

func New() (Client, error) {
	ctx := context.Background()
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC822}
	log := zerolog.New(consoleWriter).With().Timestamp().Logger()

	develop := os.Getenv("DEVELOP")
	if develop != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("develop")
	}

	key := os.Getenv("KEY")
	if key == "" {
		return Client{}, fmt.Errorf("key not found")
	}

	data := fmt.Sprintf("Bot %s", key)
	s, err := discordgo.New(data)
	if err != nil {
		return Client{}, err
	}

	db := mongodb.New(ctx, &log)

	client := Client{
		Session: s,
		DB:      &db,
		Log:     &log,
	}

	return client, nil
}

func (c *Client) Start() {
	c.Session.Open()
	defer c.Session.Close()
	defer c.DB.Disconnect()

	c.Log.Info().Msg("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (c *Client) InitHandlers() {
	t := trigger.New(c.Log)

	c.Session.AddHandler(t.OnNewMessage)
}
