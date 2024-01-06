package bot

import (
	"fmt"
	"os"
	"os/signal"
	"pidorator-bot/bot/trigger"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog"
)

type Client struct {
	Session *discordgo.Session
	Log     *zerolog.Logger
}

func New() (Client, error) {
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

	client := Client{
		Session: s,
		Log:     &log,
	}

	return client, nil
}

func (c *Client) Start() {
	c.Session.Open()
	defer c.Session.Close()

	c.Log.Info().Msg("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func (c *Client) InitHandlers() {
	t := trigger.New(c.Log)

	c.Session.AddHandler(t.OnNewMessage)
}
